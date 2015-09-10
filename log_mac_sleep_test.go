package log_mac_sleep

import (
	"bufio"
	"fmt"
	"strings"
	"testing"
)

func testTypical(t *testing.T) {
	data := `start:2015-09-03 10:00:00
sleep:2015-09-03 10:50:10
wakeup:2015-09-03 10:51:30
sleep:2015-09-03 10:50:10
wakeup:2015-09-03 10:51:30
sleep:2015-09-03 19:51:10
wakeup:2015-09-04 09:55:30
sleep:2015-09-04 14:50:10
wakeup:2015-09-04 18:50:10
poweroff:2015-09-04 18:50:10
wakeup:2015-09-05 10:00:00`

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
		if wt.ToCsvLine() != "2015-09-03,10:00:00,19:51:10" {
			t.Fatal(fmt.Sprintf("invalid csv line %s", wt.ToCsvLine()))
		}
	}

	{
		wt := result[1]
		if wt.ToCsvLine() != "2015-09-04,09:55:30,18:50:10" {
			t.Fatal(fmt.Sprintf("invalid csv line %s", wt.ToCsvLine()))
		}
	}

	{
		wt := result[2]
		if wt.ToCsvLine() != "2015-09-05,10:00:00," {
			t.Fatal(fmt.Sprintf("invalid csv line %s", wt.ToCsvLine()))
		}
	}

}

func testSerialStarts(t *testing.T) {
	data := `start:2015-09-03 10:00:00
start:2015-09-04 10:00:00
sleep:2015-09-04 10:50:10`

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
		if wt.ToCsvLine() != "2015-09-03,10:00:00," {
			t.Fatal(fmt.Sprintf("invalid csv line %s", wt.ToCsvLine()))
		}
	}

	{
		wt := result[1]
		if wt.ToCsvLine() != "2015-09-04,10:00:00,10:50:10" {
			t.Fatal(fmt.Sprintf("invalid csv line %s", wt.ToCsvLine()))
		}
	}
}

func testSerialEnds(t *testing.T) {
	data := `sleep:2015-09-03 10:00:00
sleep:2015-09-04 10:50:10`

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
		if wt.ToCsvLine() != "2015-09-03,,10:00:00" {
			t.Fatal(fmt.Sprintf("invalid csv line %s", wt.ToCsvLine()))
		}
	}

	{
		wt := result[1]
		if wt.ToCsvLine() != "2015-09-04,,10:50:10" {
			t.Fatal(fmt.Sprintf("invalid csv line %s", wt.ToCsvLine()))
		}
	}

}
