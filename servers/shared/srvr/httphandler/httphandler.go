package httphandler

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/ministryofjustice/opg-reports/servers/shared/urls"
)

// TIMEOUT is used to control the max duration of an external call
const TIMEOUT time.Duration = time.Second * 4

type HttpHandler struct {
	DataSource *HttpGet
	Response   *http.Response
	StatusCode int
	Duration   float64
	Url        *url.URL
}

func (a *HttpHandler) Result() (response *http.Response, err error) {
	err = a.Get()
	if err == nil {
		response = a.Response
	}
	return
}

// Get calls remotedata source methof to fetch it and
func (a *HttpHandler) Get() (err error) {
	if err = a.DataSource.Get(); err == nil {
		a.Response = a.DataSource.Response
		a.StatusCode = a.Response.StatusCode
	}
	a.Duration = a.DataSource.Duration

	return
}

// New creates a pre configures HttpHandler
func New(scheme string, addr string, path string) *HttpHandler {
	path = Path(path)
	uri := urls.Parse(scheme, addr, path)
	return &HttpHandler{
		DataSource: &HttpGet{Url: uri},
		Response:   nil,
		Duration:   0.0,
		StatusCode: 0,
		Url:        uri,
	}
}

// Get combines creating a new instance of HttpHandler and then also calls
// .Get() method to retrive the data directly
func Get(scheme string, addr string, path string) (response *HttpHandler, err error) {
	response = New(scheme, addr, path)
	err = response.Get()
	return
}

// GetAll fetches multiple api responses at once tracking them in the map keys passed along
//
// the map is from config.nav datasources
func GetAll(scheme string, addr string, paths map[string]string) (responses map[string]*HttpHandler, err error) {
	responses = map[string]*HttpHandler{}

	for key, path := range paths {
		response, rErr := Get(scheme, addr, path)
		// if theres an error, log it and skip
		if rErr != nil || response.StatusCode != http.StatusOK {
			slog.Error("api call failed",
				slog.String("err", fmt.Sprintf("%+v", rErr)),
				slog.String("key", key),
				slog.Int("status", response.StatusCode),
				slog.String("url", response.Url.String()))
			continue
		}
		responses[key] = response
	}
	if len(responses) == 0 && len(paths) > 0 {
		err = fmt.Errorf("failed to get any api responses")
	}
	return
}
