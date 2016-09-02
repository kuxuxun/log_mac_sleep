package log_mac_sleep

/* example: output comment to console on sleep/wake up */
import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/tacogips/mac_switch_watch"
)

var (
	mutex                = &sync.Mutex{}
	jst                  = time.FixedZone("Asia/Tokyo", 9*60*60)
	LogTimeFmt           = "2006-01-02 15_04_05"
	outDateFmt           = "2006-01-02"
	outTimeFmt           = "15:04:05"
	LogFileName          = ".sleeplog/log"
	AggrFileName         = ".sleeplog/daily_active"
	endToBeginThreshold  = 6 * time.Hour           // rolling a day if interval of end to begin over this term.
	sameActionsThreshold = 3 * endToBeginThreshold // rolling a day if interval of begin to begin over this term. (suggesting there were no "end" action between this term)
	locationJP, _        = time.LoadLocation("Asia/Tokyo")
)

const (
	IS_START int = iota
	IS_END
)

var actionTypes = map[string]int{
	"start":    IS_START,
	"wakeup":   IS_START,
	"sleep":    IS_END,
	"poweroff": IS_END,
}

func Aggregate() {
	aggrLog()
}

func Start() {
	mac_switch_watch.SetHandler(mac_switch_watch.KeyOnSystemWillSleep, OnSleep)
	mac_switch_watch.SetHandler(mac_switch_watch.KeyOnSystemWillPowerOn, OnWakeup)
	mac_switch_watch.SetHandler(mac_switch_watch.KeyOnSystemWillPowerOff, OnPowerOff)

	OnStart()

	mac_switch_watch.Watch()
}

func OnStart() {
	logTimeToFile("start")
}

func OnSleep() {
	logTimeToFile("sleep")
}

func OnPowerOff() {
	logTimeToFile("poweroff")
}

func OnWakeup() {
	logTimeToFile("wakeup")
}

func fmtJst(t time.Time, format string) string {
	return t.In(jst).Format(format)
}

func logTimeToFile(msg string) {
	usr, _ := user.Current()
	homeDir := usr.HomeDir
	logFilePath := filepath.Join(homeDir, LogFileName)

	now := fmtJst(time.Now(), LogTimeFmt)

	mutex.Lock()
	defer mutex.Unlock()

	if _, err := os.Stat(filepath.Dir(logFilePath)); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(logFilePath), 0744)
	}

	log := fmt.Sprintf("%s:%s\n", msg, now)
	file, _ := os.OpenFile(logFilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	file.WriteString(log)

	defer file.Close()
}

func aggrLog() {
	usr, _ := user.Current()
	homeDir := usr.HomeDir
	logFilePath := filepath.Join(homeDir, LogFileName)

	if _, err := os.Stat(filepath.Dir(logFilePath)); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(logFilePath), 0744)
	}

	logFile, err := os.Open(logFilePath)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	logFileScanner := bufio.NewScanner(logFile)

	workingTimeADays, err := aggregate(logFileScanner)
	if err != nil {
		panic(err)
	}

	aggrFilePath := filepath.Join(homeDir, AggrFileName)
	aggrfile, _ := os.OpenFile(aggrFilePath, os.O_RDWR|os.O_CREATE, 0644)
	defer aggrfile.Close()
	for _, workTime := range workingTimeADays {
		aggrfile.WriteString(workTime.ToCsvLine())
		print(workTime.ToCsvLine())
	}
}

type WorkingTimeADay struct {
	Start time.Time
	End   time.Time
}

func (w WorkingTimeADay) ToCsvLine() string {
	var dt string
	if w.Start.IsZero() {
		dt = fmtJst(w.End, outDateFmt)
	} else {
		dt = fmtJst(w.Start, outDateFmt)
	}

	s, e := "", ""

	if !w.Start.IsZero() {
		s = fmtJst(w.Start, outTimeFmt)
	}

	if !w.End.IsZero() {
		e = fmtJst(w.End, outTimeFmt)
	}

	return fmt.Sprintf("%s,%s,%s\n", dt, s, e)
}

func aggregate(logFileScanner *bufio.Scanner) ([]WorkingTimeADay, error) {
	var workingTimes []WorkingTimeADay
	workingTime := WorkingTimeADay{}

	for logFileScanner.Scan() {
		line := logFileScanner.Text()
		line = strings.TrimSpace(line)
		cols := strings.Split(logFileScanner.Text(), ":")
		if len(cols) != 2 {
			panic(errors.New(fmt.Sprintf("invalid line %#v", cols)))
		}
		action, timeStr := cols[0], cols[1]
		actionType, ok := actionTypes[action]
		if !ok {
			panic(errors.New(fmt.Sprintf("invalid action %s", action)))
		}

		logedTime, err := time.ParseInLocation(LogTimeFmt, timeStr, locationJP)
		if err != nil {
			return nil, err
		}

		switch actionType {
		case IS_START:
			if workingTime.Start.IsZero() {
				workingTime.Start = logedTime
			}
			if workingTime.End.IsZero() {
				if logedTime.After(workingTime.Start.Add(sameActionsThreshold)) {
					workingTimes = append(workingTimes, workingTime)
					workingTime = WorkingTimeADay{}
					workingTime.Start = logedTime
				}
			} else if logedTime.After(workingTime.End.Add(endToBeginThreshold)) {
				workingTimes = append(workingTimes, workingTime)
				workingTime = WorkingTimeADay{}
				workingTime.Start = logedTime
			} else {
				continue
			}

		case IS_END:
			if !workingTime.End.IsZero() && logedTime.After(workingTime.End.Add(sameActionsThreshold)) {
				workingTimes = append(workingTimes, workingTime)
				workingTime = WorkingTimeADay{}
			}
			workingTime.End = logedTime

		default:
			panic("not impl")

		}
	}

	workingTimes = append(workingTimes, workingTime)
	return workingTimes, nil
}
