# log_mac_sleep

- Take sleep / wakeup log
- Report daily active time with the log

### Download
```bash
go get github.com/kuxuxun/log_mac_sleep/cmd/log_mac_sleep
```

### Output file path
- log file path:	~/.sleeplog/log
- daily_report_file:	~/.sleeplog/daily_active

### Run
```bash
# start log
log_mac_sleep -watch
# output daily report
log_mac_sleep -daily
```

register to launchd
```bash
go get github.com/kuxuxun/log_mac_sleep
cd $GOPATH/src/github.com/kuxuxun/log_mac_sleep/_plist
sh register_launchd.sh
```

