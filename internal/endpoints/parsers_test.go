package endpoints

import (
	"fmt"
	"testing"
	"time"

	"github.com/ministryofjustice/opg-reports/internal/dateformats"
	"github.com/ministryofjustice/opg-reports/internal/dateintervals"
	"github.com/ministryofjustice/opg-reports/internal/dateutils"
)

// TestEndpointParserYear checks that the {year} keyword
// is correctly replaced within an ApiEndpoint
func TestEndpointParserYear(t *testing.T) {

	var uri ApiEndpoint = "/test/{year}/date"
	var now = dateutils.Reset(time.Now().UTC(), dateintervals.Year).Format(dateformats.YMD)
	var expected = "/test/" + now + "/date"
	var pg = uri.parserGroups()

	if len(pg) != 1 {
		t.Errorf("finding parse groups failed")
	}

	if len(pg[0].Arguments) != 0 {
		t.Errorf("arguments failed")
		fmt.Println(pg[0])
	}

	actual := year(string(uri), pg[0], nil)

	if expected != actual {
		t.Errorf("year parsed failed - expected [%s] actual [%v]", expected, actual)
	}

}

// TestEndpointParserMonthCurrent checks that the {month}
// keyword is replaced with the first of the current month
// correctly
func TestEndpointParserMonthCurrent(t *testing.T) {

	var uri ApiEndpoint = "/test/{month:0,2024-01-01}/date"
	var expected = "/test/2024-01-01/date"
	var pg = uri.parserGroups()

	if len(pg) != 1 {
		t.Errorf("finding parse groups failed")
	}

	if len(pg[0].Arguments) != 2 {
		t.Errorf("arguments failed")
	}

	actual := month(string(uri), pg[0], nil)

	if expected != actual {
		t.Errorf("month parsed failed - expected [%s] actual [%v]", expected, actual)
	}

}

// TestEndpointParserDay checks that the {day} keyword is
// replaced with the start of yesterday in yyyy-mm-dd format
func TestEndpointParserDay(t *testing.T) {

	var uri ApiEndpoint = "/test/{day}/date"
	var now = dateutils.Reset(time.Now().UTC(), dateintervals.Day).AddDate(0, 0, -1).Format(dateformats.YMD)
	var expected = "/test/" + now + "/date"
	var pg = uri.parserGroups()

	if len(pg) != 1 {
		t.Errorf("finding parse groups failed")
	}

	if len(pg[0].Arguments) != 0 {
		t.Errorf("arguments failed")
		fmt.Println(pg[0])
	}

	actual := day(string(uri), pg[0], nil)

	if expected != actual {
		t.Errorf("day parsed failed - expected [%s] actual [%v]", expected, actual)
	}

}
