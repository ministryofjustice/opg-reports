package cnf

import (
	"fmt"
	"testing"
)

func TestFrontCnfFlat(t *testing.T) {
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
	conf, _ := Load([]byte(content))
	f := map[string]*SiteSection{}
	FlatSections(conf.Sections, f)

	if len(f) != 8 {
		t.Errorf("flatten failed")
	}

}

func TestFrontCnfLoadSimple(t *testing.T) {
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

func TestFrontCnfLoadWithSections(t *testing.T) {
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

func TestFrontCnfLoadNested(t *testing.T) {
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
