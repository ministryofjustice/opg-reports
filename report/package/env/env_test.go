package env

import (
	"fmt"
	"opg-reports/report/package/dump"
	"os"
	"testing"
)

type mock struct {
	TestID   string `json:"test_id"`
	TestName string `json:"test_name"`
}

func TestUtilsEnvOverwriteStruct(t *testing.T) {
	var err error

	os.Setenv("TEST_ID", "100")
	m := &mock{TestName: "A"}
	err = OverwriteStruct(&m)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}
	if m.TestID != "100" || m.TestName != "A" {
		t.Errorf("incorrect value")
		fmt.Println(dump.Any(m))
	}

	os.Setenv("TEST_NAME", "B")
	err = OverwriteStruct(&m)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}
	if m.TestID != "100" || m.TestName != "B" {
		t.Errorf("incorrect value")
		fmt.Println(dump.Any(m))
	}

}
