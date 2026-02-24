package teamapiall

import (
	"context"
	"net/http"
	"opg-reports/report/package/overwrite"
	"opg-reports/report/package/rest"
	"time"
)

// Get calls the apiHost with known enpoint and returns data formatted for
// the team navigation in the front end
func Get(ctx context.Context, apiHost string, current *http.Request, overwrites ...*rest.Param) (teams []string, err error) {
	var (
		apiResponse *Response
		timeout     time.Duration = (2 * time.Second)
		params      []*rest.Param = []*rest.Param{}
		endpoint    string        = ENDPOINT
	)
	teams = []string{}
	// overwrite any params
	params = overwrite.Overwrite(params, overwrites...)
	// call the api
	apiResponse, _, err = rest.Get[*Response](ctx, current, &rest.Request{
		Host:     apiHost,
		Endpoint: endpoint,
		Timeout:  timeout,
		Params:   params,
	})

	if err != nil {
		return
	}
	teams = apiResponse.Data

	return
}
