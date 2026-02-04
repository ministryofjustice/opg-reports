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
	var ght string = os.Getenv("GITHUB_TOKEN")
	// set dummy values
	os.Setenv("GITHUB_TOKEN", "test-token")
	os.Setenv("DB_DRIVER", "sqlite2")

	cfg := New()

	if cfg.DB.Driver != "sqlite2" {
		t.Errorf("instance db config does not match")
	}
	if cfg.GithubToken != "test-token" {
		t.Errorf("instance github token does not match")
	}

	os.Setenv("DB_DRIVER", "sqlite3")
	os.Setenv("GITHUB_TOKEN", ght)

}
