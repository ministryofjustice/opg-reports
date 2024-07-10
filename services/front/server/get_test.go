package server

import (
	"net/http"
	th "opg-reports/internal/testhelpers"
	"testing"
)

func TestFrontServerGetApiMockedDetails(t *testing.T) {
	content := mockAwsCostMonthlyTotals()
	ms := th.MockServer(content, http.StatusOK)
	defer ms.Close()
	url := ms.URL
	_, err := GetUrl(url)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestFrontServerGetFromApiGetUrl(t *testing.T) {

	u := Url("", ":8081", "/aws/costs/v1/monthly")
	str := u.String()
	if str != "http://localhost:8081/aws/costs/v1/monthly/" {
		t.Errorf("url mapping failed")
	}

}
