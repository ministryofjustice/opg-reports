package utils

import (
	"net/http"
	"strings"
)

// MergeRequestWithDefaults is used to pass query string values like start_date from the front end into the
// urls used to query the api, allowing easy way to look at varying dates / teams etc
func MergeRequestWithDefaults(request *http.Request, defaults map[string]string) (merged map[string]string) {
	var include = request.URL.Query()

	merged = defaults

	for key, set := range include {
		var value = strings.TrimSuffix(strings.Join(set, ","), ",")
		merged[key] = value
	}

	return
}
