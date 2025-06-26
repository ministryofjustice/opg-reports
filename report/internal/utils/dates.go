package utils

import "time"

type dateFormats struct {
	Full string
	Y    string
	YM   string
	YMD  string
}

var DATE_FORMATS = &dateFormats{
	Full: time.RFC3339,
	Y:    "2006",
	YM:   "2006-01",
	YMD:  "2006-01-02",
}
