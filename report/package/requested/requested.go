package requested

import (
	"context"
	"net/http"
	"opg-reports/report/package/cnv"
)

// Parse finds all the keys in T in both the path & query string data sets
// and updates the dest with their values
//
// Used to convert an incoming request to an input data set
func Parse[T any](ctx context.Context, request *http.Request, dest T) (err error) {

	var specMap = map[string]string{}
	var resultMap = map[string]string{}

	// var found = map[string][]string{}
	// convert to a map
	err = cnv.Convert(dest, &specMap)
	if err != nil {
		return
	}
	// fetch from the path
	for k, _ := range specMap {
		if v := request.PathValue(k); v != "" {
			resultMap[k] = v
		}
	}
	// fetch from query string
	queryValues := request.URL.Query()
	for k, _ := range specMap {
		if v := queryValues.Get(k); v != "" {
			resultMap[k] = v
		}
	}
	// now convert back to T
	err = cnv.Convert(resultMap, dest)
	return
}
