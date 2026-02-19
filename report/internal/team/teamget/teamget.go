package teamget

import (
	"context"
	"net/http"
	"opg-reports/report/internal/team/teamapi/teamapiall"
	"opg-reports/report/package/overwrite"
	"opg-reports/report/package/rest"
	"slices"
	"time"
)

var timeout = (2 * time.Second)

// NavigationData calls the apiHost with known enpoint and returns data formatted for
// the team navigation in the front end
func NavigationData(ctx context.Context, apiHost string, current *http.Request, overwrites ...*rest.Param) (teams []string, err error) {
	var (
		apiResponse *teamapiall.Response
		params      []*rest.Param = []*rest.Param{}
		endpoint    string        = teamapiall.ENDPOINT
		exclude     []string      = []string{"org", "legacy"}
	)
	teams = []string{}
	// overwrite any params
	params = overwrite.Overwrite(params, overwrites...)
	// call the api
	apiResponse, _, err = rest.Get[*teamapiall.Response](ctx, current, &rest.Request{
		Host:     apiHost,
		Endpoint: endpoint,
		Timeout:  timeout,
		Params:   params,
	})

	if err != nil {
		return
	}
	// loop over the api data and return the names as strings
	for _, team := range apiResponse.Data {
		if !slices.Contains(exclude, team.Name) {
			teams = append(teams, team.Name)
		}
	}

	return
}
