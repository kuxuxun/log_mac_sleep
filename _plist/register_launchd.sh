#!/bin/sh
cp  $GOPATH/src/github.com/kuxuxun/log_mac_sleep/_plist/github.com.kuxuxun.log_mac_sleep.plist  ~/Library/LaunchAgents/
launchctl load ~/Library/LaunchAgents/github.com.kuxuxun.log_mac_sleep.plist
launchctl start github.com.kuxuxun.log_mac_sleep.plist
