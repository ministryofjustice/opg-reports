package page_test

import (
	"testing"

	"github.com/ministryofjustice/opg-reports/servers/front/config/page"
)

func TestServerFrontConfigPage(t *testing.T) {

	simple := page.Data{
		"list": "/",
	}

	if len(simple) != 1 {
		t.Errorf("page data failed")
	}

}
