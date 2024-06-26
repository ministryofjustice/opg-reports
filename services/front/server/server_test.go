package server

import (
	"net/http"
	"net/http/httptest"
)

const mockAwsCostTotalsResponse string = `{
    "timings": {
        "start": "2024-06-26T16:26:04.226572Z",
        "end": "2024-06-26T16:26:04.232055Z",
        "duration": 5483000
    },
    "status": 200,
    "errors": [],
    "result": {
        "headings": {
            "cells": [
                {
                    "name": "",
                    "value": ""
                },
                {
                    "name": "2023-12",
                    "value": ""
                },
                {
                    "name": "2024-01",
                    "value": ""
                },
                {
                    "name": "2024-02",
                    "value": ""
                },
                {
                    "name": "2024-03",
                    "value": ""
                }
            ]
        },
        "rows": [
            {
                "cells": [
                    {
                        "name": "Without Tax",
                        "value": "Without Tax"
                    },
                    {
                        "name": "2023-12",
                        "value": "1329101.583654"
                    },
                    {
                        "name": "2024-01",
                        "value": "1252176.507914"
                    },
                    {
                        "name": "2024-02",
                        "value": "1117055.982422"
                    },
                    {
                        "name": "2024-03",
                        "value": "0.000000"
                    }
                ]
            },
            {
                "cells": [
                    {
                        "name": "With Tax",
                        "value": "With Tax"
                    },
                    {
                        "name": "2023-12",
                        "value": "1664994.069984"
                    },
                    {
                        "name": "2024-01",
                        "value": "1481791.853741"
                    },
                    {
                        "name": "2024-02",
                        "value": "1414703.158002"
                    },
                    {
                        "name": "2024-03",
                        "value": "0.000000"
                    }
                ]
            }
        ]
    }
}`
const mockAwsCostUnitsResponse string = `{
    "timings": {
        "start": "2024-06-26T16:25:27.264395Z",
        "end": "2024-06-26T16:25:27.268054Z",
        "duration": 3659000
    },
    "status": 200,
    "errors": [],
    "result": {
        "headings": {
            "cells": [
                {
                    "name": "Unit",
                    "value": ""
                },
                {
                    "name": "2023-12",
                    "value": ""
                },
                {
                    "name": "2024-01",
                    "value": ""
                },
                {
                    "name": "2024-02",
                    "value": ""
                },
                {
                    "name": "2024-03",
                    "value": ""
                }
            ]
        },
        "rows": [
            {
                "cells": [
                    {
                        "name": "teamTwo",
                        "value": "teamTwo"
                    },
                    {
                        "name": "2023-12",
                        "value": "463581.575155"
                    },
                    {
                        "name": "2024-01",
                        "value": "399935.015446"
                    },
                    {
                        "name": "2024-02",
                        "value": "446120.042567"
                    },
                    {
                        "name": "2024-03",
                        "value": "0.000000"
                    }
                ]
            },
            {
                "cells": [
                    {
                        "name": "teamThree",
                        "value": "teamThree"
                    },
                    {
                        "name": "2023-12",
                        "value": "606809.650879"
                    },
                    {
                        "name": "2024-01",
                        "value": "537143.712587"
                    },
                    {
                        "name": "2024-02",
                        "value": "488254.288767"
                    },
                    {
                        "name": "2024-03",
                        "value": "0.000000"
                    }
                ]
            },
            {
                "cells": [
                    {
                        "name": "teamOne",
                        "value": "teamOne"
                    },
                    {
                        "name": "2023-12",
                        "value": "547004.109814"
                    },
                    {
                        "name": "2024-01",
                        "value": "446731.122453"
                    },
                    {
                        "name": "2024-02",
                        "value": "446662.885689"
                    },
                    {
                        "name": "2024-03",
                        "value": "0.000000"
                    }
                ]
            }
        ]
    }
}`
const mockAwsCostUnitEnvsResponse string = `{
    "timings": {
        "start": "2024-06-26T16:24:49.149363Z",
        "end": "2024-06-26T16:24:49.149415Z",
        "duration": 52000
    },
    "status": 200,
    "errors": [],
    "result": {
        "headings": {
            "cells": [
                {
                    "name": "Unit",
                    "value": ""
                },
                {
                    "name": "Environment",
                    "value": ""
                },
                {
                    "name": "2023-12",
                    "value": ""
                },
                {
                    "name": "2024-01",
                    "value": ""
                },
                {
                    "name": "2024-02",
                    "value": ""
                },
                {
                    "name": "2024-03",
                    "value": ""
                }
            ]
        },
        "rows": [
            {
                "cells": [
                    {
                        "name": "teamTwo",
                        "value": "teamTwo"
                    },
                    {
                        "name": "dev",
                        "value": "dev"
                    },
                    {
                        "name": "2023-12",
                        "value": "3111.488494"
                    },
                    {
                        "name": "2024-01",
                        "value": "0.000000"
                    },
                    {
                        "name": "2024-02",
                        "value": "0.000000"
                    },
                    {
                        "name": "2024-03",
                        "value": "0.000000"
                    }
                ]
            },
            {
                "cells": [
                    {
                        "name": "teamOne",
                        "value": "teamOne"
                    },
                    {
                        "name": "prod",
                        "value": "prod"
                    },
                    {
                        "name": "2023-12",
                        "value": "0.000000"
                    },
                    {
                        "name": "2024-01",
                        "value": "2704.267815"
                    },
                    {
                        "name": "2024-02",
                        "value": "0.000000"
                    },
                    {
                        "name": "2024-03",
                        "value": "0.000000"
                    }
                ]
            },
            {
                "cells": [
                    {
                        "name": "teamTwo",
                        "value": "teamTwo"
                    },
                    {
                        "name": "prod",
                        "value": "prod"
                    },
                    {
                        "name": "2023-12",
                        "value": "0.000000"
                    },
                    {
                        "name": "2024-01",
                        "value": "2820.031982"
                    },
                    {
                        "name": "2024-02",
                        "value": "0.000000"
                    },
                    {
                        "name": "2024-03",
                        "value": "0.000000"
                    }
                ]
            },
            {
                "cells": [
                    {
                        "name": "teamThree",
                        "value": "teamThree"
                    },
                    {
                        "name": "dev",
                        "value": "dev"
                    },
                    {
                        "name": "2023-12",
                        "value": "2042.579808"
                    },
                    {
                        "name": "2024-01",
                        "value": "8825.299206"
                    },
                    {
                        "name": "2024-02",
                        "value": "3119.397026"
                    },
                    {
                        "name": "2024-03",
                        "value": "0.000000"
                    }
                ]
            },
            {
                "cells": [
                    {
                        "name": "teamThree",
                        "value": "teamThree"
                    },
                    {
                        "name": "prod",
                        "value": "prod"
                    },
                    {
                        "name": "2023-12",
                        "value": "9073.416082"
                    },
                    {
                        "name": "2024-01",
                        "value": "6504.427462"
                    },
                    {
                        "name": "2024-02",
                        "value": "0.000000"
                    },
                    {
                        "name": "2024-03",
                        "value": "0.000000"
                    }
                ]
            },
            {
                "cells": [
                    {
                        "name": "teamOne",
                        "value": "teamOne"
                    },
                    {
                        "name": "dev",
                        "value": "dev"
                    },
                    {
                        "name": "2023-12",
                        "value": "0.000000"
                    },
                    {
                        "name": "2024-01",
                        "value": "8449.203764"
                    },
                    {
                        "name": "2024-02",
                        "value": "0.000000"
                    },
                    {
                        "name": "2024-03",
                        "value": "0.000000"
                    }
                ]
            }
        ]
    }
}`
const mockAwsCostUnitEnvServicesResponse string = `{
    "timings": {
        "start": "2024-06-26T16:23:59.882394Z",
        "end": "2024-06-26T16:23:59.88246Z",
        "duration": 66000
    },
    "status": 200,
    "errors": [],
    "result": {
        "headings": {
            "cells": [
                {
                    "name": "Account",
                    "value": ""
                },
                {
                    "name": "Unit",
                    "value": ""
                },
                {
                    "name": "Environment",
                    "value": ""
                },
                {
                    "name": "Service",
                    "value": ""
                },
                {
                    "name": "2023-12",
                    "value": ""
                },
                {
                    "name": "2024-01",
                    "value": ""
                },
                {
                    "name": "2024-02",
                    "value": ""
                },
                {
                    "name": "2024-03",
                    "value": ""
                }
            ]
        },
        "rows": [
            {
                "cells": [
                    {
                        "name": "515146",
                        "value": "515146"
                    },
                    {
                        "name": "teamOne",
                        "value": "teamOne"
                    },
                    {
                        "name": "dev",
                        "value": "dev"
                    },
                    {
                        "name": "VlHuY9Yh",
                        "value": "VlHuY9Yh"
                    },
                    {
                        "name": "2023-12",
                        "value": "0.000000"
                    },
                    {
                        "name": "2024-01",
                        "value": "7206.463369"
                    },
                    {
                        "name": "2024-02",
                        "value": "0.000000"
                    },
                    {
                        "name": "2024-03",
                        "value": "0.000000"
                    }
                ]
            },
            {
                "cells": [
                    {
                        "name": "783083",
                        "value": "783083"
                    },
                    {
                        "name": "teamThree",
                        "value": "teamThree"
                    },
                    {
                        "name": "dev",
                        "value": "dev"
                    },
                    {
                        "name": "eZGddgAY",
                        "value": "eZGddgAY"
                    },
                    {
                        "name": "2023-12",
                        "value": "0.000000"
                    },
                    {
                        "name": "2024-01",
                        "value": "0.000000"
                    },
                    {
                        "name": "2024-02",
                        "value": "7699.819039"
                    },
                    {
                        "name": "2024-03",
                        "value": "0.000000"
                    }
                ]
            },
            {
                "cells": [
                    {
                        "name": "322933",
                        "value": "322933"
                    },
                    {
                        "name": "teamThree",
                        "value": "teamThree"
                    },
                    {
                        "name": "preprod",
                        "value": "preprod"
                    },
                    {
                        "name": "Tj5ELRYC",
                        "value": "Tj5ELRYC"
                    },
                    {
                        "name": "2023-12",
                        "value": "1337.020225"
                    },
                    {
                        "name": "2024-01",
                        "value": "0.000000"
                    },
                    {
                        "name": "2024-02",
                        "value": "0.000000"
                    },
                    {
                        "name": "2024-03",
                        "value": "0.000000"
                    }
                ]
            },
            {
                "cells": [
                    {
                        "name": "683652",
                        "value": "683652"
                    },
                    {
                        "name": "teamOne",
                        "value": "teamOne"
                    },
                    {
                        "name": "dev",
                        "value": "dev"
                    },
                    {
                        "name": "tXlnSCF6",
                        "value": "tXlnSCF6"
                    },
                    {
                        "name": "2023-12",
                        "value": "0.000000"
                    },
                    {
                        "name": "2024-01",
                        "value": "0.000000"
                    },
                    {
                        "name": "2024-02",
                        "value": "5493.534212"
                    },
                    {
                        "name": "2024-03",
                        "value": "0.000000"
                    }
                ]
            },
            {
                "cells": [
                    {
                        "name": "128316",
                        "value": "128316"
                    },
                    {
                        "name": "teamTwo",
                        "value": "teamTwo"
                    },
                    {
                        "name": "preprod",
                        "value": "preprod"
                    },
                    {
                        "name": "eJw5raua",
                        "value": "eJw5raua"
                    },
                    {
                        "name": "2023-12",
                        "value": "0.000000"
                    },
                    {
                        "name": "2024-01",
                        "value": "9382.370540"
                    },
                    {
                        "name": "2024-02",
                        "value": "0.000000"
                    },
                    {
                        "name": "2024-03",
                        "value": "0.000000"
                    }
                ]
            },
            {
                "cells": [
                    {
                        "name": "368206",
                        "value": "368206"
                    },
                    {
                        "name": "teamThree",
                        "value": "teamThree"
                    },
                    {
                        "name": "prod",
                        "value": "prod"
                    },
                    {
                        "name": "7f0BQzZV",
                        "value": "7f0BQzZV"
                    },
                    {
                        "name": "2023-12",
                        "value": "0.000000"
                    },
                    {
                        "name": "2024-01",
                        "value": "0.000000"
                    },
                    {
                        "name": "2024-02",
                        "value": "8675.389262"
                    },
                    {
                        "name": "2024-03",
                        "value": "0.000000"
                    }
                ]
            },
            {
                "cells": [
                    {
                        "name": "173261",
                        "value": "173261"
                    },
                    {
                        "name": "teamOne",
                        "value": "teamOne"
                    },
                    {
                        "name": "dev",
                        "value": "dev"
                    },
                    {
                        "name": "ifoX4hVV",
                        "value": "ifoX4hVV"
                    },
                    {
                        "name": "2023-12",
                        "value": "7247.300916"
                    },
                    {
                        "name": "2024-01",
                        "value": "0.000000"
                    },
                    {
                        "name": "2024-02",
                        "value": "0.000000"
                    },
                    {
                        "name": "2024-03",
                        "value": "0.000000"
                    }
                ]
            },
            {
                "cells": [
                    {
                        "name": "768978",
                        "value": "768978"
                    },
                    {
                        "name": "teamOne",
                        "value": "teamOne"
                    },
                    {
                        "name": "dev",
                        "value": "dev"
                    },
                    {
                        "name": "yecSR2eo",
                        "value": "yecSR2eo"
                    },
                    {
                        "name": "2023-12",
                        "value": "0.000000"
                    },
                    {
                        "name": "2024-01",
                        "value": "1131.035523"
                    },
                    {
                        "name": "2024-02",
                        "value": "0.000000"
                    },
                    {
                        "name": "2024-03",
                        "value": "0.000000"
                    }
                ]
            },
            {
                "cells": [
                    {
                        "name": "292724",
                        "value": "292724"
                    },
                    {
                        "name": "teamTwo",
                        "value": "teamTwo"
                    },
                    {
                        "name": "dev",
                        "value": "dev"
                    },
                    {
                        "name": "MmeEB0xG",
                        "value": "MmeEB0xG"
                    },
                    {
                        "name": "2023-12",
                        "value": "0.000000"
                    },
                    {
                        "name": "2024-01",
                        "value": "2414.480628"
                    },
                    {
                        "name": "2024-02",
                        "value": "0.000000"
                    },
                    {
                        "name": "2024-03",
                        "value": "0.000000"
                    }
                ]
            }
        ]
    }
}`

func mockServer(resp string, status int) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Write([]byte(resp))
	}))
	return server
}
func testMux() *http.ServeMux {
	return http.NewServeMux()
}
func testWRGet(route string) (*httptest.ResponseRecorder, *http.Request) {
	return httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, route, nil)
}
