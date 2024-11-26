package dateutils

import (
	"log/slog"
	"time"

	"github.com/ministryofjustice/opg-reports/internal/dateformats"
)

// Format will return the format to use for the date string passed, using
// time.RFC3339 as base.
//
// Passing 2024 would return 2006, passing 2024-12-01 would return 2006-01-02
// and so on
func Format(value string) string {
	var layout = dateformats.Full
	max := len(layout)
	l := len(value)
	if l > max {
		return layout
	}
	f := layout[:l]
	slog.Debug("[dates] Format", slog.String("format", f))
	return f
}

func Reformat(s string, layout string) (date string) {
	if t, err := time.Parse(Format(s), s); err == nil {
		date = t.Format(layout)
	}
	return
}
