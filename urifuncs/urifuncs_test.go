package urifuncs_test

import (
	"testing"
	"time"

	"github.com/ministryofjustice/opg-reports/consts"
	"github.com/ministryofjustice/opg-reports/convert"
	"github.com/ministryofjustice/opg-reports/urifuncs"
)

type sdFixture struct {
	Uri      string
	Args     []interface{}
	Expected string
}

// TestUriFuncsStartDate checks parsing out of start_date from
// uris
func TestUriFuncsStartDate(t *testing.T) {
	var fixed = time.Date(2024, 3, 4, 0, 0, 0, 0, time.UTC)
	var now = time.Now().UTC()
	var fourAgoD = convert.DateResetMonth(now).AddDate(0, -4, 0).Format(consts.DateFormatYearMonthDay)
	var fourAgo = "/{version}/test1/" + fourAgoD + "/{end_date}"

	var tests = []*sdFixture{
		{
			Uri:      "/{version}/test1/{start_date}/{end_date}",
			Args:     []interface{}{-4},
			Expected: fourAgo,
		},
		// test a leap year
		{
			Uri:      "/{version}/test1/{start_date}/{end_date}",
			Args:     []interface{}{-4, "day", fixed},
			Expected: "/{version}/test1/2024-02-29/{end_date}",
		},
		{
			Uri:      "/{version}/test1/{start_date}/{end_date}",
			Args:     []interface{}{-4, "year", fixed},
			Expected: "/{version}/test1/2020-01-01/{end_date}",
		},
	}

	for _, test := range tests {
		var expected = test.Expected
		var actual = urifuncs.StartDate(test.Uri, test.Args...)

		if expected != actual {
			t.Errorf("url returned does not match - expected [%s] actual [%v]", expected, actual)
		}

	}

}

// TestUriFuncsEndDate checks parsing out of end_date from uris
func TestUriFuncsEndDate(t *testing.T) {
	var fixed = time.Date(2024, 3, 4, 0, 0, 0, 0, time.UTC)
	var now = time.Now().UTC()
	var noneAgoD = convert.DateResetMonth(now).Format(consts.DateFormatYearMonthDay)
	var noneAgo = "/{version}/test1/{start_date}/" + noneAgoD

	var tests = []*sdFixture{
		{
			Uri:      "/{version}/test1/{start_date}/{end_date}",
			Args:     []interface{}{},
			Expected: noneAgo,
		},
		// test a leap year
		{
			Uri:      "/{version}/test1/{start_date}/{end_date}",
			Args:     []interface{}{-4, "day", fixed},
			Expected: "/{version}/test1/{start_date}/2024-02-29",
		},
		{
			Uri:      "/{version}/test1/{start_date}/{end_date}",
			Args:     []interface{}{-4, "year", fixed},
			Expected: "/{version}/test1/{start_date}/2020-01-01",
		},
	}

	for _, test := range tests {
		var expected = test.Expected
		var actual = urifuncs.EndDate(test.Uri, test.Args...)

		if expected != actual {
			t.Errorf("url returned does not match - expected [%s] actual [%v]", expected, actual)
		}

	}

}

// TestUriFuncsBillingEndDate checks that {end_date}
// is replaced with a correct billing end date
// that is based on the current month
func TestUriFuncsBillingEndDate(t *testing.T) {
	var now = time.Now().UTC()
	var source = "/{version}/test/{start_date}/{end_date}"
	var expected = "/{version}/test/{start_date}/"
	if now.Day() < consts.CostsBillingDay {
		expected += convert.DateResetMonth(now).AddDate(0, -2, 0).Format(consts.DateFormatYearMonthDay)
	} else {
		expected += convert.DateResetMonth(now).AddDate(0, -1, 0).Format(consts.DateFormatYearMonthDay)
	}

	actual := urifuncs.BillingEndDate(source)
	if expected != actual {
		t.Errorf("url returned does not match - expected [%s] actual [%v]", expected, actual)
	}
}

// TestUriFuncsBillingStartDate checks that {start_date}
// is replace based on either the current billing month
// like billing end date, or with a modifier to go further
// back in the past
func TestUriFuncsBillingStartDate(t *testing.T) {
	var now = time.Now().UTC()
	var source = "/{version}/test/{start_date}"
	var base = "/{version}/test/"
	var expected = base
	var expDate = convert.DateResetMonth(now).AddDate(0, -2, 0)
	if now.Day() >= consts.CostsBillingDay {
		expDate = convert.DateResetMonth(now).AddDate(0, -1, 0)
	}

	expected += expDate.Format(consts.DateFormatYearMonthDay)
	actual := urifuncs.BillingStartDate(source)
	if expected != actual {
		t.Errorf("url returned does not match - expected [%s] actual [%v]", expected, actual)
	}

	// check the modifier works
	actual = urifuncs.BillingStartDate(source, -2)
	expected = base + expDate.AddDate(0, -2, 0).Format(consts.DateFormatYearMonthDay)
	if expected != actual {
		t.Errorf("url returned does not match - expected [%s] actual [%v]", expected, actual)
	}

}

func TestUriFuncsInterval(t *testing.T) {
	var source = "/{version}/{interval}/test"
	var expected = "/{version}/month/test"
	var actual = urifuncs.Interval(source)

	if expected != actual {
		t.Errorf("url returned does not match - expected [%s] actual [%v]", expected, actual)
	}

	expected = "/{version}/day/test"
	actual = urifuncs.Interval(source, "day")
	if expected != actual {
		t.Errorf("url returned does not match - expected [%s] actual [%v]", expected, actual)
	}

	expected = "/{version}/year/test"
	actual = urifuncs.Interval(source, "year")
	if expected != actual {
		t.Errorf("url returned does not match - expected [%s] actual [%v]", expected, actual)
	}

	expected = "/{version}/month/test"
	actual = urifuncs.Interval(source, "hour")
	if expected != actual {
		t.Errorf("url returned does not match - expected [%s] actual [%v]", expected, actual)
	}
	expected = "/{version}/month/test"
	actual = urifuncs.Interval(source, -1)
	if expected != actual {
		t.Errorf("url returned does not match - expected [%s] actual [%v]", expected, actual)
	}

}
