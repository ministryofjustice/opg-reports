package server

import (
	"net/http"
	"net/http/httptest"
	"opg-reports/shared/server"
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
                            "api": "/aws/costs/{apiVersion}/monthly/{startYM}/{endYM}/",
                            "handler": "MapMapSlice",
                            "template": "dynamic-aws-costs-totals"
                        },
                        {
                            "name": "By Unit",
                            "href": "/costs/aws/units/",
                            "api": "/aws/costs/{apiVersion}/monthly/{startYM}/{endYM}/units/",
                            "handler": "MapSlice"
                        },
                        {
                            "name": "By Unit & Environment",
                            "href": "/costs/aws/units-envs/",
                            "api": "/aws/costs/{apiVersion}/monthly/{startYM}/{endYM}/units/envs/",
                            "handler": "MapSlice"
                        },
                        {
                            "name": "Detailed Breakdown",
                            "href": "/costs/aws/detailed/",
                            "api": "/aws/costs/{apiVersion}/monthly/{startYM}/{endYM}/units/envs/services/",
                            "handler": "MapSlice"
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
var mockServerType server.ApiResponseConstraintString = "map[string]map[string][]*cost.Cost"
var mockServerResponse string = `{
    "timings": {
        "start": "2024-06-25T17:13:54.516722Z",
        "end": "2024-06-25T17:13:54.516773Z",
        "duration": 51000
    },
    "status": 200,
    "errors": [],
    "result_type": "map[string]map[string][]*cost.Cost",
    "result": {
        "with_tax": {
            "month^2023-12.": [
                {
                    "uuid": "c911b825-9be8-43ea-83c8-8fb5b7d3e1d4",
                    "account_organsiation": "Cn3",
                    "account_id": "244002",
                    "account_environment": "DKb",
                    "account_name": "gkYWi6PXWs",
                    "account_unit": "gMIV",
                    "account_label": "K82RR",
                    "service": "ecs",
                    "region": "1gZC",
                    "date": "2023-12-22T01:24:12Z",
                    "cost": "5464.788727"
                }
            ],
            "month^2024-01.": [
                {
                    "uuid": "61392250-18ea-4323-aaf6-00b37be881c2",
                    "account_organsiation": "pm7",
                    "account_id": "386289",
                    "account_environment": "bcX",
                    "account_name": "fxYOfS5Pmn",
                    "account_unit": "mPIb",
                    "account_label": "4hGmO",
                    "service": "ec2",
                    "region": "RaNf",
                    "date": "2024-01-12T19:11:55Z",
                    "cost": "5096.116557"
                },
                {
                    "uuid": "ff41b2a4-c669-49fa-a8c4-841897999345",
                    "account_organsiation": "gpP",
                    "account_id": "794712",
                    "account_environment": "CSn",
                    "account_name": "JOORuHI3Is",
                    "account_unit": "GgqT",
                    "account_label": "5ZaNx",
                    "service": "ec2",
                    "region": "wb1z",
                    "date": "2024-01-22T08:59:34Z",
                    "cost": "7397.055649"
                },
                {
                    "uuid": "a53a79af-1fa5-45c8-910a-77c7d9075f4c",
                    "account_organsiation": "9Ly",
                    "account_id": "129324",
                    "account_environment": "N8M",
                    "account_name": "aJlYMsuoKd",
                    "account_unit": "FyKY",
                    "account_label": "YrM59",
                    "service": "ecs",
                    "region": "UlG2",
                    "date": "2024-01-10T21:45:37Z",
                    "cost": "5327.629128"
                }
            ],
            "month^2024-02.": [
                {
                    "uuid": "54b8846e-9272-40fb-85d9-c04da9c29ba4",
                    "account_organsiation": "nU0",
                    "account_id": "821882",
                    "account_environment": "gzF",
                    "account_name": "eeZwO5L5Oy",
                    "account_unit": "N4O1",
                    "account_label": "LjB3R",
                    "service": "r53",
                    "region": "ShdA",
                    "date": "2024-02-27T11:12:21Z",
                    "cost": "3116.061929"
                }
            ],
            "month^2024-03.": [
                {
                    "uuid": "fcc75cb0-317a-4f56-8881-911011430bfa",
                    "account_organsiation": "tO5",
                    "account_id": "847563",
                    "account_environment": "CXB",
                    "account_name": "nTC1xYD8Nl",
                    "account_unit": "PGyT",
                    "account_label": "f6ntN",
                    "service": "rds",
                    "region": "BieP",
                    "date": "2024-03-20T11:08:45Z",
                    "cost": "3696.366160"
                }
            ],
            "month^2024-05.": [
                {
                    "uuid": "5558cb54-6bf7-4f79-ab49-4b273fa3dff3",
                    "account_organsiation": "NTw",
                    "account_id": "385018",
                    "account_environment": "DbO",
                    "account_name": "LefPLZeTRx",
                    "account_unit": "3Ys0",
                    "account_label": "giX9P",
                    "service": "tax",
                    "region": "duWW",
                    "date": "2024-05-15T19:09:24+01:00",
                    "cost": "1258.826527"
                }
            ],
            "month^2024-06.": [
                {
                    "uuid": "6a15f37f-9882-400e-a04b-2fee41f7ef01",
                    "account_organsiation": "09A",
                    "account_id": "415612",
                    "account_environment": "Kra",
                    "account_name": "BCkUqTQNB1",
                    "account_unit": "4k7M",
                    "account_label": "fAW3f",
                    "service": "r53",
                    "region": "YOAd",
                    "date": "2024-06-23T05:25:24+01:00",
                    "cost": "9596.483581"
                },
                {
                    "uuid": "1a05cb54-7997-4559-8857-f126d5350671",
                    "account_organsiation": "ulA",
                    "account_id": "628321",
                    "account_environment": "n1O",
                    "account_name": "2K7odlcLzl",
                    "account_unit": "cQLk",
                    "account_label": "qppYs",
                    "service": "r53",
                    "region": "QepB",
                    "date": "2024-06-30T02:12:43+01:00",
                    "cost": "5000.559594"
                }
            ]
        },
        "without_tax": {
            "month^2023-12.": [
                {
                    "uuid": "c911b825-9be8-43ea-83c8-8fb5b7d3e1d4",
                    "account_organsiation": "Cn3",
                    "account_id": "244002",
                    "account_environment": "DKb",
                    "account_name": "gkYWi6PXWs",
                    "account_unit": "gMIV",
                    "account_label": "K82RR",
                    "service": "ecs",
                    "region": "1gZC",
                    "date": "2023-12-22T01:24:12Z",
                    "cost": "5464.788727"
                }
            ],
            "month^2024-01.": [
                {
                    "uuid": "ff41b2a4-c669-49fa-a8c4-841897999345",
                    "account_organsiation": "gpP",
                    "account_id": "794712",
                    "account_environment": "CSn",
                    "account_name": "JOORuHI3Is",
                    "account_unit": "GgqT",
                    "account_label": "5ZaNx",
                    "service": "ec2",
                    "region": "wb1z",
                    "date": "2024-01-22T08:59:34Z",
                    "cost": "7397.055649"
                },
                {
                    "uuid": "a53a79af-1fa5-45c8-910a-77c7d9075f4c",
                    "account_organsiation": "9Ly",
                    "account_id": "129324",
                    "account_environment": "N8M",
                    "account_name": "aJlYMsuoKd",
                    "account_unit": "FyKY",
                    "account_label": "YrM59",
                    "service": "ecs",
                    "region": "UlG2",
                    "date": "2024-01-10T21:45:37Z",
                    "cost": "5327.629128"
                },
                {
                    "uuid": "61392250-18ea-4323-aaf6-00b37be881c2",
                    "account_organsiation": "pm7",
                    "account_id": "386289",
                    "account_environment": "bcX",
                    "account_name": "fxYOfS5Pmn",
                    "account_unit": "mPIb",
                    "account_label": "4hGmO",
                    "service": "ec2",
                    "region": "RaNf",
                    "date": "2024-01-12T19:11:55Z",
                    "cost": "5096.116557"
                }
            ],
            "month^2024-02.": [
                {
                    "uuid": "54b8846e-9272-40fb-85d9-c04da9c29ba4",
                    "account_organsiation": "nU0",
                    "account_id": "821882",
                    "account_environment": "gzF",
                    "account_name": "eeZwO5L5Oy",
                    "account_unit": "N4O1",
                    "account_label": "LjB3R",
                    "service": "r53",
                    "region": "ShdA",
                    "date": "2024-02-27T11:12:21Z",
                    "cost": "3116.061929"
                }
            ],
            "month^2024-03.": [
                {
                    "uuid": "fcc75cb0-317a-4f56-8881-911011430bfa",
                    "account_organsiation": "tO5",
                    "account_id": "847563",
                    "account_environment": "CXB",
                    "account_name": "nTC1xYD8Nl",
                    "account_unit": "PGyT",
                    "account_label": "f6ntN",
                    "service": "rds",
                    "region": "BieP",
                    "date": "2024-03-20T11:08:45Z",
                    "cost": "3696.366160"
                }
            ],
            "month^2024-06.": [
                {
                    "uuid": "1a05cb54-7997-4559-8857-f126d5350671",
                    "account_organsiation": "ulA",
                    "account_id": "628321",
                    "account_environment": "n1O",
                    "account_name": "2K7odlcLzl",
                    "account_unit": "cQLk",
                    "account_label": "qppYs",
                    "service": "r53",
                    "region": "QepB",
                    "date": "2024-06-30T02:12:43+01:00",
                    "cost": "5000.559594"
                },
                {
                    "uuid": "6a15f37f-9882-400e-a04b-2fee41f7ef01",
                    "account_organsiation": "09A",
                    "account_id": "415612",
                    "account_environment": "Kra",
                    "account_name": "BCkUqTQNB1",
                    "account_unit": "4k7M",
                    "account_label": "fAW3f",
                    "service": "r53",
                    "region": "YOAd",
                    "date": "2024-06-23T05:25:24+01:00",
                    "cost": "9596.483581"
                }
            ]
        }
    }
}`

func mockServerAWSCostTotals() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(server.ResponseTypeHeader, string(mockServerType))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockServerResponse))
	}))
	return server
}

func testMux() *http.ServeMux {
	return http.NewServeMux()
}
func testWRGet(route string) (*httptest.ResponseRecorder, *http.Request) {
	return httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, route, nil)
}
