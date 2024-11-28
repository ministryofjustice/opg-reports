package inputs

import (
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/ministryofjustice/opg-reports/internal/dateformats"
	"github.com/ministryofjustice/opg-reports/internal/dateintervals"
	"github.com/ministryofjustice/opg-reports/pkg/convert"
)

// VersionInput is the version part of the uri path
type VersionInput struct {
	Version string `json:"version,omitempty" path:"version" required:"true" doc:"Version prefix for the api" default:"v1" enum:"v1"`
}

// VersionUnitInput is the version part of the uri path and option to filter bu unit
type VersionUnitInput struct {
	Version string `json:"version,omitempty" path:"version" db:"-" required:"true" doc:"Version prefix for the api" default:"v1" enum:"v1"`
	Unit    string `json:"unit,omitempty" query:"unit" db:"unit" doc:"Unit name to filter data by"`
}

// // DateRangeInput describes input parameters to capture the start and end of a date range
// // The date range is >= start_date < end_date so in the case covering all of Jan 2024
// // you would want 2024-01-01 & 2024-02-01
// // This makes find last days easier for users
// type DateRangeInput struct {
// 	Version   string `json:"version" db:"-" path:"version" required:"true" doc:"Version prefix for the api" default:"v1" enum:"v1"`
// 	StartDate string `json:"start_date" db:"start_date" path:"start_date" required:"true" doc:"Earliest date to start the data (uses >=). YYYY-MM-DD." example:"2022-01-01" pattern:"([0-9]{4}-[0-9]{2}-[0-9]{2})"`
// 	EndDate   string `json:"end_date" db:"end_date" path:"end_date" required:"true" doc:"Latest date to capture the data for (uses <). YYYY-MM-DD."  example:"2024-04-01" pattern:"([0-9]{4}-[0-9]{2}-[0-9]{2})"`
// }

// OptionalDateRangeInput
type OptionalDateRangeInput struct {
	Version   string `json:"version,omitempty" db:"-" path:"version" required:"true" doc:"Version prefix for the api" default:"v1" enum:"v1"`
	Unit      string `json:"unit,omitempty" query:"unit" db:"unit" doc:"Unit name to filter data by"`
	StartDate string `json:"start_date,omitempty" db:"start_date" query:"start_date" doc:"Earliest date to start the data (uses >=). YYYY-MM-DD." example:"2022-01-01" pattern:"([0-9]{4}-[0-9]{2}-[0-9]{2})"`
	EndDate   string `json:"end_date,omitempty" db:"end_date" query:"end_date" doc:"Latest date to capture the data for (uses <). YYYY-MM-DD."  example:"2024-04-01" pattern:"([0-9]{4}-[0-9]{2}-[0-9]{2})"`
}

func (self *OptionalDateRangeInput) Start() (t time.Time) {
	if val, err := convert.ToTime(self.StartDate); err == nil {
		t = val
	}
	return
}
func (self *OptionalDateRangeInput) End() (t time.Time) {
	if val, err := convert.ToTime(self.EndDate); err == nil {
		t = val
	}
	return
}

// RequiredGroupedDateRangeInput must have start and end dates as well as interval details
type RequiredGroupedDateRangeInput struct {
	Version    string `json:"version,omitempty" db:"-" path:"version" required:"true" doc:"Version prefix for the api" default:"v1" enum:"v1"`
	StartDate  string `json:"start_date,omitempty" db:"start_date" path:"start_date" required:"true" doc:"Earliest date to start the data (uses >=). YYYY-MM-DD." example:"2022-01-01" pattern:"([0-9]{4}-[0-9]{2}-[0-9]{2})"`
	EndDate    string `json:"end_date,omitempty" db:"end_date" path:"end_date" required:"true" doc:"Latest date to capture the data for (uses <). YYYY-MM-DD."  example:"2024-04-01" pattern:"([0-9]{4}-[0-9]{2}-[0-9]{2})"`
	Interval   string `json:"interval,omitempty" db:"-" path:"interval" default:"month" enum:"year,month,day" doc:"Group the data by this type of interval."`
	DateFormat string `json:"date_format,omitempty" db:"date_format"`
}

func (self *RequiredGroupedDateRangeInput) Start() (t time.Time) {
	if val, err := convert.ToTime(self.StartDate); err == nil {
		t = val
	}
	return
}

func (self *RequiredGroupedDateRangeInput) End() (t time.Time) {
	if val, err := convert.ToTime(self.EndDate); err == nil {
		t = val
	}
	return
}

func (self *RequiredGroupedDateRangeInput) GetInterval() dateintervals.Interval {
	return dateintervals.Interval(self.Interval)
}

// Resolve to convert the interval input param into a date format for the db call
func (self *RequiredGroupedDateRangeInput) Resolve(ctx huma.Context) []error {
	self.DateFormat = string(dateformats.SqliteYM)

	if self.Interval == "year" {
		self.DateFormat = string(dateformats.SqliteY)
	} else if self.Interval == "day" {
		self.DateFormat = string(dateformats.SqliteYMD)
	}
	return nil
}

// RequiredGroupedDateRangeUnitInput must have start and end dates as well as interval details and has optional
type RequiredGroupedDateRangeUnitInput struct {
	Version    string `json:"version,omitempty" db:"-" path:"version" required:"true" doc:"Version prefix for the api" default:"v1" enum:"v1"`
	Unit       string `json:"unit,omitempty" query:"unit" db:"unit" doc:"Unit name to filter data by"`
	StartDate  string `json:"start_date,omitempty" db:"start_date" path:"start_date" required:"true" doc:"Earliest date to start the data (uses >=). YYYY-MM-DD." example:"2022-01-01" pattern:"([0-9]{4}-[0-9]{2}-[0-9]{2})"`
	EndDate    string `json:"end_date,omitempty" db:"end_date" path:"end_date" required:"true" doc:"Latest date to capture the data for (uses <). YYYY-MM-DD."  example:"2024-04-01" pattern:"([0-9]{4}-[0-9]{2}-[0-9]{2})"`
	Interval   string `json:"interval,omitempty" db:"-" path:"interval" default:"month" enum:"year,month,day" doc:"Group the data by this type of interval."`
	DateFormat string `json:"date_format,omitempty" db:"date_format"`
}

func (self *RequiredGroupedDateRangeUnitInput) Start() (t time.Time) {
	if val, err := convert.ToTime(self.StartDate); err == nil {
		t = val
	}
	return
}

func (self *RequiredGroupedDateRangeUnitInput) End() (t time.Time) {
	if val, err := convert.ToTime(self.EndDate); err == nil {
		t = val
	}
	return
}

func (self *RequiredGroupedDateRangeUnitInput) GetInterval() dateintervals.Interval {
	return dateintervals.Interval(self.Interval)
}

// Resolve to convert the interval input param into a date format for the db call
func (self *RequiredGroupedDateRangeUnitInput) Resolve(ctx huma.Context) []error {
	self.DateFormat = string(dateformats.SqliteYM)

	if self.Interval == "year" {
		self.DateFormat = string(dateformats.SqliteY)
	} else if self.Interval == "day" {
		self.DateFormat = string(dateformats.SqliteYMD)
	}
	return nil
}
