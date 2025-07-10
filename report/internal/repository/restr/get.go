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

	if u, err = url.Parse(uri); err != nil {
		return
	}

	if u.Host != "" {
		u.Host += "/"
	}
	if u.RawQuery != "" {
		u.RawQuery = "?" + u.RawQuery
	}
	parsed = fmt.Sprintf("%s://%s%s%s",
		u.Scheme,
		u.Host,
		strings.TrimPrefix(u.Path, "/"),
		u.RawQuery,
	)
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
