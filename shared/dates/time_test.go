package dates_test

import (
	"testing"
	"time"

	"github.com/ministryofjustice/opg-reports/shared/dates"
)

func TestSharedDatesBillingDates(t *testing.T) {

	// billing date of the 5th
	// current date of the 4th
	// ask for 3 months
	var (
		billingday int       = 5
		current    time.Time = time.Date(2023, 10, 4, 0, 0, 0, 0, time.UTC)
		months     int       = 3
		expectedS  time.Time = time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)
		expectedE  time.Time = time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC)
	)
	s, e := dates.BillingDates(current, billingday, months)

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
	current = time.Date(2023, 11, 28, 0, 0, 0, 0, time.UTC)
	months = 6
	expectedS = time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC)
	expectedE = time.Date(2023, 10, 31, 0, 0, 0, 0, time.UTC)
	s, e = dates.BillingDates(current, billingday, months)

	if s.Format(dates.FormatYMD) != expectedS.Format(dates.FormatYMD) {
		t.Errorf("start date mis match expected [%s] actual [%s]", expectedS.String(), s.String())

	}
	if e.Format(dates.FormatYMD) != expectedE.Format(dates.FormatYMD) {
		t.Errorf("end date mis match expected [%s] actual [%s]", expectedE.String(), e.String())

	}

}
