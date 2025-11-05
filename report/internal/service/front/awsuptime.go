package front

import (
	"net/http"
	"opg-reports/report/internal/endpoints"
	"opg-reports/report/internal/repository/restr"
	"opg-reports/report/internal/service/front/datatable"
	"opg-reports/report/internal/utils"
	"time"
)

// awsUptimeHeaders used for the headings on the teables
type awsUptimeHeaders struct {
	Columns []string `json:"columns"` // the row headings - used in the table body and table header & footer at the start
	Data    []string `json:"data"`    // the data headings - should be the date ranges in order
	Extras  []string `json:"extras"`  // extra headers at the end of the table (trend & total)
}

// awsUptimeTable contains the uptime data, but formatted as a table for easier front end handling (less work per request)
type awsUptimeTable struct {
	Headers *awsUptimeHeaders   `json:"headers"`
	Rows    []map[string]string `json:"rows"`
	Footer  map[string]string   `json:"footer"`
}

// apiResponseAwsUptimeGrouped represents the api data structure returned
// by the aws costs grouped end points
//
// endpoint: `/v1/awsuptime/grouped/{granularity}/{start_date}/{end_date}`
// interface: datatable.ResponseBody
type apiResponseAwsUptimeGrouped struct {
	Count   int                 `json:"count,omityempty"`
	Dates   []string            `json:"dates,omitempty"`
	Groups  []string            `json:"groups,omitempty"`
	Data    []map[string]string `json:"data"`
	Tabular *awsCostsTable      `json:"tabular"`
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

// parseAwsUptimeGroupedF
//
// The conversion done by datatable.New is messy as it groups into table rows
func parseAwsUptimeGroupedF(response *apiResponseAwsUptimeGrouped) (dt *datatable.DataTable, err error) {
	dt = &datatable.DataTable{
		Body:         response.Tabular.Rows,
		RowHeaders:   response.Tabular.Headers.Columns,
		DataHeaders:  response.Tabular.Headers.Data,
		ExtraHeaders: response.Tabular.Headers.Extras,
		Footer:       response.Tabular.Footer,
		Others:       map[string]interface{}{},
	}
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
