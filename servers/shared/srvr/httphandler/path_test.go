package httphandler_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/httphandler"
	"github.com/ministryofjustice/opg-reports/shared/consts"
	"github.com/ministryofjustice/opg-reports/shared/dates"
	"github.com/ministryofjustice/opg-reports/shared/logger"
)

func TestSharedSrvrHttphandlerPathNoReplacements(t *testing.T) {
	logger.LogSetup()

	simple := "/test/all"
	actual := httphandler.Path(simple)

	if actual != simple {
		t.Errorf("should be a match, wasnt")
	}
}

func TestSharedSrvrHttphandlerPathMonth(t *testing.T) {
	logger.LogSetup()
	var n = time.Now().UTC()
	var now = n.Format(dates.FormatYM)
	var expected string = fmt.Sprintf("/test/%s/all", now)
	var simple = "/test/{month}/all"

	p := httphandler.Path(simple)
	if p != expected {
		t.Errorf("date parsing failed: expected [%s] actual [%s]", expected, p)
	}

	expected = fmt.Sprintf("/test/all?start=%s", now)
	simple = "/test/all?start={month}"
	p = httphandler.Path(simple)
	if p != expected {
		t.Errorf("date parsing failed: expected [%s] actual [%s]", expected, p)
	}

	m := n.AddDate(0, -1, 0).Format(dates.FormatYM)
	expected = fmt.Sprintf("/test/%s/%s", m, now)
	simple = "/test/{month:-1}/{month}"
	p = httphandler.Path(simple)
	if p != expected {
		t.Errorf("date parsing failed: expected [%s] actual [%s]", expected, p)
	}

}

func TestSharedSrvrHttphandlerPathDay(t *testing.T) {
	logger.LogSetup()
	var n = time.Now().UTC()
	var now = n.Format(dates.FormatYMD)
	var expected string = fmt.Sprintf("/test/%s/all", now)
	var simple = "/test/{day}/all"

	p := httphandler.Path(simple)
	if p != expected {
		t.Errorf("date parsing failed: expected [%s] actual [%s]", expected, p)
	}

	expected = fmt.Sprintf("/test/all?start=%s", now)
	simple = "/test/all?start={day}"
	p = httphandler.Path(simple)
	if p != expected {
		t.Errorf("date parsing failed: expected [%s] actual [%s]", expected, p)
	}

	m := n.AddDate(0, 0, -1).Format(dates.FormatYMD)
	expected = fmt.Sprintf("/test/%s/%s", m, now)
	simple = "/test/{day:-1}/{day}"
	p = httphandler.Path(simple)
	if p != expected {
		t.Errorf("date parsing failed: expected [%s] actual [%s]", expected, p)
	}

}

func TestSharedSrvrHttphandlerPathBillingMonth(t *testing.T) {
	logger.LogSetup()
	var n = time.Now().UTC()
	var b = dates.BillingEndDate(n, consts.BILLING_DATE)
	var s = b.AddDate(0, -1, 0).Format(dates.FormatYM)
	var now = b.Format(dates.FormatYM)
	var expected string = fmt.Sprintf("/test/%s/%s/all", s, now)
	var simple string = "/test/{billingMonth:-1}/{billingMonth}/all"

	p := httphandler.Path(simple)
	if p != expected {
		t.Errorf("date parsing failed: expected [%s] actual [%s]", expected, p)
	}

}

func TestSharedSrvrHttphandlerPathBillingDay(t *testing.T) {
	logger.LogSetup()
	var n = time.Now().UTC()
	var b = dates.BillingEndDate(n, consts.BILLING_DATE)
	var s = b.AddDate(0, 0, -30).Format(dates.FormatYMD)
	var now = b.Format(dates.FormatYMD)
	var expected string = fmt.Sprintf("/test/%s/%s/all", s, now)
	var simple = "/test/{billingDay:-30}/{billingDay}/all"

	p := httphandler.Path(simple)
	if p != expected {
		t.Errorf("date parsing failed: expected [%s] actual [%s]", expected, p)
	}

}
