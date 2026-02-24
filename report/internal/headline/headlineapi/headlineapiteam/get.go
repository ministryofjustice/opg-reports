package headlineapiteam

import (
	"context"
	"net/http"
	"opg-reports/report/package/overwrite"
	"opg-reports/report/package/rest"
	"opg-reports/report/package/times"
	"time"
)

func Get(ctx context.Context, apiHost string, current *http.Request, overwrites ...*rest.Param) (data *Response, err error) {

	var (
		response *Response
		params   []*rest.Param = []*rest.Param{}
		timeout  time.Duration = (2 * time.Second)
		endpoint string        = ENDPOINT
		now      time.Time     = times.Today()
		end      time.Time     = times.Add(times.ResetMonth(now), -1, times.MONTH)
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
