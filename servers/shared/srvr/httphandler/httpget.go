package httphandler

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/ministryofjustice/opg-reports/shared/timer"
)

type HttpGet struct {
	Url      *url.URL
	Response *http.Response
	Duration float64
}

// Get calls the remote data source in the url and sets the response
func (h *HttpGet) Get() (err error) {
	var (
		request *http.Request
		client  http.Client
		uri     = h.Url.String()
	)

	tick := timer.New()
	if request, err = http.NewRequest(http.MethodGet, uri, nil); err == nil {
		client = http.Client{Timeout: TIMEOUT}
		h.Response, err = client.Do(request)
	}
	tick.Stop()
	h.Duration = tick.Duration()
	slog.Debug("fetched data from remote",
		slog.String("url", uri),
		slog.Float64("duration", h.Duration),
		slog.String("err", fmt.Sprintf("%+v", err)))
	return
}
