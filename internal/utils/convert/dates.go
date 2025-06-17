package convert

import (
	"time"

	"github.com/ministryofjustice/opg-reports/internal/utils/dates"
)

// StringToTime converts string to a time.Time using time.Parse
// and guesses the format via LayoutToUse
func StringToTime(s string) (t time.Time, err error) {
	layout := dates.GuessDateFormat(s)
	t, err = time.Parse(layout, s)
	return
}
