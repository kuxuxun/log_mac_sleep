package log_mac_sleep

import (
	"bufio"
	"fmt"
	"strings"
	"testing"
)

func TestTypical(t *testing.T) {
	data := `start:2015-09-03 10_00_00
sleep:2015-09-03 10_50_10
wakeup:2015-09-03 10_51_30
sleep:2015-09-03 10_50_10
wakeup:2015-09-03 10_51_30
sleep:2015-09-03 19_51_10
wakeup:2015-09-04 09_55_30
sleep:2015-09-04 14_50_10
wakeup:2015-09-04 18_50_10
poweroff:2015-09-04 18_50_10
wakeup:2015-09-05 10_00_00`

	sc := bufio.NewScanner(strings.NewReader(data))
	result, err := aggregate(sc)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 3 {
		t.Fatal("work days num is not 3")
	}

	{
		wt := result[0]
		if wt.ToCsvLine() != "2015-09-03,10:00:00,19:51:10\n" {
			t.Fatal(fmt.Sprintf("invalid csv line %s", wt.ToCsvLine()))
		}
	}

	{
		wt := result[1]
		if wt.ToCsvLine() != "2015-09-04,09:55:30,18:50:10\n" {
			t.Fatal(fmt.Sprintf("invalid csv line %s", wt.ToCsvLine()))
		}
	}

	{
		wt := result[2]
		if wt.ToCsvLine() != "2015-09-05,10:00:00,\n" {
			t.Fatal(fmt.Sprintf("invalid csv line %s", wt.ToCsvLine()))
		}
	}

}

func TestSerialStarts(t *testing.T) {
	data := `start:2015-09-03 10_00_00
start:2015-09-04 10_00_00
sleep:2015-09-04 10_50_10`

	sc := bufio.NewScanner(strings.NewReader(data))
	result, err := aggregate(sc)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatal("work days num is not 2")
	}

	{
		wt := result[0]
		if wt.ToCsvLine() != "2015-09-03,10:00:00,\n" {
			t.Fatal(fmt.Sprintf("invalid csv line %s", wt.ToCsvLine()))
		}
	}

	{
		wt := result[1]
		if wt.ToCsvLine() != "2015-09-04,10:00:00,10:50:10\n" {
			t.Fatal(fmt.Sprintf("invalid csv line %s", wt.ToCsvLine()))
		}
	}
}

func TestOneday(t *testing.T) {
	data := `start:2015-09-03 10_00_00
sleep:2015-09-03 10_50_10
wakeup:2015-09-03 10_51_30
sleep:2015-09-03 10_51_20`

	sc := bufio.NewScanner(strings.NewReader(data))
	result, err := aggregate(sc)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatal("work days num is not 1")
	}

	{
		wt := result[0]
		if wt.ToCsvLine() != "2015-09-03,10:00:00,10:51:20\n" {
			t.Fatal(fmt.Sprintf("invalid csv line %s", wt.ToCsvLine()))
		}
	}

}

func TestSerialEnds(t *testing.T) {
	data := `sleep:2015-09-03 10_00_00
sleep:2015-09-04 10_50_10`

	sc := bufio.NewScanner(strings.NewReader(data))
	result, err := aggregate(sc)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatal("work days num is not 2")
	}

	{
		wt := result[0]
		if wt.ToCsvLine() != "2015-09-03,,10:00:00\n" {
			t.Fatal(fmt.Sprintf("invalid csv line %s", wt.ToCsvLine()))
		}
	}

	{
		wt := result[1]
		if wt.ToCsvLine() != "2015-09-04,,10:50:10\n" {
			t.Fatal(fmt.Sprintf("invalid csv line %s", wt.ToCsvLine()))
		}
	}

}
