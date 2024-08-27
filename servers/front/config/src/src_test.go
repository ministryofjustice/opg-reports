package src_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ministryofjustice/opg-reports/servers/front/config/src"
	"github.com/ministryofjustice/opg-reports/shared/consts"
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
	var n = time.Now().UTC()
	var now = n.Format(dates.FormatYM)
	var expected string = fmt.Sprintf("/test/%s/all", now)
	var simple src.ApiUrl = "/test/{month}/all"

	p := simple.Parsed()
	if p != expected {
		t.Errorf("date parsing failed: expected [%s] actual [%s]", expected, p)
	}

	expected = fmt.Sprintf("/test/all?start=%s", now)
	simple = "/test/all?start={month}"
	p = simple.Parsed()
	if p != expected {
		t.Errorf("date parsing failed: expected [%s] actual [%s]", expected, p)
	}

	m := n.AddDate(0, -1, 0).Format(dates.FormatYM)
	expected = fmt.Sprintf("/test/%s/%s", m, now)
	simple = "/test/{month:-1}/{month}"
	p = simple.Parsed()
	if p != expected {
		t.Errorf("date parsing failed: expected [%s] actual [%s]", expected, p)
	}

}

func TestServersFrontConfigSrcDay(t *testing.T) {
	logger.LogSetup()
	var n = time.Now().UTC()
	var now = n.Format(dates.FormatYMD)
	var expected string = fmt.Sprintf("/test/%s/all", now)
	var simple src.ApiUrl = "/test/{day}/all"

	p := simple.Parsed()
	if p != expected {
		t.Errorf("date parsing failed: expected [%s] actual [%s]", expected, p)
	}

	expected = fmt.Sprintf("/test/all?start=%s", now)
	simple = "/test/all?start={day}"
	p = simple.Parsed()
	if p != expected {
		t.Errorf("date parsing failed: expected [%s] actual [%s]", expected, p)
	}

	m := n.AddDate(0, 0, -1).Format(dates.FormatYMD)
	expected = fmt.Sprintf("/test/%s/%s", m, now)
	simple = "/test/{day:-1}/{day}"
	p = simple.Parsed()
	if p != expected {
		t.Errorf("date parsing failed: expected [%s] actual [%s]", expected, p)
	}

}

func TestServersFrontConfigSrcBillingMonth(t *testing.T) {
	logger.LogSetup()
	var n = time.Now().UTC()
	var b = dates.BillingEndDate(n, consts.BILLING_DATE)
	var s = b.AddDate(0, -1, 0).Format(dates.FormatYM)
	var now = b.Format(dates.FormatYM)
	var expected string = fmt.Sprintf("/test/%s/%s/all", s, now)
	var simple src.ApiUrl = "/test/{billingMonth:-1}/{billingMonth}/all"

	p := simple.Parsed()
	if p != expected {
		t.Errorf("date parsing failed: expected [%s] actual [%s]", expected, p)
	}

}

func TestServersFrontConfigSrcBillingDay(t *testing.T) {
	logger.LogSetup()
	var n = time.Now().UTC()
	var b = dates.BillingEndDate(n, consts.BILLING_DATE)
	var s = b.AddDate(0, 0, -30).Format(dates.FormatYMD)
	var now = b.Format(dates.FormatYMD)
	var expected string = fmt.Sprintf("/test/%s/%s/all", s, now)
	var simple src.ApiUrl = "/test/{billingDay:-30}/{billingDay}/all"

	p := simple.Parsed()
	if p != expected {
		t.Errorf("date parsing failed: expected [%s] actual [%s]", expected, p)
	}

}
