package releasesio

import (
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/ministryofjustice/opg-reports/pkg/convert"
	"github.com/ministryofjustice/opg-reports/pkg/datastore"
	"github.com/ministryofjustice/opg-reports/sources/releases"
)

type ReleasesTeamsInput struct {
	Version string `json:"version" path:"version" required:"true" doc:"Version prefix for the api" default:"v1" enum:"v1"`
}

type ReleasesListAllInput struct {
	Version   string `json:"version" path:"version" required:"true" doc:"Version prefix for the api" default:"v1" enum:"v1"`
	StartDate string `json:"start_date" db:"start_date" path:"start_date" required:"true" doc:"Earliest date to start the data (uses >=). YYYY-MM-DD." example:"2022-01-01" pattern:"([0-9]{4}-[0-9]{2}-[0-9]{2})"`
	EndDate   string `json:"end_date" db:"end_date" path:"end_date" required:"true" doc:"Latest date to capture the data for (uses <). YYYY-MM-DD."  example:"2024-04-01" pattern:"([0-9]{4}-[0-9]{2}-[0-9]{2})"`
	Unit      string `json:"unit" db:"unit" query:"unit" doc:"optional unit name to filter data by"`
}

type ReleasesInput struct {
	Version   string `json:"version" path:"version" required:"true" doc:"Version prefix for the api" default:"v1" enum:"v1"`
	StartDate string `json:"start_date" db:"start_date" path:"start_date" required:"true" doc:"Earliest date to start the data (uses >=). YYYY-MM-DD." example:"2022-01-01" pattern:"([0-9]{4}-[0-9]{2}-[0-9]{2})"`
	EndDate   string `json:"end_date" db:"end_date" path:"end_date" required:"true" doc:"Latest date to capture the data for (uses <). YYYY-MM-DD."  example:"2024-04-01" pattern:"([0-9]{4}-[0-9]{2}-[0-9]{2})"`
	Interval  string `json:"interval" db:"-" path:"interval" default:"month" enum:"year,month,day" doc:"Group the data by this type of interval."`
	Unit      string `json:"unit" db:"unit" query:"unit" doc:"optional unit name to filter data by"`

	DateFormat string `json:"date_format,omitempty" db:"date_format"`
}

func (self *ReleasesInput) StartTime() (t time.Time) {
	if val, err := convert.ToTime(self.StartDate); err == nil {
		t = val
	}
	return
}
func (self *ReleasesInput) EndTime() (t time.Time) {
	if val, err := convert.ToTime(self.EndDate); err == nil {
		t = val
	}
	return
}

// Resolve to convert the interval input param into a date format for the db call
func (self *ReleasesInput) Resolve(ctx huma.Context) []error {
	self.DateFormat = datastore.Sqlite.YearMonthFormat
	if self.Interval == "year" {
		self.DateFormat = datastore.Sqlite.YearFormat
	} else if self.Interval == "day" {
		self.DateFormat = datastore.Sqlite.YearMonthDayFormat
	}
	return nil
}

type ReleasesTeamsBody struct {
	Type      string                            `json:"type" doc:"States what type of data this is for front end handling"`
	Request   *ReleasesTeamsInput               `json:"request" doc:"The public parameters originaly specified in the request to this API."`
	Result    []*releases.Team                  `json:"result" doc:"List of all items found matching query."`
	TableRows map[string]map[string]interface{} `json:"-"`
}

type ReleaseTeamsOutput struct {
	Body *ReleasesTeamsBody
}

type ReleasesListAllBody struct {
	Type      string                            `json:"type" doc:"States what type of data this is for front end handling"`
	Request   *ReleasesListAllInput             `json:"request" doc:"The public parameters originaly specified in the request to this API."`
	Result    []*releases.Release               `json:"result" doc:"List of all items found matching query."`
	TableRows map[string]map[string]interface{} `json:"-"`
}

type ReleaseListAllOutput struct {
	Body *ReleasesListAllBody
}

type ReleasesBody struct {
	Type         string                            `json:"type" doc:"States what type of data this is for front end handling"`
	Request      *ReleasesInput                    `json:"request" doc:"The public parameters originaly specified in the request to this API."`
	Result       []*releases.Release               `json:"result" doc:"List of all items found matching query."`
	ColumnOrder  []string                          `json:"column_order" doc:"List of columns set in the order they should be rendered for each row."`
	ColumnValues map[string][]interface{}          `json:"column_values" doc:"Contains all of the ordered columns possible values, to help display rendering."`
	DateRange    []string                          `json:"date_range" doc:"list of string dates between the start and end date"`
	TableRows    map[string]map[string]interface{} `json:"-"`
}

type ReleaseOutput struct {
	Body *ReleasesBody
}
