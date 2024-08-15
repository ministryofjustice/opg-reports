package env

import (
	"os"
	"testing"

	"github.com/ministryofjustice/opg-reports/shared/logger"
)

func TestGet(t *testing.T) {
	logger.LogSetup()

	if Get("PATH", "123") != os.Getenv("PATH") {
		t.Errorf("path mismatch")
	}

	if Get("IFTHISEXISTSTHENTHETESTWILLFAILBUTITSHOULDNOT", "123") != "123" {
		t.Errorf("default failed")
	}

}
