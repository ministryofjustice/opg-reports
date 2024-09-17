package apihandlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/ministryofjustice/opg-reports/servers/api/github_standards"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/httphandler"
	"github.com/ministryofjustice/opg-reports/shared/convert"
)

type GithubStandards struct{}

func (gh *GithubStandards) Handles(url string) bool {
	segment := "/github-standards/"
	return strings.Contains(url, segment)
}
func (gh *GithubStandards) Handle(data map[string]interface{}, key string, remote *httphandler.HttpHandler, w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		response *http.Response
		api      *github_standards.GHSResponse
	)
	response, err = remote.Result()
	slog.Info("called api",
		slog.Float64("duration", remote.Duration),
		slog.Int("status_code", remote.StatusCode),
		slog.String("uri", remote.Url.String()))

	if err != nil || remote.StatusCode != http.StatusOK {
		slog.Error("error with api response", slog.String("err", fmt.Sprintf("%+v", err)), slog.Int("status_code", remote.StatusCode))
		return
	}

	api, err = convert.UnmarshalR[*github_standards.GHSResponse](response)
	if err != nil {
		slog.Error("error with unmarshal response", slog.String("err", err.Error()))
		return
	}
	data[key] = api
}
