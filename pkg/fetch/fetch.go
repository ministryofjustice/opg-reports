// Package fetch is used to handle making outbound calls driven from navigation data
package fetch

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ministryofjustice/opg-reports/pkg/consts"
	"github.com/ministryofjustice/opg-reports/pkg/endpoints"
)

// Fetch gets data from the uri passsed and returns the repsonse, response code and
// any error
// This is used to go get data from the api
//
// Uses Response to make the call to get the content
func Fetch(host string, uri endpoints.ApiEndpoint) (content []byte, code int, err error) {
	var (
		response *http.Response
		url      string = host + uri.Parse()
	)

	response, err = Response(url, consts.FetchTimeout)
	if err != nil {
		return
	}
	defer response.Body.Close()

	code = response.StatusCode
	content, err = io.ReadAll(response.Body)

	return
}

// Response creates a new request and calls the url (as http get) returning the
// http.Response
// If the status code is not a 200 then an error is returned
func Response(url string, timeout time.Duration) (response *http.Response, err error) {
	var (
		request *http.Request
		client  http.Client
	)

	request, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	client = http.Client{Timeout: timeout}
	response, err = client.Do(request)

	if err == nil && response.StatusCode != http.StatusOK {
		err = fmt.Errorf("expected status [%d] actual [%v]", http.StatusOK, response.StatusCode)
	}

	return
}
