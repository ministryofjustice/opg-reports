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
	"time"
)

var ErrRequestFailed = errors.New("request failed.")

func FromApi[R any](ctx context.Context, apiHost string, endpoint string, current *http.Request, params ...*Param) (response R, err error) {
	var timeout time.Duration = (2 * time.Second)

	response, _, err = Get[R](ctx, current, &Request{
		Host:     apiHost,
		Endpoint: endpoint,
		Timeout:  timeout,
		Params:   params,
	})
	return
}

// Get is a helper to fetch json based data from an endpoint, mixing default parameters
// (both path & query string) with values from the current request - allowing the front
// end to overwrite things like start dates directly.
//
// Data returned is converted into R via json unmarshal
func Get[R any](ctx context.Context, current *http.Request, req *Request) (result R, statusCode int, err error) {
	var content []byte
	var called string = ""
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "rest", "func", "Get")

	log.Info("Starting ...")
	content, statusCode, called, err = get(ctx, current, req)
	log.Info("called url", "url", called)
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

// Get is a helper to fetch json based data from an endpoint, mixing default parameters
// (both path & query string) with values from the current request - allowing the front
// end to overwrite things like start dates directly.
//
// Data returned is converted to a string
func GetStr(ctx context.Context, current *http.Request, req *Request) (result string, statusCode int, err error) {
	var content []byte
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "rest", "func", "Get")

	log.Debug("Starting ...")
	content, statusCode, _, err = get(ctx, current, req)
	if err != nil {
		return
	}
	result = string(content)
	log.Debug("done")
	return
}

// get is a helper to fetch json based data from an endpoint, mixing default parameters
// (both path & query string) with values from the current request - allowing the front
// end to voerwrite things like start dates directly.
func get(ctx context.Context, current *http.Request, req *Request) (content []byte, statusCode int, calledURL string, err error) {
	var (
		request  *http.Request
		response *http.Response
		uri      string
		client   http.Client  = http.Client{Timeout: req.Timeout}
		log      *slog.Logger = cntxt.GetLogger(ctx).With("package", "rest", "func", "get")
	)
	// generate the url to call
	uri, err = req.URL(current)
	if err != nil {
		return
	}
	calledURL = uri
	log.Debug("calling uri ...", "uri", uri)
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
	return
}
