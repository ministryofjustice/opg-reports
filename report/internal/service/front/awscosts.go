package front

import (
	"net/http"
	"opg-reports/report/internal/endpoints"
	"opg-reports/report/internal/repository/restr"
	"opg-reports/report/internal/service/front/datatable"
	"opg-reports/report/internal/utils"
	"time"
)

// apiResponseAwsCostsGrouped represents the api data structure returned
// by the aws costs grouped end points
//
// endpoint: `/v1/awscosts/grouped/{granularity}/{start_date}/{end_date}`
// interface: datatable.ResponseBody
type apiResponseAwsCostsGrouped struct {
	Count  int                 `json:"count,omityempty"`
	Dates  []string            `json:"dates,omitempty"`
	Groups []string            `json:"groups,omitempty"`
	Data   []map[string]string `json:"data"`
}

// check interface
var _ datatable.ResponseBody = &apiResponseAwsCostsGrouped{}

// DataHeaders returns the column headings used for the core data - generally Dates
func (self *apiResponseAwsCostsGrouped) DataHeaders() (dh []string) {
	return self.Dates
}
func (self *apiResponseAwsCostsGrouped) DataRows() (data []map[string]string) {
	return self.Data
}
func (self *apiResponseAwsCostsGrouped) PaddedDataRows() (all []map[string]string) {
	padding := utils.DummyRows(self.Dates, "date")
	all = append(self.Data, padding...)
	return
}
func (self *apiResponseAwsCostsGrouped) Identifiers() (identifiers []string) {
	return self.Groups
}
func (self *apiResponseAwsCostsGrouped) Cells() (cells []string) {
	return self.Dates
}
func (self *apiResponseAwsCostsGrouped) TransformColumn() string {
	return "date"
}
func (self *apiResponseAwsCostsGrouped) ValueColumn() string {
	return "cost"
}
func (self *apiResponseAwsCostsGrouped) RowTotalKeyName() string {
	return "total"
}
func (self *apiResponseAwsCostsGrouped) TrendKeyName() string {
	return "trend"
}
func (self *apiResponseAwsCostsGrouped) SumColumns() (cols []string) {
	cols = self.Dates
	if c := self.RowTotalKeyName(); c != "" {
		cols = append(cols, c)
	}
	return
}

type preCallF func(params map[string]string)

// GetAwsCostsGrouped call the api (based on config values), convert that data into tabluar
// form for rendering, ensuring there are no empty / missing columns.
//
// The request object is used to merge the query string parameters from the front end with the default ones
// used to call the api, therefore allowing the front end to directly change start_dates, team names etc
//
// the options set of `adjusters` functions allows a different way to overwrite parameters, by
// running a function against the parameters before the end point is generated - this allows user
// to adjust a value with out knowing its original value
func (self *Service) GetAwsCostsGrouped(client restr.RepositoryRestGetter, request *http.Request, apiParameters map[string]string, adjusters ...preCallF) (table *datatable.DataTable, err error) {
	var (
		log      = self.log.With("operation", "GetAwsCostsGrouped")
		defaults = awsCostsParams(self.conf.Aws.BillingDate)
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
	log.Debug("returning data table ... ")
	return
}

// parseAwsCostsGroupedF calls the datatable.New which handles all of the conversion from a response
// into a datatable
//
// The conversion is messy!
func parseAwsCostsGroupedF(response *apiResponseAwsCostsGrouped) (dt *datatable.DataTable, err error) {
	dt, err = datatable.New(response)
	return
}

// awsCostsParams returns a map of the values that aws costs endpoints can accept.
// See `awscosts.GroupedAwsCostsInput` for the input struct.
func awsCostsParams(billingDate int) map[string]string {
	var (
		now   = time.Now().UTC()
		end   = utils.BillingMonth(now, billingDate)
		start = time.Date(end.Year(), end.Month()-5, 1, 0, 0, 0, 0, time.UTC)
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
	}
}
