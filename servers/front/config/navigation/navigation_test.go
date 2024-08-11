package navigation_test

import (
	"testing"

	"github.com/ministryofjustice/opg-reports/servers/front/config/navigation"
	"github.com/ministryofjustice/opg-reports/shared/testhelpers"
)

func TestServersFrontConfigNavigationNewSuccess(t *testing.T) {
	working := `[
        {
            "name": "BAR",
            "uri": "/bar",
            "is_header": true,
            "data_sources": {
                "list": "/list/"
            },
            "template": "template1",
			"navigation": [
				{
					"name": "Sub",
					"uri": "/bar/sub",
					"is_header": true,
					"template": "tp1"
				}
			]
        },
		{
            "name": "FOO",
            "uri": "/foo",
            "is_header": true,
            "data_sources": {
                "list": "/list/"
            },
            "template": "template2"
        }
    ]`

	n := navigation.New([]byte(working))

	if len(n) != 2 {
		t.Errorf("length error: expected: [2] actual: [%d]", len(n))
	}
	_, r := testhelpers.WRGet("/bar/sub?test=123&abc=123")
	all := navigation.Flat(n, r)

	if !all["/bar"].Active {
		t.Errorf("active check failed")
	}
	if !all["/bar/sub"].Active {
		t.Errorf("active sub check failed")
	}
	if all["/foo"].Active {
		t.Errorf("active check failed")
	}

	lvl0, a := navigation.Level(n, r)
	if len(lvl0) != 2 {
		t.Errorf("level failed")
	}
	if a.Name != "BAR" {
		t.Errorf("active failed")
	}

	found := navigation.ForTemplate("tp1", n)
	if found == nil || found.Name != "Sub" {
		t.Errorf("get by template failed")
	}

}
