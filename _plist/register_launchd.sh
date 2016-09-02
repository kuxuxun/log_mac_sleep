#!/bin/sh
cp $GOPATH/src/github.com/tacogips/log_mac_sleep/_plist/github.com.tacogips.log_mac_sleep.plist  ~/Library/LaunchAgents/
launchctl load ~/Library/LaunchAgents/github.com.tacogips.log_mac_sleep.plist
launchctl start github.com.tacogips.log_mac_sleep.plist
