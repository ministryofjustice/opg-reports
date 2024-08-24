package getter

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/ministryofjustice/opg-reports/servers/front/config/navigation"
	"github.com/ministryofjustice/opg-reports/servers/shared/urls"
	"github.com/ministryofjustice/opg-reports/shared/consts"
	"github.com/ministryofjustice/opg-reports/shared/env"
	"github.com/ministryofjustice/opg-reports/shared/timer"
)

func GetUrl(url *url.URL) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return
	}
	apiClient := http.Client{Timeout: consts.API_TIMEOUT}
	resp, err = apiClient.Do(req)
	return
}

// ApiHttpResponses gets api data
func ApiHttpResponses(nav *navigation.NavigationItem, r *http.Request) (responses map[string]*http.Response, requestErr error) {
	apiScheme := env.Get("API_SCHEME", consts.API_SCHEME)
	apiAddr := env.Get("API_ADDR", consts.API_ADDR)
	dataSources := nav.DataSources
	responses = map[string]*http.Response{}

	for key, source := range dataSources {
		endpoint := source.Parsed()
		url := urls.Parse(apiScheme, apiAddr, endpoint)
		tick := timer.New()

		slog.Debug("getting data from api",
			slog.String("key", key),
			slog.String("endpoint", endpoint),
			slog.String("url", url.String()))

		httpResponse, err := GetUrl(url)
		tick.Stop()
		if err != nil || httpResponse.StatusCode != http.StatusOK {
			requestErr = err
			slog.Error("api call failed",
				slog.String("err", fmt.Sprintf("%+v", err)),
				slog.String("key", key),
				slog.Int("status", httpResponse.StatusCode),
				slog.String("endpoint", endpoint),
				slog.String("url", url.String()))
			continue
		}
		responses[key] = httpResponse
	}
	return
}
