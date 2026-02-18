package requested

import (
	"context"
	"net/http"
	"opg-reports/report/package/cnv"
)

// Parse finds all the keys in T in both the path & query string data sets
// and updates the dest with their values
//
// Converts incoming request to an input struct; mappping path segments
// like `/v1/{month}/` or `/v1/?month=jan` to `dest.Month` - presuming dest has
// a property with suitable json annotation.
//
// Generally used in api handlers to map parts of http request into a struct
// which can then be used for filtering / triggering actions in the function
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
