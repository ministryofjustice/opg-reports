package aws_costs

import (
	"net/http"
	"time"

	"github.com/ministryofjustice/opg-reports/servers/shared/query"
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
	GroupByD GroupBy
	GroupByT GroupBy

	Unit string
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
	a.GroupByT = GroupBy(a.GroupBy)
	// -- filters
	values = a.unit.Values(r)
	a.Unit = must.First(values)
}

func NewRequest(start time.Time, end time.Time, interval dates.Interval, group GroupBy) *ApiRequest {
	return &ApiRequest{
		start:    query.Get("start"),
		end:      query.Get("end"),
		interval: query.Get("interval"),
		groupBy:  query.Get("group"),
		unit:     query.Get("unit"),
		// -- defaults
		StartD:    start,
		EndD:      end,
		IntervalD: interval,
		GroupByD:  group,
	}
}
