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
	r.SetResults("test")

	if r.GetResults() != "test" {
		t.Errorf("failed with string res")
	}

	r.Result = []string{"test1", "test2"}
	if len(r.GetResults().([]string)) != 2 {
		t.Errorf("failed with strings")
	}

	r.SetResults([]int{1, 2, 3})
	if len(r.GetResults().([]int)) != 3 {
		t.Errorf("failed with []int")
	}
}

func TestServicesApiAwsCostMonthlyResponseErrs(t *testing.T) {
	r := &ApiResponse{}
	r.SetErrors([]error{errors.New("test error")})

	if len(r.GetErrors()) != len(r.Errors) && len(r.GetErrors()) != 1 {
		t.Errorf("errors not set")
	}

	r.AddError(errors.New("test"))
	if len(r.GetErrors()) != 2 {
		t.Errorf("error add failed")
	}
}
