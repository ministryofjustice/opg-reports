// Package dateformats provices a custom type (wrap of string) and
// commonly used date formats as constants
package dateformats

type Format string

const (
	SqliteY   Format = "%Y"
	SqliteYM  Format = "%Y-%m"
	SqliteYMD Format = "%Y-%m-%d"
)
