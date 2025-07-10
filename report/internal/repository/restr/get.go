package restr

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"opg-reports/report/internal/utils"
	"strings"
)

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

// Get calls the URI passed in, uses json marshaling to then convert it into
// `result` and returns status code and errors
func (self *Repository) Get(client http.Client, uri string, result interface{}) (statuscode int, err error) {
	var (
		request  *http.Request
		response *http.Response
		content  []byte
		parsed   string
	)
	parsed, err = parseURI(uri)
	if err != nil {
		return
	}
	self.log.With("uri", parsed).Info("calling uri")
	// make the request instance
	request, err = http.NewRequest(http.MethodGet, parsed, nil)
	if err != nil {
		self.log.Error("error creating request", "err", err.Error())
		return
	}

	// make the get call and check it worked
	response, err = client.Do(request)
	if err != nil {
		self.log.Error("error running request", "err", err.Error())
		return
	}
	defer response.Body.Close()

	statuscode = response.StatusCode
	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("recieved unexpected http status: %v", response.StatusCode)
	}

	// now read the content from the response & try convert to T
	content, err = io.ReadAll(response.Body)
	if err != nil {
		return
	}
	err = utils.Unmarshal(content, &result)

	return
}
