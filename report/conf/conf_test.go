package conf

import (
	"os"
	"testing"
)

var confName = "default"

// TestConfLoadViaSetup checks the db config is loaded with expected
// db values from a direct setup() call
func TestConfLoadViaSetup(t *testing.T) {

	vp, conf, err := setup(confName)

	if err != nil {
		t.Errorf("unexpected error getting config")
		t.FailNow()
	}

	if conf.DB.Driver != "sqlite3" {
		t.Errorf("setup db config does not match")
	}

	if vp.Get("db.driver") != "sqlite3" {
		t.Errorf("viper get db config did not match default")
	}

}

func TestConfViaNewWithEnvVars(t *testing.T) {

	os.Setenv("DB_DRIVER", "sqlite2")
	conf := New()
	if conf.DB.Driver != "sqlite2" {
		t.Errorf("instance db config does not match")
	}

}
