package server

import (
	"testing"
)

func TestFrontServerGetFromApiGetUrl(t *testing.T) {

	u := Url("", ":8081", "/home?test=1&team=1&group=test")
	str := u.String()
	if str != "http://localhost:8081/home/?test=1&team=1&group=test" {
		t.Errorf("url mapping failed: %v", str)
	}

	u = Url("https://", "", "/home/test/?team=1&group=test")
	str = u.String()
	if str != "https://localhost/home/test/?team=1&group=test" {
		t.Errorf("url mapping failed: %v", str)
	}

	u = Url("https://", "localhost", "https://localhost/home/test?team=1&group=test")
	str = u.String()
	if str != "https://localhost/home/test/?team=1&group=test" {
		t.Errorf("url mapping failed: %v", str)
	}

	u = Url("", ":8081", "/aws/costs/v1/monthly")
	str = u.String()
	if str != "http://localhost:8081/aws/costs/v1/monthly/" {
		t.Errorf("url mapping failed: %v", str)
	}

}
