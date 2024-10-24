package endpoints

import (
	"fmt"
	"testing"
	"time"

	"github.com/ministryofjustice/opg-reports/consts"
	"github.com/ministryofjustice/opg-reports/convert"
)

func TestEndpointParserYear(t *testing.T) {

	var uri ApiEndpoint = "/test/{year}/date"
	var now = convert.DateResetYear(time.Now().UTC()).Format(consts.DateFormatYearMonthDay)
	var expected = "/test/" + now + "/date"
	var pg = uri.parserGroups()

	if len(pg) != 1 {
		t.Errorf("finding parse groups failed")
	}

	if len(pg[0].Arguments) != 0 {
		t.Errorf("arguments failed")
		fmt.Println(pg[0])
	}

	actual := year(string(uri), pg[0])

	if expected != actual {
		t.Errorf("year parsed failed - expected [%s] actual [%v]", expected, actual)
	}

}

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

	actual := month(string(uri), pg[0])

	if expected != actual {
		t.Errorf("month parsed failed - expected [%s] actual [%v]", expected, actual)
	}

}

func TestEndpointParserDay(t *testing.T) {

	var uri ApiEndpoint = "/test/{day}/date"
	var now = convert.DateResetDay(time.Now().UTC()).AddDate(0, 0, -1).Format(consts.DateFormatYearMonthDay)
	var expected = "/test/" + now + "/date"
	var pg = uri.parserGroups()

	if len(pg) != 1 {
		t.Errorf("finding parse groups failed")
	}

	if len(pg[0].Arguments) != 0 {
		t.Errorf("arguments failed")
		fmt.Println(pg[0])
	}

	actual := day(string(uri), pg[0])

	if expected != actual {
		t.Errorf("day parsed failed - expected [%s] actual [%v]", expected, actual)
	}

}
