package inout

import (
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/ministryofjustice/opg-reports/internal/dateformats"
	"github.com/ministryofjustice/opg-reports/internal/dateintervals"
	"github.com/ministryofjustice/opg-reports/internal/dateutils"
)

// VersionInput is the version part of the uri path and has the following fields for the uri:
//
//   - version (path, required)
type VersionInput struct {
	Version string `json:"version,omitempty" path:"version" required:"true" doc:"Version prefix for the api" default:"v1" enum:"v1"`
}

// VersionUnitInput is the version part of the uri path and option to filter bu unit and has the following field setup:
//
//   - version (path, required)
//   - unit (query, optional)
type VersionUnitInput struct {
	Version string `json:"version,omitempty" path:"version" db:"-" required:"true" doc:"Version prefix for the api" default:"v1" enum:"v1"`
	Unit    string `json:"unit,omitempty" query:"unit" db:"unit" doc:"Unit name to filter data by"`
}

// DateRangeUnitInput allows for following fields:
//
//   - version (path, required)
//   - unit (query, optional)
//   - start_date (path, required)
//   - end_date (path, required)
type DateRangeUnitInput struct {
	Version   string `json:"version,omitempty" db:"-" path:"version" required:"true" doc:"Version prefix for the api" default:"v1" enum:"v1"`
	Unit      string `json:"unit,omitempty" query:"unit" db:"unit" doc:"Unit name to filter data by"`
	StartDate string `json:"start_date,omitempty" db:"start_date" required:"true" path:"start_date" doc:"Earliest date to start the data (uses >=). YYYY-MM-DD." example:"2022-01-01" pattern:"([0-9]{4}-[0-9]{2}-[0-9]{2})"`
	EndDate   string `json:"end_date,omitempty" db:"end_date" required:"true" path:"end_date" doc:"Latest date to capture the data for (uses <). YYYY-MM-DD."  example:"2024-04-01" pattern:"([0-9]{4}-[0-9]{2}-[0-9]{2})"`
}

func (self *DateRangeUnitInput) Start() (t time.Time) {
	if val, err := dateutils.Time(self.StartDate); err == nil {
		t = val
	}
	return
}
func (self *DateRangeUnitInput) End() (t time.Time) {
	if val, err := dateutils.Time(self.EndDate); err == nil {
		t = val
	}
	return
}

// OptionalDateRangeInput has the following fields for input and their setup:
//
//   - version (path, required)
//   - unit (query, optional)
//   - start_date (query, optional)
//   - end_date (query, optional)
type OptionalDateRangeInput struct {
	Version   string `json:"version,omitempty" db:"-" path:"version" required:"true" doc:"Version prefix for the api" default:"v1" enum:"v1"`
	Unit      string `json:"unit,omitempty" query:"unit" db:"unit" doc:"Unit name to filter data by"`
	StartDate string `json:"start_date,omitempty" db:"start_date" query:"start_date" doc:"Earliest date to start the data (uses >=). YYYY-MM-DD." example:"2022-01-01" pattern:"([0-9]{4}-[0-9]{2}-[0-9]{2})"`
	EndDate   string `json:"end_date,omitempty" db:"end_date" query:"end_date" doc:"Latest date to capture the data for (uses <). YYYY-MM-DD."  example:"2024-04-01" pattern:"([0-9]{4}-[0-9]{2}-[0-9]{2})"`
}

func (self *OptionalDateRangeInput) Start() (t time.Time) {
	if val, err := dateutils.Time(self.StartDate); err == nil {
		t = val
	}
	return
}
func (self *OptionalDateRangeInput) End() (t time.Time) {
	if val, err := dateutils.Time(self.EndDate); err == nil {
		t = val
	}
	return
}

// RequiredGroupedDateRangeInput must have start and end dates as well as interval details:
//
//   - version (path, required)
//   - start_date (path, required)
//   - end_date (path, required)
//   - interval (path, required)
//
// Note: no unit filtering for this.
type RequiredGroupedDateRangeInput struct {
	Version    string `json:"version,omitempty" db:"-" path:"version" required:"true" doc:"Version prefix for the api" default:"v1" enum:"v1"`
	StartDate  string `json:"start_date,omitempty" db:"start_date" path:"start_date" required:"true" doc:"Earliest date to start the data (uses >=). YYYY-MM-DD." example:"2022-01-01" pattern:"([0-9]{4}-[0-9]{2}-[0-9]{2})"`
	EndDate    string `json:"end_date,omitempty" db:"end_date" path:"end_date" required:"true" doc:"Latest date to capture the data for (uses <). YYYY-MM-DD."  example:"2024-04-01" pattern:"([0-9]{4}-[0-9]{2}-[0-9]{2})"`
	Interval   string `json:"interval,omitempty" db:"-" path:"interval" required:"true" default:"month" enum:"year,month,day" doc:"Group the data by this type of interval."`
	DateFormat string `json:"date_format,omitempty" db:"date_format"`
}

func (self *RequiredGroupedDateRangeInput) Start() (t time.Time) {
	if val, err := dateutils.Time(self.StartDate); err == nil {
		t = val
	}
	return
}

func (self *RequiredGroupedDateRangeInput) End() (t time.Time) {
	if val, err := dateutils.Time(self.EndDate); err == nil {
		t = val
	}
	return
}

func (self *RequiredGroupedDateRangeInput) GetInterval() dateintervals.Interval {
	return dateintervals.Interval(self.Interval)
}

// Resolve to convert the interval input param into a date format for the db call
func (self *RequiredGroupedDateRangeInput) Resolve(ctx huma.Context) []error {
	self.DateFormat = dateformats.SqliteYM

	if self.Interval == "year" {
		self.DateFormat = dateformats.SqliteY
	} else if self.Interval == "day" {
		self.DateFormat = dateformats.SqliteYMD
	}
	return nil
}

// RequiredGroupedDateRangeUnitInput must have start and end dates as well as interval details and has optional unit
// filtering. Setup:
//
//   - version (path, required)
//   - unit (query, optional)
//   - start_date (path, required)
//   - end_date (path, required)
//   - interval (path, required)
type RequiredGroupedDateRangeUnitInput struct {
	Version    string `json:"version,omitempty" db:"-" path:"version" required:"true" doc:"Version prefix for the api" default:"v1" enum:"v1"`
	Unit       string `json:"unit,omitempty" query:"unit" db:"unit" doc:"Unit name to filter data by"`
	StartDate  string `json:"start_date,omitempty" db:"start_date" path:"start_date" required:"true" doc:"Earliest date to start the data (uses >=). YYYY-MM-DD." example:"2022-01-01" pattern:"([0-9]{4}-[0-9]{2}-[0-9]{2})"`
	EndDate    string `json:"end_date,omitempty" db:"end_date" path:"end_date" required:"true" doc:"Latest date to capture the data for (uses <). YYYY-MM-DD."  example:"2024-04-01" pattern:"([0-9]{4}-[0-9]{2}-[0-9]{2})"`
	Interval   string `json:"interval,omitempty" db:"-" path:"interval" default:"month" enum:"year,month,day" doc:"Group the data by this type of interval."`
	DateFormat string `json:"date_format,omitempty" db:"date_format"`
}

func (self *RequiredGroupedDateRangeUnitInput) Start() (t time.Time) {
	if val, err := dateutils.Time(self.StartDate); err == nil {
		t = val
	}
	return
}

func (self *RequiredGroupedDateRangeUnitInput) End() (t time.Time) {
	if val, err := dateutils.Time(self.EndDate); err == nil {
		t = val
	}
	return
}

func (self *RequiredGroupedDateRangeUnitInput) GetInterval() dateintervals.Interval {
	return dateintervals.Interval(self.Interval)
}

// Resolve to convert the interval input param into a date format for the db call
func (self *RequiredGroupedDateRangeUnitInput) Resolve(ctx huma.Context) []error {
	self.DateFormat = dateformats.SqliteYM

	if self.Interval == "year" {
		self.DateFormat = dateformats.SqliteY
	} else if self.Interval == "day" {
		self.DateFormat = dateformats.SqliteYMD
	}
	return nil
}
