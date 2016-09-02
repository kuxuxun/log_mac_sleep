package main

import (
	"flag"
	"fmt"

	lw "github.com/tacogips/log_mac_sleep"
)

var (
	watch     = flag.Bool("watch", false, "")
	aggregate = flag.Bool("daily", false, "")
)

func main() {
	flag.Parse()
	if *watch {
		lw.Start()
	} else if *aggregate {
		lw.Aggregate()
	} else {
		fmt.Printf(" Error: No argument. \n usage : log_mac_sleep [-watch | -daily] '")
	}
}
