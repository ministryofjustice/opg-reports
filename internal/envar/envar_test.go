package envar_test

import (
	"os"
	"testing"

	"github.com/ministryofjustice/opg-reports/pkg/envar"
)

func TestGet(t *testing.T) {

	if envar.Get("PATH", "123") != os.Getenv("PATH") {
		t.Errorf("path mismatch")
	}

	if envar.Get("IFTHISEXISTSTHENTHETESTWILLFAILBUTITSHOULDNOT", "123") != "123" {
		t.Errorf("default failed")
	}

}
