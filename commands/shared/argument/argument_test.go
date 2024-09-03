package argument_test

import (
	"flag"
	"testing"
	"time"

	"github.com/ministryofjustice/opg-reports/commands/shared/argument"
	"github.com/ministryofjustice/opg-reports/shared/dates"
)

func TestSharedArgumentDatePassed(t *testing.T) {
	var (
		now       = time.Now().UTC()
		yesterday = dates.ResetDay(now).AddDate(0, 0, -1)
		testDay   = yesterday.AddDate(0, -1, 0)
		group     = flag.NewFlagSet("test_day", flag.ExitOnError)
		day       = argument.NewDate(group, "day", yesterday, dates.FormatYMD, "Day (YYYY-MM-DD) to fetch uptime data for")
	)

	// test with correct date added
	group.Parse([]string{"-day=" + testDay.Format(dates.FormatYMD)})
	if day.Value.Format(dates.FormatYMD) != testDay.Format(dates.FormatYMD) {
		t.Errorf("day should be testDay, actual [%v]", day.Value)
	}
}

func TestSharedArgumentDateEmpty(t *testing.T) {
	var (
		now       = time.Now().UTC()
		yesterday = dates.ResetDay(now).AddDate(0, 0, -1)
		group     = flag.NewFlagSet("test_day", flag.ExitOnError)
		day       = argument.NewDate(group, "day", yesterday, dates.FormatYMD, "Day (YYYY-MM-DD) to fetch uptime data for")
	)

	// test with now date added - so should be default
	group.Parse([]string{})
	if day.Value.Format(dates.FormatYMD) != yesterday.Format(dates.FormatYMD) {
		t.Errorf("day should be yesterday, actual [%v]", day.Value)
	}
}

func TestSharedArgumentDateDash(t *testing.T) {
	var (
		now       = time.Now().UTC()
		yesterday = dates.ResetDay(now).AddDate(0, 0, -1)
		group     = flag.NewFlagSet("test_day", flag.ExitOnError)
		day       = argument.NewDate(group, "day", yesterday, dates.FormatYMD, "Day (YYYY-MM-DD) to fetch uptime data for")
	)

	// test with now date added - so should be default
	group.Parse([]string{"-day=-"})
	if day.Value.Format(dates.FormatYMD) != yesterday.Format(dates.FormatYMD) {
		t.Errorf("day should be yesterday, actual [%v]", day.Value)
	}
}
