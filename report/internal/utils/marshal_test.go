package utils_test

import (
	"fmt"
	"testing"

	"opg-reports/report/internal/utils"
)

type testItem struct {
	Name string `json:"name,omitempty"`
}

func TestMarshal(t *testing.T) {
	var (
		err    error
		result []byte
		source *testItem = &testItem{Name: "test-item"}
		dest   string    = `{
  "name": "test-item"
}`
	)
	// make sure marshal on a valid struct matches the expected outcome
	result, err = utils.Marshal(source)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if string(result) != dest {
		t.Errorf("marshaled version not as expected: expected [%s] actual [%s]", dest, string(result))
	}
	// wrap in must
	result = utils.MustMarshal(source)
	if string(result) != dest {
		t.Errorf("mustmarshaled version not as expected: expected [%s] actual [%s]", dest, string(result))
	}

	// now do a direct check with MarshalStr
	str := utils.MarshalStr(source)
	if str != dest {
		t.Errorf("marshalStr version not as expected: expected [%s] actual [%s]", dest, str)
	}

	// test marshal an invalid value and make sure if fails
	result, err = utils.Marshal(func() { fmt.Print("here") })
	if err == nil {
		t.Errorf("expected an error when marshaling a func")
	}
	// check must marshal returns empty on a error
	result = utils.MustMarshal(func() { fmt.Print("here") })
	if len(result) != 0 {
		t.Errorf("expected empty result when must marshaling a func")
	}
}

func TestUnmarshal(t *testing.T) {
	var (
		err    error
		result *testItem = &testItem{}
		dest   *testItem = &testItem{Name: "test-item"}
		source []byte    = []byte(`{"name": "test-item"}`)
	)

	// unmarshal valid data into a struct
	err = utils.Unmarshal(source, &result)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if result.Name != dest.Name {
		t.Errorf("result doesnt match expected values: expected [%s] actual [%s]", dest, result)
	}

	// unmarshal mis-matching, but valid json
	result = &testItem{}
	err = utils.Unmarshal([]byte(`{"foo":"bar"}`), &result)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if result.Name != "" {
		t.Errorf("expected the result data to be empty: [%s]", result)
	}

}
