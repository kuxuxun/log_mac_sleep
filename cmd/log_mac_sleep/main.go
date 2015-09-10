package main

import (
	"flag"
	"fmt"

	lw "github.com/kuxuxun/log_mac_sleep"
)

var (
	watch     = flag.Bool("watch", false, "")
	aggregate = flag.Bool("aggr", false, "")
)

func main() {
	flag.Parse()
	if *watch {
		lw.Start()
	} else if *aggregate {
		lw.Aggregate()
	} else {
		fmt.Printf(" Error: No argument. \n usage : log_mac_sleep [-watch | -aggr] '  ")
	}
}