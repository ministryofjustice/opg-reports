package config_test

import (
	"fmt"
	"testing"

	"github.com/ministryofjustice/opg-reports/servers/front/config"
)

func TestServersFrontConfigNewSuccess(t *testing.T) {
	working := `{
    "organisation": "FOO",
    "navigation": [
        {
            "name": "BAR",
            "uri": "/",
            "is_header": true,
            "data_sources": {
                "list": "/list/"
            },
            "template": "template"
        }
    ]
}`
	c := config.New([]byte(working))

	if c.Organisation != "FOO" {
		t.Errorf("failed to convert")
		fmt.Printf("%+v\n", c)
	}
}

func TestServersFrontConfigNewFail(t *testing.T) {
	working := `{
    "orgs": "FOO",
    "navs": [
        {
            "name": "BAR",
            "uri": "/",
            "is_header": true,
            "data_sources": {
                "list": "/list/"
            },
            "template": "template"
        }
    ]
}`
	c := config.New([]byte(working))
	if c.Organisation != "" {
		t.Errorf("expected an empty config")
		fmt.Printf("%+v\n", c)
	}
}
