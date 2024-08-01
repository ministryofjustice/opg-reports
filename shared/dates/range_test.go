package dates

import (
	"opg-reports/shared/logger"
	"testing"
	"time"
)

func TestSharedDatesMonths(t *testing.T) {
	logger.LogSetup()
	s := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	e := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	set := Months(s, e)
	if len(set) != 13 {
		t.Errorf("error with date range")
	}

	s = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	e = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	set = Months(s, e)
	if len(set) != 1 {
		t.Errorf("error with date range")
	}

	s = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	e = time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	set = Months(s, e)
	if len(set) != 3 {
		t.Errorf("error with date range")
	}
}

func TestSharedDatesDays(t *testing.T) {
	logger.LogSetup()
	s := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	e := time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC)
	set := Days(s, e)
	if len(set) != 29 {
		t.Errorf("error with date range: %d", len(set))
	}
}

func TestSharedDatesInMonth(t *testing.T) {
	logger.LogSetup()
	months := []string{"2024-01", "2024-02", "2024-03"}

	d := "2024-02"
	if InMonth(d, months) != true {
		t.Errorf("in month should be true")
	}

	d = "2023-12"
	if InMonth(d, months) == true {
		t.Errorf("in month should be false")
	}
}
