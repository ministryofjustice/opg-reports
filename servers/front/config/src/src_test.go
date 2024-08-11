package src_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ministryofjustice/opg-reports/servers/front/config/src"
	"github.com/ministryofjustice/opg-reports/shared/dates"
	"github.com/ministryofjustice/opg-reports/shared/logger"
)

func TestServersFrontConfigSrcNoReplacements(t *testing.T) {
	logger.LogSetup()

	var simple src.ApiUrl = "/test/all"

	if (simple.Parsed()) != string(simple) {
		t.Errorf("parse failed")
	}
}

func TestServersFrontConfigSrcMonth(t *testing.T) {
	logger.LogSetup()
	var now = time.Now().UTC().Format(dates.FormatYM)
	var expected string = fmt.Sprintf("/test/%s/all", now)
	var simple src.ApiUrl = "/test/{month}/all"

	p := simple.Parsed()
	if p != expected {
		t.Errorf("date parsing failed: expected [%s] actual [%s]", expected, p)
	}
}
