package getter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"opg-reports/report/packages/slogx"
	"strings"
	"time"
)

var ErrRequestFailed = errors.New("request failed.")

type QueryStringer interface {
	QueryString() string
}

type Source struct {
	Host    string
	Path    string
	Timeout time.Duration
}

func (self *Source) Url() string {
	return fmt.Sprintf("%s%s", self.Host, self.Path)
}

// Get fetches the data from a remote source and tries to
// convert to a struct via json unmarshaling
func Get[T any](ctx context.Context, source *Source, qs QueryStringer) (result T, err error) {
	var content []byte
	var called string
	var log = slogx.FromContext(ctx)

	log.Info(ctx, "fetching data from source ... ")
	content, _, called, err = get(ctx, source, qs)
	log.Info(ctx, "called url", "url", called)

	if err != nil {
		return
	}

	err = json.Unmarshal(content, &result)
	if err != nil {
		return
	}

	log.Info(ctx, "data fetching complete.")
	return

}

// get is a helper to fetch data from an endpoint
func get(ctx context.Context, source *Source, qs QueryStringer) (content []byte, statusCode int, calledURL string, err error) {
	var (
		request  *http.Request
		response *http.Response
		uri      string
		client   http.Client = http.Client{Timeout: source.Timeout}
		log                  = slogx.FromContext(ctx)
	)

	uri, err = parseURI(source.Url() + qs.QueryString())
	if err != nil {
		return
	}
	calledURL = uri
	log.Debug(ctx, "calling uri ...", "uri", uri)
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
