package server

const testRealisticServerCnf string = `{
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
                            "handler": "MapSlice",
                            "template": "dynamic-aws-costs-units"
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
const testServerCfg string = `{
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

// func TestFrontServerNavigationTop(t *testing.T) {
// 	conf, _ := cnf.Load([]byte(testCfg))
// 	s := New(conf, nil, "", "")

// 	r := httptest.NewRequest(http.MethodGet, "/s1/2/", nil)
// 	active, all := s.Nav.Top(r)

// 	if len(all) != 3 {
// 		t.Errorf("failed to get top nav items")
// 		fmt.Println(all)
// 	}

// 	if active == nil {
// 		t.Errorf("failed to find active top nav")
// 	} else if active.Href != "/s1/" {
// 		t.Errorf("found incorrect top nav")
// 	}

// 	r = httptest.NewRequest(http.MethodGet, "/random/123/", nil)
// 	active, all = s.Nav.Top(r)
// 	if len(all) != 3 {
// 		t.Errorf("failed to get top nav items")
// 	}
// 	if active != nil {
// 		t.Errorf("top should not have been found")
// 	}

// }

// func TestFrontServerNavigationSide(t *testing.T) {
// 	conf, _ := cnf.Load([]byte(testCfg))
// 	s := New(conf, nil, "", "")

// 	r := httptest.NewRequest(http.MethodGet, "/s2/2/1/", nil)
// 	activeT, _ := s.Nav.Top(r)
// 	active, side := s.Nav.Side(r, activeT)

// 	if len(side) != 2 {
// 		t.Errorf("failed to get side nav items")
// 	}

// 	if active == nil {
// 		t.Errorf("should have found active side item")
// 	} else if active.Href != "/s2/2/" {
// 		t.Errorf("found incorrect side item")
// 	}

// }

// func TestFrontServerNavigationActive(t *testing.T) {
// 	conf, _ := cnf.Load([]byte(testCfg))
// 	s := New(conf, nil, "", "")

// 	r := httptest.NewRequest(http.MethodGet, "/s2/2/1/", nil)
// 	active := s.Nav.Active(r)

// 	if active == nil {
// 		t.Errorf("failed to get active item")
// 	} else if active.Href != "/s2/2/1/" || active.Name != "S2.2.1" {
// 		t.Errorf("failed to get correct active item")
// 	}

// 	r = httptest.NewRequest(http.MethodGet, "/s5/2/1/", nil)
// 	active = s.Nav.Active(r)
// 	if active != nil {
// 		t.Errorf("found an item when should not have")
// 	}
// }
