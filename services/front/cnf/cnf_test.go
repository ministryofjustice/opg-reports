package cnf

import (
	"fmt"
	"opg-reports/shared/dates"
	"strings"
	"testing"
	"time"
)

func TestServicesFrontCnfSubs(t *testing.T) {
	sample := &SiteSection{
		Name:         "example",
		Href:         "/home/",
		Header:       false,
		Exclude:      false,
		TemplateName: "any",
		Api: map[string]string{
			"costs-now":   "/costs/v1/{month}/{month}/",
			"costs-range": "/costs/v1/{month:-6}/{month:-1}/",
		},
	}
	now := time.Now().UTC()
	res, _ := sample.ApiUrls()
	urlNow := res["costs-now"]

	if strings.Contains(urlNow, "{") || strings.Contains(urlNow, "}") {
		t.Errorf("standard month sub failed")
	}
	if !strings.Contains(urlNow, "/"+now.Format(dates.FormatYM)+"/") {
		t.Errorf("incorrect date: [%s] = [%s]", urlNow, now.Format(dates.FormatYM))
	}

	ago := res["costs-range"]
	ago6 := now.AddDate(0, -6, 0)
	ago1 := now.AddDate(0, -1, 0)

	if !strings.Contains(ago, "/"+ago6.Format(dates.FormatYM)+"/") || !strings.Contains(ago, "/"+ago1.Format(dates.FormatYM)+"/") {
		t.Errorf("relative month sub failed")
	}

}

func TestServicesFrontCnfFlat(t *testing.T) {
	content := `{
		"organisation": "test-org",
		"sections": [
			{
				"name": "Home",
				"href": "/",
				"exclude": true
			},
			{
				"name": "Section 1",
				"href": "/s1/",
				"sections": [
					{"Name": "S1.1", "href": "/s1/1/"},
					{"Name": "S1.2", "href": "/s1/2/"}
				]
			},
			{
				"name": "Section 2",
				"href": "/s2/",
				"sections": [
					{"Name": "S2.1", "href": "/s2/1/"},
					{
						"Name": "S2.2",
						"href": "/s2/2/",
						"sections": [
							{"Name": "S2.2.1", "href": "/s2/2/1/"}
						]
					}
				]
			}
		]
	}`
	conf, err := Load([]byte(content))
	if err != nil {
		t.Errorf("failed: %v", err)
	}
	f := map[string]*SiteSection{}
	FlatSections(conf.Sections, f)

	if len(f) != 8 {
		t.Errorf("flatten failed")
		fmt.Println(len(f))
		fmt.Println(f)
	}

}

func TestServicesFrontCnfLoadSimple(t *testing.T) {
	content := `{"organisation": "test"}`
	cnf, err := Load([]byte(content))

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if cnf.Organisation != "test" {
		t.Errorf("parse failed")
		fmt.Printf("%+v\n", cnf)
	}
	if len(cnf.Sections) != 0 {
		t.Errorf("sections should be empty")
	}
}

func TestServicesFrontCnfLoadWithSections(t *testing.T) {
	content := `{"organisation": "test", "sections": [{"name":"s1"}] }`
	cnf, err := Load([]byte(content))

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if cnf.Organisation != "test" {
		t.Errorf("parse failed")
		fmt.Printf("%+v\n", cnf)
	}
	if len(cnf.Sections) != 1 {
		t.Errorf("sections should have an entry")
	}

	if cnf.Sections[0].Name != "s1" {
		t.Errorf("section name failed")
	}
}

func TestServicesFrontCnfLoadNested(t *testing.T) {
	content := `{
		"organisation": "test",
		"sections": [
			{"name":"s1"},
			{
				"name":"s2",
				"sections": [
					{"name": "s2.1", "href": "/"},
					{"name": "s2.3"}
				]
			}
		]
	}`
	cnf, err := Load([]byte(content))

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if cnf.Organisation != "test" {
		t.Errorf("parse failed")
		fmt.Printf("%+v\n", cnf)
	}
	if len(cnf.Sections) != 2 {
		t.Errorf("sections should have entries")
	}

	if cnf.Sections[0].Name != "s1" {
		t.Errorf("section name failed")
	}

	if len(cnf.Sections[1].Sections) == 0 {
		t.Errorf("s2 should have sub sections")
	}
}

func TestServicesFrontCnfBillingMonth(t *testing.T) {
	postBilling := time.Date(2024, 7, 31, 17, 0, 0, 0, time.UTC)
	res := billingMonth("{billingMonth}", "", "/{billingMonth}/", &postBilling)
	if res != "/2024-06/" {
		t.Errorf("unepxected month returned for billing: %v", res)
	}

	postBilling = time.Date(2024, 2, 29, 17, 0, 0, 0, time.UTC)
	res = billingMonth("{billingMonth}", "", "/{billingMonth}/", &postBilling)
	if res != "/2024-01/" {
		t.Errorf("unepxected month returned for billing: %v", res)
	}

	postBilling = time.Date(2024, 2, billingDay-1, 17, 0, 0, 0, time.UTC)
	res = billingMonth("{billingMonth}", "", "/{billingMonth}/", &postBilling)
	if res != "/2023-12/" {
		t.Errorf("unepxected month returned for billing: %v", res)
	}

}

func TestServicesFrontCnfMonth(t *testing.T) {
	m := time.Date(2024, 7, 31, 17, 0, 0, 0, time.UTC)
	res := month("{m:-1}", "-1", "/{m:-1}/", &m)
	if res != "/2024-06/" {
		t.Errorf("unepxected month returned: %v", res)
	}
	m = time.Date(2024, 2, 29, 23, 0, 0, 0, time.UTC)
	res = month("{m:-1}", "-1", "/{m:-1}/", &m)
	if res != "/2024-01/" {
		t.Errorf("unepxected month returned: %v", res)
	}
	res = month("{m:-2}", "-2", "/{m:-2}/", &m)
	if res != "/2023-12/" {
		t.Errorf("unepxected month returned: %v", res)
	}
}
