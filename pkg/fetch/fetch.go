// Package fetch is used to handle making outbound calls driven from navigation data
package fetch

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/ministryofjustice/opg-reports/pkg/consts"
	"github.com/ministryofjustice/opg-reports/pkg/endpoints"
)

// Fetch gets data from the uri passsed and returns the repsonse, response code and
// any error
// This is used to go get data from the api
func Fetch(host string, uri endpoints.ApiEndpoint) (content []byte, code int, err error) {
	var (
		request  *http.Request
		response *http.Response
		client   http.Client
		url      string = host + uri.Parse()
	)

	request, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		slog.Error("[fetch.Fetch] error generating request")
		return
	}

	client = http.Client{Timeout: consts.ApiTimeout}
	response, err = client.Do(request)
	if err != nil {
		slog.Error("[fetch.Fetch] error calling do")
		return
	}

	code = response.StatusCode
	content, err = io.ReadAll(response.Body)
	if err != nil {
		slog.Error("[fetch.Fetch] stringify failed", slog.String("err", err.Error()))
	}

	return
}
