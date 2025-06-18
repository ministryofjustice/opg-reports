package utils_test

import (
	"os"
	"testing"

	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

func TestGet(t *testing.T) {

	if utils.GetEnvVar("PATH", "123") != os.Getenv("PATH") {
		t.Errorf("path mismatch")
	}

	if utils.GetEnvVar("IFTHISEXISTSTHENTHETESTWILLFAILBUTITSHOULDNOT", "123") != "123" {
		t.Errorf("default failed")
	}

}
