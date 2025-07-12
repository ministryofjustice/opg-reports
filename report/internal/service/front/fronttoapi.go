package front

import (
	"net/http"
	"strings"
)

// mergeRequestWithMaps is used to pass query string values like start_date from the front end into the
// urls used to query the api, allowing easy way to look at varying dates / teams etc
func mergeRequestWithMaps(request *http.Request, maps ...map[string]string) (merged map[string]string) {
	var include = request.URL.Query()
	merged = map[string]string{}
	// loop over all the maps and set those first
	for _, mp := range maps {
		for key, val := range mp {
			merged[key] = val
		}
	}
	// then replace them with any values in the query string
	for key, set := range include {
		var value = strings.TrimSuffix(strings.Join(set, ","), ",")
		merged[key] = value
	}

	return
}
