package convert

import (
	"opg-reports/app/internal/fmtx"
	"testing"
)

type testStructA struct {
	ID   string
	Name string
}

func TestConvertConvertStructToMap(t *testing.T) {
	model := &testStructA{ID: "1001", Name: "A"}
	mapped := map[string]string{}

	err := Convert(model, &mapped)
	if err != nil {
		t.Errorf("unexpected error converting: %v\n", err.Error())
	}

	if id, ok := mapped["ID"]; !ok || id != "1001" {
		t.Errorf("ID not as expected:\n")
		fmtx.Printj(mapped)
	}

	if nm, ok := mapped["Name"]; !ok || nm != "A" {
		t.Errorf("Name not as expected:\n")
		fmtx.Printj(mapped)
	}

}
