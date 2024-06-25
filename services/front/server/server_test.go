package server

import (
	"net/http"
	"net/http/httptest"
)

var testRealisticCfg = `{
    "organisation": "OPG",
    "sections": [
        {
            "name": "Home",
            "href": "/",
            "exclude": true,
            "template": "static-home"
        },
        {
            "name": "Costs",
            "template": "static-costs-home",
			"href": "/costs/",
            "sections": [
                {
                    "name": "AWS Costs",
                    "href": "/costs/aws/",
                    "header": true,
                    "sections": [
                        {
                            "name": "Totals",
                            "href": "/costs/aws/totals/",
                            "api": "/aws/costs/{apiVersion}/monhtly/{startDate}/{endDate}/"
                        },
                        {
                            "name": "By Unit",
                            "href": "/costs/aws/units/",
                            "api": "/aws/costs/{apiVersion}/monhtly/{startDate}/{endDate}/units/"
                        },
                        {
                            "name": "By Unit & Environment",
                            "href": "/costs/aws/units-envs/",
                            "api": "/aws/costs/{apiVersion}/monhtly/{startDate}/{endDate}/units/envs/"
                        },
                        {
                            "name": "Detailed Breakdown",
                            "href": "/costs/aws/detailed/",
                            "api": "/aws/costs/{apiVersion}/monhtly/{startDate}/{endDate}/units/envs/services/"
                        }
                    ]
                }
            ]
        }
    ]
}`

var testCfg = `{
	"organisation": "test-org",
	"sections": [
		{
			"name": "Home",
			"href": "/",
			"exclude": true,
			"template": "static-home"
		},
		{
			"name": "Section 1",
			"href": "/s1/",
			"template": "static-home",
			"sections": [
				{
					"Name": "S1.1",
					"href": "/s1/1/"
				},
				{
					"Name": "S1.2",
					"href": "/s1/2/"
				}
			]
		},
		{
			"name": "Section 2",
			"href": "/s2/",
			"sections": [
				{
					"Name": "S2.1",
					"href": "/s2/1/"
				},
				{
					"Name": "S2.2",
					"href": "/s2/2/",
					"sections": [
						{
							"Name": "S2.2.1",
							"href": "/s2/2/1/"
						}
					]
				}
			]
		}
	]
}`

func testMux() *http.ServeMux {
	return http.NewServeMux()
}
func testWRGet(route string) (*httptest.ResponseRecorder, *http.Request) {
	return httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, route, nil)
}
