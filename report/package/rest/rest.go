package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"opg-reports/report/package/cntxt"
)

var ErrRequestFailed = errors.New("request failed.")

// Get is a helper to fetch json based data from an endpoint, mixing default parameters
// (both path & query string) with values from the current request - allowing the front
// end to voerwrite things like start dates directly.
//
// Data returned is converted into R via json unmarshal
func Get[R any](ctx context.Context, current *http.Request, req *Request) (result R, statusCode int, err error) {
	var (
		request  *http.Request
		response *http.Response
		content  []byte
		uri      string
		client   http.Client  = http.Client{Timeout: req.Timeout}
		log      *slog.Logger = cntxt.GetLogger(ctx).With("package", "rest", "func", "Get")
	)
	// generate the url to call
	uri, err = req.URL(current)
	if err != nil {
		return
	}
	log.Info("calling uri ...", "uri", uri)
	// create request
	request, err = http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return
	}
	// req request
	response, err = client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	// check status
	statusCode = response.StatusCode
	if statusCode != http.StatusOK {
		err = errors.Join(ErrRequestFailed, fmt.Errorf("returned status code [%d]", statusCode))
		return
	}
	// read
	content, err = io.ReadAll(response.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(content, &result)
	if err != nil {
		return
	}

	log.Info("done")
	return
}
