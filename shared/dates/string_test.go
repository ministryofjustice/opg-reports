package dates_test

import (
	"testing"

	"github.com/ministryofjustice/opg-reports/shared/dates"
)

func TestSharedDatesTimeError(t *testing.T) {
	var fail = "2026 was a great yeaar"
	ts := dates.Time(fail)

	if ts != dates.ErrorTime {
		t.Errorf("should be match to failed time")
	}

}

func TestSharedDatesTimeWorking(t *testing.T) {
	var working = "2024-02-29"
	ts := dates.Time(working)

	if ts.Format(dates.FormatYMD) != working {
		t.Errorf("failed to convert time")
	}
}
