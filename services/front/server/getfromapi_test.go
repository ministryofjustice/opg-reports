package server

import (
	"testing"
)

func TestFrontServerGetFromApiMockedDetails(t *testing.T) {
	ms := mockServerAWSCostTotals()
	defer ms.Close()
	url := ms.URL
	resType, _, err := GetFromApi(url)

	if resType != mockServerType {
		t.Errorf("res type failed")
	}
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
