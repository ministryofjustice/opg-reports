package rest

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Request provides structure around the outbound api call
type Request struct {
	Host     string        // api server address
	Endpoint string        // the endpoint pattern to call - like `/v1/infracosts/between/{start_date}/{end_date}`
	Params   []*Param      // list of parameters to pass along to the api
	Timeout  time.Duration // api timeout
}

// URL return string to call, parses the path and query params
// to form the final url
func (self *Request) URL(current *http.Request) (uri string, err error) {
	uri = self.Host + self.Endpoint + "?"

	for _, p := range self.Params {

		var val = p.GetValue(current)
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

// parse and clean up the passsed along uri to be valid host
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
