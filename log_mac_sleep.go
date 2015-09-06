package main

/* example: output comment to console on sleep/wake up */
import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/kuxuxun/mac_switch_watch"
)

var (
	mutex       = &sync.Mutex{}
	jst         = time.FixedZone("Asia/Tokyo", 9*60*60)
	LogTimeFmt  = "2006-01-02 15:04:05"
	LogFilePath = "~/.sleeplog/log"
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
	now := fmtJst(time.Now(), LogTimeFmt)

	mutex.Lock()
	defer mutex.Unlock()

	if _, err := os.Stat(LogFilePath); os.IsNotExist(err) {
		os.Mkdir(filepath.Dir(LogFilePath), 0666)
		_, err = os.Create(LogFilePath)
		if err != nil {
			panic(err)
		}
	}

	log := fmt.Sprintf("%s:%s", msg, now)
	file, _ := os.OpenFile(LogFilePath, os.O_APPEND, 0666)
	file.WriteString(log)
	defer file.Close()
}
