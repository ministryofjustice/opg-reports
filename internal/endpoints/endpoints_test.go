package endpoints

import (
	"fmt"
	"testing"

	"github.com/ministryofjustice/opg-reports/pkg/bi"
)

// TestEndpointParserGroupsCountMatch checks that
// patterns within strings are found in correctly
func TestEndpointParserGroupsCountMatch(t *testing.T) {
	var checks = map[string]int{
		"/nothing/in/here":                           0,
		"/fixed/{version}/{month:-5}/{month:-1}/end": 3,
		"/version/{version}/{thing:-1, 2, test}":     2,
		"whats-this":                                 0,
		"hello {my-name}!":                           1,
		"hello {your name}!":                         1,
	}

	for url, expected := range checks {
		var ep = ApiEndpoint(url)
		var pg = ep.parserGroups()
		var actual = len(pg)

		if expected != actual {
			t.Errorf("incorrect number of groups - expected [%d] actual [%v]\n[%s]", expected, actual, url)
		}
	}

}

// TestEndpointParserGroupsArgsCountMatch checks that for
// a single matched chunk the argsuments are parsed correctly
func TestEndpointParserGroupsArgsCountMatch(t *testing.T) {
	var checks = map[string]int{
		"{month}":                      0,
		"{month:}":                     0,
		"{month:-5}":                   1,
		"/{thing:-1, 2, test}/foobar/": 3,
	}

	for url, expected := range checks {
		var ep = ApiEndpoint(url)
		var pg = ep.parserGroups()
		var actual = len(pg[0].Arguments)

		if expected != actual {
			t.Errorf("incorrect number args - expected [%d] actual [%v]\n[%s]", expected, actual, url)
			fmt.Println(pg[0].Arguments)
		}
	}

}

func TestEndpointParsing(t *testing.T) {
	var checks = map[string]string{
		"/test/{month:0,2024-01-01}/end":                            "/test/2024-01-01/end",
		"/test/{month:1,2024-01-20}/end":                            "/test/2024-02-01/end",
		"/{version}/{month:-1,2024-03-15}/end":                      "/" + bi.ApiVersion + "/2024-02-01/end",
		"/test/{year:0,2024-11-09}/end":                             "/test/2024-01-01/end",
		"/test/{day:-1,2024-03-01}/end":                             "/test/2024-02-29/end",
		"/test/{day:1,2024-02-28}/end":                              "/test/2024-02-29/end",
		"/{billing_date:0,2024-04-16}":                              "/2024-03-01",
		"/{billing_date:0,2024-04-14}":                              "/2024-02-01",
		"/{billing_date:-4,2024-08-16}/{billing_date:0,2024-08-16}": "/2024-03-01/2024-07-01",
	}

	for uri, expected := range checks {
		var ep = ApiEndpoint(uri)
		var actual = ep.Parse()
		if expected != actual {
			t.Errorf("url parse failed - expected [%s] actual [%s]", expected, actual)
		}
	}

}
