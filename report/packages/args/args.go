// Package args provides the standard arguments that are used by
// multiple packages
//
// Provides consistent structs used by all handlers and the `Default`
// func to generate its standard values tha can then be overridden
// via environment vars of cobra config values.
package args

import (
	"fmt"
	"opg-reports/report/packages/reset"
	"os"
	"path/filepath"
	"time"
)

// API is used for api registeration arguments
type API struct {
	DB       *DB       // Database details
	Versions *Versions // Version data
	Hosts    *Hosts    // host data
}

// Import contains all the values that any import command
// needs to operate in hierarchy structure.
//
// Generally attached to a cobra command flags
type Import struct {
	DB      *DB      // Database details
	Filters *Filters // Data range and other filters
	Aws     *AWS     // AWS related settings
	Github  *GitHub  // Github related settings
	File    *File    // Raw data file
}

// GitHub settings used to fetch data from the api
type GitHub struct {
	Organisation string `json:"org"`    // github org (--org)
	Parent       string `json:"parent"` // github parent team (--parent)
}

// AWS data used to fetch from the api
type AWS struct {
	AccountID string `json:"id"`     // set from the client
	Region    string `json:"region"` // AWS region (--region)
}

// Dates captures date inputs
type Dates struct {
	Start      time.Time `json:"date_start"`       // start date; normally first day of the current month (--date-start)
	StartCosts time.Time `json:"date_start_costs"` // start date for costs; this should be 2 months ago to get stable cost data (--date-start-costs)
	End        time.Time `json:"date_end"`         // end date; normally start of current day (--date-end)
}

// Filters are generic filters that can be used between any
type Filters struct {
	Dates  *Dates
	Filter string `json:"filter"` // Filtering with an command (--filter)
}

type File struct {
	Path string `json:"file_path"` // raw file path location (--src-file)
}

type Versions struct {
	Version string `json:"version"` // (--version)
	SHA     string `json:"sha"`     // (--sha)
}

type Hosts struct {
	Front string `json:"front_host"` // (--front-host)
	API   string `json:"api_host"`   // (--api-host)
}

// DB contains the db connection details
type DB struct {
	Driver string `json:"driver"` // DB driver - generally sqlite (--driver)
	DB     string `json:"db"`     // DB location (--db)
	Params string `json:"params"` // DB connection parameters (--params)
}

func (self *DB) Connection() (driver string, conn string) {
	os.MkdirAll(filepath.Dir(self.DB), os.ModePerm)

	if self.Driver == "sqlite3" {
		conn = fmt.Sprintf("%s%s", self.DB, self.Params)
	}
	driver = self.Driver
	return
}

// argType used as constraints
type argType interface {
	*Import |
		*API |
		*DB |
		*Dates |
		*Filters |
		*AWS |
		*GitHub |
		*File |
		*Versions |
		*Hosts
}

// Default returns T with populated default values that typically used
// with the type.
//
// For *Dates - if its the first of the month, reset to use the first
// day of the previous month as the start to ensure its fully captured.
//
// *Import and similar will recursively call its sub structures.
func Default[T argType](now time.Time) T {
	var def interface{}
	var arg T

	switch any(arg).(type) {
	case *API:
		def = &API{
			DB:       Default[*DB](now),
			Versions: Default[*Versions](now),
			Hosts:    Default[*Hosts](now),
		}
	case *Import:
		def = &Import{
			DB:      Default[*DB](now),
			Filters: Default[*Filters](now),
			Aws:     Default[*AWS](now),
			Github:  Default[*GitHub](now),
			File:    Default[*File](now),
		}
	case *Hosts:
		def = &Hosts{
			Front: ":8080",
			API:   ":8081",
		}
	case *Versions:
		def = &Versions{
			Version: "0.0.1",
			SHA:     "abcdef",
		}
	case *File:
		def = &File{
			Path: "",
		}
	case *DB:
		def = &DB{
			Driver: "sqlite3",
			DB:     "./database/api.db",
			Params: "?_journal=WAL&_busy_timeout=5000&_vacuum=incremental&_synchronous=NORMAL&_cache_size=1000000000",
		}

	case *Dates:
		var start = reset.Month(&now)
		// if its the first day of the month, start date
		// goes back one day to get data for that month
		// as well
		if now.Day() == 1 {
			start = start.AddDate(0, -1, 0)
		}
		def = &Dates{
			End:        reset.Day(&now),
			Start:      start,
			StartCosts: start.AddDate(0, -2, 0),
		}
	case *Filters:
		def = &Filters{
			Filter: "",
			Dates:  Default[*Dates](now),
		}
	case *AWS:
		def = &AWS{
			Region: "eu-west-1",
		}
	case *GitHub:
		def = &GitHub{
			Organisation: "ministryofjustice",
			Parent:       "opg",
		}
	}
	return def.(T)
}
