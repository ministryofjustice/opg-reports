package endpoints

import (
	"testing"
)

type epTest struct {
	Values   map[string]string
	Ep       string
	Expected string
}

func TestFrontComponentParseEndpoint(t *testing.T) {
	var (
		tests = []epTest{
			{
				Ep: "/v1/awscosts/grouped/{granularity}/{start_date}/{end_date}",
				Values: map[string]string{
					"granularity": "month",
					"start_date":  "2025-01-01",
					"end_date":    "2025-02-01",
					"team":        "true",
				},
				Expected: "/v1/awscosts/grouped/month/2025-01-01/2025-02-01?team=true",
			},
			{
				Ep: "/v1/awscosts/grouped/{granularity}/{start_date}/{end_date}",
				Values: map[string]string{
					"granularity": "month",
					"start_date":  "2025-01-01",
					"end_date":    "2025-02-01",
					"team":        "NAME",
					"account":     "true",
				},
				Expected: "/v1/awscosts/grouped/month/2025-01-01/2025-02-01?account=true&team=NAME",
			},
		}
	)

	for _, test := range tests {
		actual := Parse(test.Ep, test.Values)
		if test.Expected != actual {
			t.Errorf("expected [%s] actual [%s]", test.Expected, actual)
		}
	}
}
