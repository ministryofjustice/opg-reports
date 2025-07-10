package restr

import (
	"fmt"
	"io"
	"net/http"
	"opg-reports/report/internal/utils"
	"strings"
)

// parse and clean up the passsed along uri
func parseURI(uri string) (parsed string, err error) {

	var (
		schema      string = "http"
		host        string = "localhost"
		path        string = ""
		queryString string = ""
		chunks      []string
	)

	// add trailing ?
	if !strings.Contains(uri, "?") {
		uri += "?"
	}
	// grab the first schema
	// 	"http://localhost:8080/test/v1?test=1"
	// 		=> "http", "localhost:8080/test/v1?test=1"
	chunks = strings.Split(uri, "://")
	// shift if we have more than 1 item (as in schema + remainder)
	if len(chunks) > 1 {
		schema, uri = chunks[0], strings.Join(chunks[1:], "")
	}

	// so the first / should show end of the hostname
	// "localhost:8080/test/v1?test=1"
	// 		=> "localhost:8080", "test/v1?test=1"
	chunks = strings.Split(uri, "/")
	if len(chunks) > 1 && chunks[0] != "" {
		host, uri = chunks[0], strings.Join(chunks[1:], "/")
	}

	// "test/v1?test=1"
	// 		=> "test/v1", "test=1"
	chunks = strings.Split(uri, "?")
	if len(chunks) > 1 {
		path, queryString = chunks[0], strings.Join(chunks[1:], "")
	}

	hostPath := fmt.Sprintf("%s/%s", host, path)
	hostPath = strings.ReplaceAll(hostPath, "//", "/")
	parsed = fmt.Sprintf("%s://%s?%s",
		schema,
		hostPath,
		queryString,
	)
	parsed = strings.TrimSuffix(parsed, "?")

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
