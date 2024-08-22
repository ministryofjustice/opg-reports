package dates_test

import (
	"testing"
	"time"

	"github.com/ministryofjustice/opg-reports/shared/dates"
	"github.com/ministryofjustice/opg-reports/shared/fake"
	"github.com/ministryofjustice/opg-reports/shared/testhelpers"
)

func TestSharedDatesBillingDates(t *testing.T) {
	var s time.Time
	var e time.Time
	// billing date of the 5th
	// current date of the 4th
	// ask for 3 months
	var (
		billingday int       = 5
		months     int       = 3
		current    time.Time = time.Date(2023, 10, 4, 0, 0, 0, 0, time.UTC)
		expectedS  time.Time = time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)
		expectedE  time.Time = time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC)
	)
	s, e = dates.BillingDates(current, billingday, months)

	if s.Format(dates.FormatYMD) != expectedS.Format(dates.FormatYMD) {
		t.Errorf("start date mis match expected [%s] actual [%s]", expectedS.String(), s.String())

	}
	if e.Format(dates.FormatYMD) != expectedE.Format(dates.FormatYMD) {
		t.Errorf("end date mis match expected [%s] actual [%s]", expectedE.String(), e.String())

	}

	// billing date of the 5th
	// current date of the 28th
	// ask for 6 months
	billingday = 5
	months = 6
	current = time.Date(2023, 11, 28, 0, 0, 0, 0, time.UTC)
	expectedS = time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC)
	expectedE = time.Date(2023, 11, 1, 0, 0, 0, 0, time.UTC)
	s, e = dates.BillingDates(current, billingday, months)

	if s.Format(dates.FormatYMD) != expectedS.Format(dates.FormatYMD) {
		t.Errorf("start date mis match expected [%s] actual [%s]", expectedS.String(), s.String())

	}
	if e.Format(dates.FormatYMD) != expectedE.Format(dates.FormatYMD) {
		t.Errorf("end date mis match expected [%s] actual [%s]", expectedE.String(), e.String())

	}

	billingday = 15
	months = 11
	current = time.Date(2024, 8, 21, 0, 0, 0, 0, time.UTC)
	expectedE = time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC)
	_, e = dates.BillingDates(current, billingday, months)

	if e.Format(dates.FormatYMD) != expectedE.Format(dates.FormatYMD) {
		t.Errorf("end date mis match expected [%s] actual [%s]", expectedE.String(), e.String())
	}
}

func TestSharedDatesResetMonth(t *testing.T) {
	ts := time.Date(2023, 10, 4, 23, 1, 0, 0, time.UTC)
	r := dates.ResetMonth(ts)

	if r.Format(dates.FormatYM) != ts.Format(dates.FormatYM) {
		t.Errorf("date was updated incorrectly")
	}

	if r.Day() != 1 {
		t.Errorf("day was not reset")
	}

}

func TestSharedDatesResetDay(t *testing.T) {
	ts := time.Date(2023, 10, 4, 23, 1, 0, 0, time.UTC)
	r := dates.ResetDay(ts)

	if r.Format(dates.FormatYMD) != ts.Format(dates.FormatYMD) {
		t.Errorf("date was updated incorrectly")
	}

	if r.Hour() != 0 {
		t.Errorf("hour was not reset")
	}

}

func TestSharedDatesRangeMonth(t *testing.T) {
	var (
		start time.Time = time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)
		end   time.Time = time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC)
	)

	rnge := dates.Range(start, end, dates.MONTH)
	if len(rnge) != 4 {
		t.Errorf("range returned incorrect amount")
	}
}

func TestSharedDatesRangeDay(t *testing.T) {
	var (
		start time.Time = time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC)
		end   time.Time = time.Date(2024, 3, 4, 0, 0, 0, 0, time.UTC)
	)

	rnge := dates.Range(start, end, dates.DAY)
	if len(rnge) != 5 {
		t.Errorf("range returned incorrect amount")
	}
}

func TestSharedDatesStrings(t *testing.T) {
	var (
		start time.Time = time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)
		end   time.Time = time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC)
	)
	rnge := dates.Range(start, end, dates.MONTH)
	strs := dates.Strings(rnge, dates.FormatYM)

	if len(strs) != len(strs) {
		t.Errorf("strings did not return same length as range")
	}

}

func TestSharedDatesMaxTime(t *testing.T) {
	min, max, _ := testhelpers.Dates()
	times := []time.Time{}
	for i := 0; i < 10; i++ {
		times = append(times, fake.Date(min, max))
	}
	big := max.AddDate(2, 0, 1)
	times = append(times, big)

	maxT := dates.MaxTime(times)

	if maxT != big {
		t.Errorf("max time did not find the largest")
	}

}
