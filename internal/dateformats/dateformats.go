// Package dateformats provices a custom type (wrap of string) and
// commonly used date formats as constants
package dateformats

import "time"

const (
	Full   string = time.RFC3339
	Y      string = "2006"
	YM     string = "2006-01"
	YMD    string = "2006-01-02"
	NameYM string = "January 2006"
)

// old formats used in earlier versions of data
const (
	Old string = "2006-01-02 15:04:05.999999 +0000 UTC"
)

// sqlite date formats
const (
	SqliteY   string = "%Y"
	SqliteYM  string = "%Y-%m"
	SqliteYMD string = "%Y-%m-%d"
)