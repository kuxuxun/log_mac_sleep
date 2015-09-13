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

### register to launchd
```bash
go get github.com/kuxuxun/log_mac_sleep
```

copy _plist/github.com.kuxuxun.log_mac_sleep.plist.example to _plist/github.com.kuxuxun.log_mac_sleep.plist.
replace log_mac_sleep path with absolute path in your environment.
```
<array>
  <string>/path/to/log_mac_sleep</string>
  <string>-watch</string>
</array>
```

then load plist

```bash
sh register_launchd.sh
```

