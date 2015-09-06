package main

/* example: output comment to console on sleep/wake up */
import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sync"
	"time"

	"github.com/kuxuxun/mac_switch_watch"
)

var (
	mutex       = &sync.Mutex{}
	jst         = time.FixedZone("Asia/Tokyo", 9*60*60)
	LogTimeFmt  = "2006-01-02 15:04:05"
	LogFileName = ".sleeplog/log"
)

func main() {
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

	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		err := os.Mkdir(filepath.Dir(logFilePath), 0666)
		if err != nil {
			panic(err)
		}

		_, err = os.Create(logFilePath)
		if err != nil {
			panic(err)
		}
	}

	log := fmt.Sprintf("%s:%s", msg, now)
	file, _ := os.OpenFile(logFilePath, os.O_APPEND, 0666)
	file.WriteString(log)

	defer file.Close()
}
