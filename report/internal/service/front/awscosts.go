package front

import (
	"net/http"
	"opg-reports/report/internal/endpoints"
	"opg-reports/report/internal/repository/restr"
	"opg-reports/report/internal/service/front/datatable"
	"opg-reports/report/internal/utils"
	"time"
)

// awsCostHeaders used for the headings on the teables
type awsCostHeaders struct {
	Columns []string `json:"columns"` // the row headings - used in the table body and table header & footer at the start
	Data    []string `json:"data"`    // the data headings - should be the date ranges in order
	Extras  []string `json:"extras"`  // extra headers at the end of the table (trend & total)
}

// awsCostsTable contains the cost data, but formatted as a table for easier front end handling (less work per request)
type awsCostsTable struct {
	Headers *awsCostHeaders     `json:"headers"`
	Rows    []map[string]string `json:"rows"`
	Footer  map[string]string   `json:"footer"`
}

// apiResponseAwsCostsGrouped represents the api data structure returned
// by the aws costs grouped end points
//
// endpoint: `/v1/awscosts/grouped/{granularity}/{start_date}/{end_date}`
// interface: datatable.ResponseBody
type apiResponseAwsCostsGrouped struct {
	Count   int                 `json:"count,omityempty"`
	Dates   []string            `json:"dates,omitempty"`
	Groups  []string            `json:"groups,omitempty"`
	Data    []map[string]string `json:"data"`
	Tabular *awsCostsTable      `json:"tabular"`
}

type awsCostsPreCallF func(params map[string]string)

// GetAwsCostsGrouped call the api (based on config values), convert that data into tabluar
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
//	service.GetAwsCostsGrouped(client, request, map[string]string{"team":"true"} )
//
// Filtering costs by team & then grouping by account (and month by default):
//
//	service.GetAwsCostsGrouped(client, request, map[string]string{"team": teamName, "account_name": "true"} )
//
// Note: Always adds billing date to the Others property for use in front end returned data
func (self *Service) GetAwsCostsGrouped(client restr.RepositoryRestGetter, request *http.Request, apiParameters map[string]string, adjusters ...awsCostsPreCallF) (table *datatable.DataTable, err error) {
	var (
		log      = self.log.With("operation", "GetAwsCostsGrouped")
		defaults = awsCostsParams()
		params   = mergeRequestWithMaps(request, defaults, apiParameters)
		endpoint string
	)
	// allow function to overwrite parameters
	for _, adjustF := range adjusters {
		adjustF(params)
	}
	endpoint = endpoints.Parse(endpoints.AWSCOSTS_GROUPED, params)

	log.With("defaults", defaults, "api", apiParameters, "merged", params).Debug("calling api for grouped cost data")
	table, err = getFromAPI(self.ctx, self.log, self.conf,
		client,
		endpoint,
		parseAwsCostsGroupedF,
	)
	if err == nil {
		table.Others["BillingDay"] = self.conf.Aws.BillingDate
	}

	log.Debug("returning data table ... ")
	return
}

// parseAwsCostsGroupedF
//
// The conversion done by datatable.New is messy as it groups into table rows
func parseAwsCostsGroupedF(response *apiResponseAwsCostsGrouped) (dt *datatable.DataTable, err error) {
	dt = &datatable.DataTable{}

	if response.Tabular != nil {
		dt = &datatable.DataTable{
			Body:         response.Tabular.Rows,
			RowHeaders:   response.Tabular.Headers.Columns,
			DataHeaders:  response.Tabular.Headers.Data,
			ExtraHeaders: response.Tabular.Headers.Extras,
			Footer:       response.Tabular.Footer,
		}
	}

	dt.Others = map[string]interface{}{}
	return
}

// awsCostsParams returns a map of the values that aws costs endpoints can accept.
// See `awscosts.AwsCostsGroupedInput` for the input struct.
func awsCostsParams() map[string]string {
	var (
		now   = time.Now().UTC()
		end   = utils.TimeReset(now, utils.TimeIntervalMonth)
		start = utils.MonthsAgo(-5) // time.Date(end.Year(), end.Month()-5, 1, 0, 0, 0, 0, time.UTC)
	)
	return map[string]string{
		"start_date":  start.Format(utils.DATE_FORMATS.YM),
		"end_date":    end.Format(utils.DATE_FORMATS.YM),
		"granularity": string(utils.GranularityMonth),
		"team":        "true",
		"region":      "-",
		"service":     "-",
		"account":     "-",
		"environment": "-",
		"tabular":     "",
	}
}
