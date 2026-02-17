package requested

import (
	"context"
	"net/http"
	"opg-reports/report/package/cnv"
)

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

// func pathVariableValues(ctx context.Context, pathPattern string, request *http.Request, positions map[string][]int) (values map[string]string) {
// 	var uri = request.PathValue()
// 	fmt.Println(uri)
// 	return
// }

// // pathVariables finds all variables in url path (`{name}`) and retusn a list of the idx
// // position its founds at
// func pathVariables(ctx context.Context, pathPattern string) (positions map[string][]int) {
// 	var chunks = strings.Split(pathPattern, "/")
// 	positions = map[string][]int{}

// 	for i, chunk := range chunks {
// 		if len(chunk) > 0 && chunk[0] == '{' && chunk[len(chunk)-1] == '}' {
// 			key := strings.TrimPrefix(strings.TrimSuffix(chunk, "}"), "{")
// 			if _, ok := positions[key]; !ok {
// 				positions[key] = []int{}
// 			}
// 			positions[key] = append(positions[key], i)
// 		}
// 	}

// 	return
// }
