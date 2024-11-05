package uptimeio

import (
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/ministryofjustice/opg-reports/pkg/convert"
	"github.com/ministryofjustice/opg-reports/pkg/datastore"
	"github.com/ministryofjustice/opg-reports/sources/uptime"
)

// UptimeInput is the version part of the uri path
type UptimeInput struct {
	Version   string `json:"version" path:"version" required:"true" doc:"Version prefix for the api" default:"v1" enum:"v1"`
	StartDate string `json:"start_date" db:"start_date" path:"start_date" required:"true" doc:"Earliest date to start the data (uses >=). YYYY-MM-DD." example:"2022-01-01" pattern:"([0-9]{4}-[0-9]{2}-[0-9]{2})"`
	EndDate   string `json:"end_date" db:"end_date" path:"end_date" required:"true" doc:"Latest date to capture the data for (uses <). YYYY-MM-DD."  example:"2024-04-01" pattern:"([0-9]{4}-[0-9]{2}-[0-9]{2})"`
	Interval  string `json:"interval" db:"-" path:"interval" default:"month" enum:"year,month,day" doc:"Group the data by this type of interval."`
	Unit      string `json:"unit" db:"unit" query:"unit" doc:"optional unit name to filter data by"`

	DateFormat string `json:"date_format,omitempty" db:"date_format"`
}

func (self *UptimeInput) StartTime() (t time.Time) {
	if val, err := convert.ToTime(self.StartDate); err == nil {
		t = val
	}
	return
}
func (self *UptimeInput) EndTime() (t time.Time) {
	if val, err := convert.ToTime(self.EndDate); err == nil {
		t = val
	}
	return
}

// Resolve to convert the interval input param into a date format for the db call
func (self *UptimeInput) Resolve(ctx huma.Context) []error {
	self.DateFormat = datastore.Sqlite.YearMonthFormat
	if self.Interval == "year" {
		self.DateFormat = datastore.Sqlite.YearFormat
	} else if self.Interval == "day" {
		self.DateFormat = datastore.Sqlite.YearMonthDayFormat
	}
	return nil
}

// UptimeBody is the normal body for a resposne
type UptimeBody struct {
	Type         string                            `json:"type" doc:"States what type of data this is for front end handling"`
	Request      *UptimeInput                      `json:"request" doc:"The public parameters originaly specified in the request to this API."`
	Result       []*uptime.Uptime                  `json:"result" doc:"List of all results found matching query."`
	ColumnOrder  []string                          `json:"column_order" doc:"List of columns set in the order they should be rendered for each row."`
	ColumnValues map[string][]interface{}          `json:"column_values" doc:"Contains all of the ordered columns possible values, to help display rendering."`
	DateRange    []string                          `json:"date_range" doc:"list of string dates between the start and end date"`
	TableRows    map[string]map[string]interface{} `json:"-"`
}

// UptimeOutput is the standard response
type UptimeOutput struct {
	Body *UptimeBody
}
