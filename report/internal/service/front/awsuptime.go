package front

import (
	"net/http"
	"opg-reports/report/internal/endpoints"
	"opg-reports/report/internal/repository/restr"
	"opg-reports/report/internal/service/front/datatable"
	"opg-reports/report/internal/utils"
	"time"
)

// apiResponseAwsUptimeGrouped represents the api data structure returned
// by the aws costs grouped end points
//
// endpoint: `/v1/awsuptime/grouped/{granularity}/{start_date}/{end_date}`
// interface: datatable.ResponseBody
type apiResponseAwsUptimeGrouped struct {
	Count  int                 `json:"count,omityempty"`
	Dates  []string            `json:"dates,omitempty"`
	Groups []string            `json:"groups,omitempty"`
	Data   []map[string]string `json:"data"`
}

// check interface
var _ datatable.ResponseBody = &apiResponseAwsUptimeGrouped{}

// DataHeaders returns the column headings used for the core data - generally Dates
func (self *apiResponseAwsUptimeGrouped) DataHeaders() (dh []string) {
	return self.Dates
}
func (self *apiResponseAwsUptimeGrouped) DataRows() (data []map[string]string) {
	return self.Data
}
func (self *apiResponseAwsUptimeGrouped) PaddedDataRows() (all []map[string]string) {
	padding := utils.DummyRows(self.Dates, "date")
	all = append(self.Data, padding...)
	return
}
func (self *apiResponseAwsUptimeGrouped) Identifiers() (identifiers []string) {
	return self.Groups
}
func (self *apiResponseAwsUptimeGrouped) Cells() (cells []string) {
	return self.Dates
}
func (self *apiResponseAwsUptimeGrouped) TransformColumn() string {
	return "date"
}
func (self *apiResponseAwsUptimeGrouped) ValueColumn() string {
	return "average"
}
func (self *apiResponseAwsUptimeGrouped) RowTotalKeyName() string {
	return "total"
}
func (self *apiResponseAwsUptimeGrouped) TrendKeyName() string {
	return "trend"
}
func (self *apiResponseAwsUptimeGrouped) SumColumns() (cols []string) {
	cols = self.Dates
	if c := self.RowTotalKeyName(); c != "" {
		cols = append(cols, c)
	}
	return
}
func (self *apiResponseAwsUptimeGrouped) RowTotalCleanup() datatable.RowTotalCleaner {
	return datatable.RowTotalsAveraged
}
func (self *apiResponseAwsUptimeGrouped) ColumnTotalCleanup() datatable.ColumnTotalCleaner {
	return datatable.ColumnTotalsAveraged
}

type awsUptimePreCallF func(params map[string]string)

// GetAwsUptimeGrouped call the api (based on config values), convert that data into tabluar
// form for rendering, ensuring there are no empty / missing columns.
//
// The request object is used to merge the query string parameters from the front end with the default ones
// used to call the api, therefore allowing the front end to directly change start_dates, team names etc
//
// the optional set of `adjusters` functions allows a different way to overwrite parameters, by
// running a function against the parameters before the end point is generated - this allows user
// to adjust a value with out knowing its original value.
//
// # Examples
//
// Grouping by team:
//
//	service.GetAwsUptimeGrouped(client, request, map[string]string{"team":"true"} )
func (self *Service) GetAwsUptimeGrouped(client restr.RepositoryRestGetter, request *http.Request, apiParameters map[string]string, adjusters ...awsUptimePreCallF) (table *datatable.DataTable, err error) {
	var (
		log      = self.log.With("operation", "GetAwsUptimeGrouped")
		defaults = awsUptimeParams()
		params   = mergeRequestWithMaps(request, defaults, apiParameters)
		endpoint string
	)
	// allow function to overwrite parameters
	for _, adjustF := range adjusters {
		adjustF(params)
	}
	// generate the full url to call
	endpoint = endpoints.Parse(endpoints.AWSUPTIME_GROUPED, params)

	log.With("defaults", defaults, "api", apiParameters, "merged", params).Debug("calling api for grouped uptime data")
	// call the api and pass along the result handler function to create a datatable
	table, err = getFromAPI(self.ctx, self.log, self.conf,
		client,
		endpoint,
		parseAwsUptimeGroupedF,
	)

	log.Debug("returning data table ... ")
	return
}

// parseAwsUptimeGroupedF calls the datatable.New which handles all of the conversion from a response
// into a datatable
//
// The conversion done by datatable.New is messy as it groups into table rows
func parseAwsUptimeGroupedF(response *apiResponseAwsUptimeGrouped) (dt *datatable.DataTable, err error) {
	dt, err = datatable.New(response)
	return
}

// awsUptimeParams returns a map of the values that aws costs endpoints can accept.
// See `awscosts.AwsUptimeGroupedInput` for the input struct.
func awsUptimeParams() map[string]string {
	var (
		now   = time.Now().UTC()
		end   = utils.TimeReset(now, utils.TimeIntervalMonth) // this month
		start = utils.MonthsAgo(-5)
	)
	return map[string]string{
		"start_date":  start.Format(utils.DATE_FORMATS.YM),
		"end_date":    end.Format(utils.DATE_FORMATS.YM),
		"granularity": string(utils.GranularityMonth),
		"team":        "true",
	}
}
