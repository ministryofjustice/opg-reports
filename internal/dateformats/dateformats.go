// Package dateformats provices a custom type (wrap of string) and
// commonly used date formats as constants
package dateformats

import "time"

type Format string

const (
	Full         string = time.RFC3339
	Year         string = "2006"
	YearMonth    string = "2006-01"
	YearMonthDay string = "2006-01-02"
	Old          string = "2006-01-02 15:04:05.999999 +0000 UTC"
)

const (
	SqliteY   Format = "%Y"
	SqliteYM  Format = "%Y-%m"
	SqliteYMD Format = "%Y-%m-%d"
)
