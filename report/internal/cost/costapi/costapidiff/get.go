package costapidiff

import (
	"context"
	"net/http"
	"opg-reports/report/package/overwrite"
	"opg-reports/report/package/rest"
	"opg-reports/report/package/times"
	"time"
)

// Get call the api and get the response
func Get(ctx context.Context, apiHost string, current *http.Request, overwrites ...*rest.Param) (data *Response, err error) {
	var (
		response *Response
		timeout  time.Duration = (2 * time.Second)
		params   []*rest.Param = []*rest.Param{}
		endpoint string        = ENDPOINT
		now      time.Time     = times.Today()
		end      time.Time     = times.Add(times.ResetMonth(now), -1, times.MONTH)
		start    time.Time     = times.Add(times.ResetMonth(now), -2, times.MONTH)
	)
	// set default start / end dates
	params = []*rest.Param{
		{Type: rest.PATH, Key: "date_a", Value: times.AsYMString(start)},
		{Type: rest.PATH, Key: "date_b", Value: times.AsYMString(end)},
		{Type: rest.QUERY, Key: "change", Value: "100"},
	}
	// overwrite any params
	params = overwrite.Overwrite(params, overwrites...)
	// call the api
	response, _, err = rest.Get[*Response](ctx, current, &rest.Request{
		Host:     apiHost,
		Endpoint: endpoint,
		Timeout:  timeout,
		Params:   params,
	})
	if err != nil {
		return
	}
	data = response
	return
}
