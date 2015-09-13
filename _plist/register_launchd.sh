#!/bin/sh
cp $GOPATH/src/github.com/kuxuxun/log_mac_sleep/_plist/github.com.kuxuxun.log_mac_sleep.plist  ~/Library/LaunchAgents/
launchctl setenv PATH $GOPATH/bin:$PATH
launchctl load ~/Library/LaunchAgents/github.com.kuxuxun.log_mac_sleep.plist
launchctl start github.com.kuxuxun.log_mac_sleep.plist
