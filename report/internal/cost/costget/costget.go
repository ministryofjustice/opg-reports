package costget

import (
	"context"
	"net/http"
	"opg-reports/report/internal/cost/costapi/costapiteam"
	"opg-reports/report/internal/cost/costapi/costapiteamfilter"
	"opg-reports/report/package/overwrite"
	"opg-reports/report/package/rest"
	"opg-reports/report/package/tabulate"
	"opg-reports/report/package/times"
	"time"
)

var timeout = (2 * time.Second)

type tableHeaders struct {
	Labels []string
	Data   []string
	Extra  []string
	End    []string
}

type Data struct {
	Headers *tableHeaders
	Data    []map[string]interface{}
}

// ByMonthFilterByTeam returns the cost data grouped by month, filtered by team
func ByMonthFilterByTeam(ctx context.Context, apiHost string, current *http.Request, overwrites ...*rest.Param) (data *Data, err error) {
	var (
		apiResponse *costapiteamfilter.Response
		params      []*rest.Param = []*rest.Param{}
		endpoint    string        = costapiteamfilter.ENDPOINT
		end         time.Time     = times.Today()
		start       time.Time     = times.Add(times.ResetMonth(end), -5, times.MONTH)
	)
	data = &Data{}
	// set default start / end dates
	params = []*rest.Param{
		{Type: rest.PATH, Key: "date_start", Value: times.AsYMDString(start)},
		{Type: rest.PATH, Key: "date_end", Value: times.AsYMDString(end)},
	}
	// overwrite any params
	params = overwrite.Overwrite(params, overwrites...)
	// call the api
	apiResponse, _, err = rest.Get[*costapiteamfilter.Response](ctx, current, &rest.Request{
		Host:     apiHost,
		Endpoint: endpoint,
		Timeout:  timeout,
		Params:   params,
	})
	if err != nil {
		return
	}
	data = &Data{
		Headers: &tableHeaders{
			Labels: apiResponse.Headers[tabulate.KEY],
			Data:   apiResponse.Headers[tabulate.DATA],
			Extra:  apiResponse.Headers[tabulate.EXTRA],
			End:    apiResponse.Headers[tabulate.END],
		},
		Data: apiResponse.Data,
	}

	return
}

// ByMonthTeam returns the cost data per team
func ByMonthTeam(ctx context.Context, apiHost string, current *http.Request, overwrites ...*rest.Param) (data *Data, err error) {
	var (
		apiResponse *costapiteam.Response
		params      []*rest.Param = []*rest.Param{}
		endpoint    string        = costapiteam.ENDPOINT
		end         time.Time     = times.Today()
		start       time.Time     = times.Add(times.ResetMonth(end), -5, times.MONTH)
	)
	data = &Data{}
	// set default start / end dates
	params = []*rest.Param{
		{Type: rest.PATH, Key: "date_start", Value: times.AsYMDString(start)},
		{Type: rest.PATH, Key: "date_end", Value: times.AsYMDString(end)},
	}
	// overwrite any params
	params = overwrite.Overwrite(params, overwrites...)
	// call the api
	apiResponse, _, err = rest.Get[*costapiteam.Response](ctx, current, &rest.Request{
		Host:     apiHost,
		Endpoint: endpoint,
		Timeout:  timeout,
		Params:   params,
	})
	if err != nil {
		return
	}
	data = &Data{
		Headers: &tableHeaders{
			Labels: apiResponse.Headers[tabulate.KEY],
			Data:   apiResponse.Headers[tabulate.DATA],
			Extra:  apiResponse.Headers[tabulate.EXTRA],
			End:    apiResponse.Headers[tabulate.END],
		},
		Data: apiResponse.Data,
	}

	return
}
