package costapidetailed

import (
	"context"
	"net/http"
	"opg-reports/report/package/overwrite"
	"opg-reports/report/package/rest"
	"opg-reports/report/package/times"
	"time"
)

// Get call the api and get the response, grouped costs by team
func Get(ctx context.Context, apiHost string, current *http.Request, overwrites ...*rest.Param) (data *Response, err error) {
	var (
		response *Response
		timeout  time.Duration = (2 * time.Second)
		params   []*rest.Param = []*rest.Param{}
		endpoint string        = ENDPOINT
		end      time.Time     = times.Today()
		start    time.Time     = times.Add(times.ResetMonth(end), -4, times.MONTH)
	)
	// set default start / end dates
	params = []*rest.Param{
		{Type: rest.PATH, Key: "date_start", Value: times.AsYMString(start)},
		{Type: rest.PATH, Key: "date_end", Value: times.AsYMString(end)},
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
