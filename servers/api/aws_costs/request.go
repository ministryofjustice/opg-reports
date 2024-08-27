package aws_costs

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ministryofjustice/opg-reports/datastore/aws_costs/awsc"
	"github.com/ministryofjustice/opg-reports/servers/shared/query"
	"github.com/ministryofjustice/opg-reports/shared/consts"
	"github.com/ministryofjustice/opg-reports/shared/dates"
	"github.com/ministryofjustice/opg-reports/shared/must"
)

type ApiRequest struct {
	start    *query.Query
	end      *query.Query
	interval *query.Query
	groupBy  *query.Query
	unit     *query.Query

	Start  string
	StartD time.Time
	StartT time.Time

	End      string
	EndD     time.Time
	EndT     time.Time
	RangeEnd time.Time

	Interval       string
	IntervalD      dates.Interval
	IntervalT      dates.Interval
	IntervalFormat string

	GroupBy  string
	GroupByD consts.GroupBy
	GroupByT consts.GroupBy
}

func (a *ApiRequest) Update(r *http.Request) {
	var values []string
	// -- interval
	values = a.interval.Values(r)
	a.Interval = must.FirstOrDefault(values, string(a.IntervalD))
	a.IntervalT = dates.Interval(a.Interval)
	a.IntervalFormat = dates.IntervalFormat(a.IntervalT)
	// -- start
	values = a.start.Values(r)
	a.Start = must.FirstOrDefault(values, a.StartD.Format(a.IntervalFormat))
	a.StartT = dates.Time(a.Start)
	// -- end
	values = a.end.Values(r)
	a.End = must.FirstOrDefault(values, a.EndD.Format(a.IntervalFormat))
	a.EndT = dates.Time(a.End)
	// -- range
	a.RangeEnd = a.EndT.AddDate(0, -1, 0)
	if a.IntervalT == dates.DAY {
		a.RangeEnd = a.EndT.AddDate(0, 0, -1)
	}
	// -- groupby
	values = a.groupBy.Values(r)
	a.GroupBy = must.FirstOrDefault(values, string(a.GroupByD))
	a.GroupByT = consts.GroupBy(a.GroupBy)
}

type queryWrapperF func(ctx context.Context, req *ApiRequest, q *awsc.Queries) (results []*CommonResult, err error)

// --- all the functions that call the correct queries
func monthlyPerUnit(ctx context.Context, req *ApiRequest, q *awsc.Queries) (results []*CommonResult, err error) {
	params := awsc.MonthlyCostsPerUnitParams{Start: req.Start, End: req.End}
	if res, err := q.MonthlyCostsPerUnit(ctx, params); err == nil {
		results = Common(res)
	}
	return
}

func monthlyPerUnitEnv(ctx context.Context, req *ApiRequest, q *awsc.Queries) (results []*CommonResult, err error) {
	params := awsc.MonthlyCostsPerUnitEnvironmentParams{Start: req.Start, End: req.End}
	if res, err := q.MonthlyCostsPerUnitEnvironment(ctx, params); err == nil {
		results = Common(res)
	}
	return
}

func monthlyDetails(ctx context.Context, req *ApiRequest, q *awsc.Queries) (results []*CommonResult, err error) {
	params := awsc.MonthlyCostsDetailedParams{Start: req.Start, End: req.End}
	if res, err := q.MonthlyCostsDetailed(ctx, params); err == nil {
		results = Common(res)
	}
	return
}

func dailyPerUnit(ctx context.Context, req *ApiRequest, q *awsc.Queries) (results []*CommonResult, err error) {
	params := awsc.DailyCostsPerUnitParams{Start: req.Start, End: req.End}
	if res, err := q.DailyCostsPerUnit(ctx, params); err == nil {
		results = Common(res)
	}
	return
}

func dailyPerUnitEnv(ctx context.Context, req *ApiRequest, q *awsc.Queries) (results []*CommonResult, err error) {
	params := awsc.DailyCostsPerUnitEnvironmentParams{Start: req.Start, End: req.End}
	if res, err := q.DailyCostsPerUnitEnvironment(ctx, params); err == nil {
		results = Common(res)
	}
	return
}

func dailyDetails(ctx context.Context, req *ApiRequest, q *awsc.Queries) (results []*CommonResult, err error) {
	params := awsc.DailyCostsDetailedParams{Start: req.Start, End: req.End}
	if res, err := q.DailyCostsDetailed(ctx, params); err == nil {
		results = Common(res)
	}
	return
}

// setIfFound
func setIfFound(r *CommonResult, columns map[string]map[string]bool) {
	if r.Unit != "" {
		if _, ok := columns["unit"]; !ok {
			columns["unit"] = map[string]bool{}
		}
		columns["unit"][r.Unit] = true
	}
	if r.Environment != nil && r.Environment.(string) != "" {
		if _, ok := columns["environment"]; !ok {
			columns["environment"] = map[string]bool{}
		}
		columns["environment"][r.Environment.(string)] = true
	}
	if r.AccountID != "" {
		if _, ok := columns["account_id"]; !ok {
			columns["account_id"] = map[string]bool{}
		}
		columns["account_id"][r.AccountID] = true
	}
	if r.Service != "" {
		if _, ok := columns["service"]; !ok {
			columns["service"] = map[string]bool{}
		}
		columns["service"][r.Service] = true
	}
}

// StandardQueryResults uses the group and interval values from the request to determine
// which db query to run
//
// The query to use is determined by the following:
// Interval set to `MONTH`
//   - group set to `unit` (unit)
//   - group set to `unit-env` (unit and environment)
//   - group set to `detailed` (unit, environment, account id and service)
//
// Interval set to `DAY`
//   - group set to `unit` (unit)
//   - group set to `unit-env` (unit and environment)
//   - group set to `detailed` (unit, environment, account id and service)
//
// Query results are converted to `[]*CommonResult` struct.
func StandardQueryResults(ctx context.Context, q *awsc.Queries, req *ApiRequest) (results []*CommonResult, err error) {
	// define the possible functions and params
	var possibleFuncs = map[string]queryWrapperF{
		string(dates.MONTH) + string(consts.GroupByUnit):            monthlyPerUnit,
		string(dates.MONTH) + string(consts.GroupByUnitEnvironment): monthlyPerUnitEnv,
		string(dates.MONTH) + string(consts.GroupByDetailed):        monthlyDetails,
		//
		string(dates.DAY) + string(consts.GroupByUnit):            dailyPerUnit,
		string(dates.DAY) + string(consts.GroupByUnitEnvironment): dailyPerUnitEnv,
		string(dates.DAY) + string(consts.GroupByDetailed):        dailyDetails,
	}

	key := req.Interval + req.GroupBy

	if f, ok := possibleFuncs[key]; ok {
		results, err = f(ctx, req, q)
	} else {
		err = fmt.Errorf("error finding query function based on get paremters")
	}
	return
}

func NewRequest(start time.Time, end time.Time, interval dates.Interval, group consts.GroupBy) *ApiRequest {
	return &ApiRequest{
		start:    query.Get("start"),
		end:      query.Get("end"),
		interval: query.Get("interval"),
		groupBy:  query.Get("group"),
		// -- defaults
		StartD:    start,
		EndD:      end,
		IntervalD: interval,
		GroupByD:  group,
	}
}
