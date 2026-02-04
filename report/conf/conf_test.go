package conf

import (
	"os"
	"testing"
)

var confName = "default"

// TestConfLoadViaSetup checks the db config is loaded with expected
// db values from a direct setup() call
func TestConfLoadViaSetup(t *testing.T) {

	_, cfg, err := setup()

	if err != nil {
		t.Errorf("unexpected error getting config: %s", err.Error())
		t.FailNow()
	}

	if cfg.DB.Driver != "sqlite3" {
		t.Errorf("setup db config does not match")
	}

}

func TestConfViaNewWithEnvVars(t *testing.T) {

	os.Setenv("DB_DRIVER", "sqlite2")
	cfg := New()

	if cfg.DB.Driver != "sqlite2" {
		t.Errorf("instance db config does not match")
	}
	os.Setenv("DB_DRIVER", "sqlite3")

}
