package endpoints

import (
	"fmt"
	"slices"
	"strings"
)

const API_VERSION string = "v1"

const HOME string = "/"

// All `teams` related endpoints
const (
	teams         string = "/" + API_VERSION + "/teams/"
	TEAMS_GET_ALL string = teams + "all"
)

// All `awsaccounts` related endpoints
const (
	awsaccounts         string = "/" + API_VERSION + "/awsaccounts/"
	AWSACCOUNTS_GET_ALL string = awsaccounts + "all"
)

// All `awscost` related endpoints
const (
	awscosts            string = "/" + API_VERSION + "/awscosts/"
	AWSCOSTS_GET_TOP_20 string = awscosts + "top20"
	AWSCOSTS_GROUPED    string = awscosts + "grouped/{granularity}/{start_date}/{end_date}"
)

// All `awsuptime` related endpoints
const (
	awsuptime         string = "/" + API_VERSION + "/awsuptime/"
	AWSUPTIME_ALL     string = awsuptime + "all"
	AWSUPTIME_GROUPED string = awsuptime + "grouped/{granularity}/{start_date}/{end_date}"
)

// Parse takes the map and inserts them into the url and subsitutes
// any `{key}` parts with the matching value in the map. Any items in
// the map that dont match a `{key}` will be appended as query string
// parameters
//
// path elements are placed in situ (/test/{name}/all => /test/VALUE/all)
// and all versions are replaced with the same value
//
// query string items are appending in alphabetical order based on their
// key - so `/test?a=value&b=value`
//
// A value of "-" will skip a query string
func Parse(endpoint string, values map[string]string) (ep string) {
	var queryStrings = []string{}
	ep = endpoint

	for key, value := range values {
		pathSub := fmt.Sprintf("{%s}", key)
		if strings.Contains(ep, pathSub) {
			ep = strings.ReplaceAll(ep, pathSub, value)
		} else {
			queryStrings = append(queryStrings, key)
		}
	}
	slices.Sort(queryStrings)

	if len(queryStrings) > 0 {
		ep += "?"
		for _, key := range queryStrings {
			value := values[key]
			if value != "-" {
				ep += fmt.Sprintf("%s=%s&", key, value)
			}
		}
		ep = strings.TrimSuffix(ep, "&")
	}

	return
}
