package env

import (
	"os"
	"testing"
)

func TestGet(t *testing.T) {

	if Get("PATH", "123") != os.Getenv("PATH") {
		t.Errorf("path mismatch")
	}

	if Get("IFTHISEXISTSTHENTHETESTWILLFAILBUTITSHOULDNOT", "123") != "123" {
		t.Errorf("default failed")
	}

}
