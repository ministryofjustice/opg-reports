package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"opg-reports/report/internal/utils/unmarshal"
	"strings"
	"time"
)

type Type string

const (
	QUERY Type = "query"
	PATH  Type = "path"
)

var ErrRequestFailed = errors.New("request failed.")

type Response interface{}

// Param
type Param struct {
	Key    string
	Type   Type
	Value  string
	Locked bool // stops value being replaced by request
}

func (self *Param) PathKey() string {
	return fmt.Sprintf("{%s}", self.Key)
}

// GetValue will give the value for this param to use in api query by checking for it
// in the current front end server request - allowing pass through
func (self *Param) GetValue(current *http.Request) (v string) {
	var values = current.URL.Query()
	v = self.Value

	// if locked, then cant be overwrittern by request
	if self.Locked {
		return self.Value
	}

	for key, set := range values {
		if key == self.Key {
			v = strings.TrimSuffix(strings.Join(set, ","), ",")
		}
	}

	return
}

type Call struct {
	Host     string        // api server address
	Endpoint string        // the endpoint pattern to call - like `/v1/infracosts/between/{start_date}/{end_date}`
	Request  *http.Request // the current request to the front end, used to grab params to pass along to api
	Params   []*Param      // list of parameters to pass along to the api
	Timeout  time.Duration // api timeout
}

// URL retusn string to call, parses the host and query params
// to form the final url
func (self *Call) URL() (uri string, err error) {
	uri = self.Host + self.Endpoint + "?"

	for _, p := range self.Params {

		var val = p.GetValue(self.Request)
		if p.Type == PATH && val != "" {
			uri = strings.ReplaceAll(uri, p.PathKey(), val)
		}
		if p.Type == QUERY && val != "" {
			uri += fmt.Sprintf("%s=%s&", p.Key, val)
		}
	}
	uri, err = parseURI(uri)
	uri = strings.TrimSuffix(uri, "&")
	uri = strings.TrimSuffix(uri, "?")

	return
}

// Get fetches data from the specified location (normally the api) and returns that detail
// converting the resposne into R on the way
func Get[R Response](ctx context.Context, log *slog.Logger, call *Call) (result R, statusCode int, err error) {
	var (
		request  *http.Request
		response *http.Response
		content  []byte
		uri      string
		client   http.Client  = http.Client{Timeout: call.Timeout}
		lg       *slog.Logger = log.With("func", "api.Get")
	)
	uri, err = call.URL()
	if err != nil {
		return
	}
	lg.Info("calling api ...", "uri", uri)
	// create request
	request, err = http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return
	}
	// call request
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
	err = unmarshal.Unmarshal(content, &result)
	if err != nil {
		return
	}

	lg.Info("done")
	return
}

// parse and clean up the passsed along uri
func parseURI(uri string) (parsed string, err error) {
	var u *url.URL

	u, err = url.Parse(uri)
	if err != nil {
		return
	}
	// if there is no path (so http://localhost), add one
	// but if there is a path and the first character is
	// not '/', then the hostname is likely in here.
	//
	// Typically this comes from a uri like
	// `localhost/path-to-file.json` so we can fix
	// the hostname and path
	if u.Path == "" {
		u.Path = "/"
	} else if u.Path[0] != '/' {
		chunks := strings.Split(u.Path, "/")
		u.Host = chunks[0]
		u.Path = strings.Join(chunks[1:], "/")
	}
	if u.Host == "" {
		u.Host = "localhost"
	}
	if u.Scheme == "" {
		u.Scheme = "http"
	}

	parsed = u.String()
	// Sometimes when the scheme is not stated the
	// String() may not add one, so look for the ://
	// deliminator, if we cant find it, then add
	// http:// as a default
	if !strings.Contains(parsed, "://") {
		parsed = "http://" + parsed
	}
	return
}
