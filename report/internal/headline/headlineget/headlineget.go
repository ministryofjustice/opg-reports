package headlineget

import (
	"context"
	"net/http"
	"opg-reports/report/internal/headline/headlineapi/headlineapihome"
	"opg-reports/report/package/overwrite"
	"opg-reports/report/package/rest"
	"opg-reports/report/package/times"
	"time"
)

var timeout = (2 * time.Second)

func ForHomepage(ctx context.Context, apiHost string, current *http.Request, overwrites ...*rest.Param) (data *headlineapihome.Response, err error) {

	var (
		response *headlineapihome.Response
		params   []*rest.Param = []*rest.Param{}
		endpoint string        = headlineapihome.ENDPOINT
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
	response, _, err = rest.Get[*headlineapihome.Response](ctx, current, &rest.Request{
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
