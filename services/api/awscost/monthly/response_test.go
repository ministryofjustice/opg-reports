package monthly

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestServicesApiAwsCostMonthlyResponseBodyMarshaling(t *testing.T) {
	r := &ApiResponse{}
	r.Start()
	r.End()
	body := r.Body()

	u := &ApiResponse{}
	json.Unmarshal(body, u)

	if u.RequestStart != r.RequestStart || u.RequestEnd != r.RequestEnd {
		t.Errorf("unmarshall error")
	}
}

func TestServicesApiAwsCostMonthlyResponseResults(t *testing.T) {

	r := &ApiResponse{}
	r.Set("test", 200, nil)

	if r.Results() != "test" {
		t.Errorf("failed with string res")
	}

	r.Res = []string{"test1", "test2"}
	if len(r.Results().([]string)) != 2 {
		t.Errorf("failed with strings")
	}

	r.Set([]int{1, 2, 3}, 200, nil)
	if len(r.Results().([]int)) != 3 {
		t.Errorf("failed with []int")
	}
}

func TestServicesApiAwsCostMonthlyResponseErrs(t *testing.T) {
	r := &ApiResponse{}
	r.Set("test", 400, []error{errors.New("test error")})

	if len(r.Errors()) != len(r.Errs) && len(r.Errors()) != 1 {
		t.Errorf("errors not set")
	}

}
