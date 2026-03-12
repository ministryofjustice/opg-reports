package arguments

import (
	"opg-reports/report/packages/types"
	"testing"
)

var (
	_ types.Versioner       = &Versions{}
	_ types.DBer            = &DB{}
	_ types.Hoster          = &host{}
	_ types.ServerConfigure = &Api{}
)

func TestPackagesArgumentsDefaults(t *testing.T) {
	var api = Default[*Api]()

	// check some of the defaults...
	if api.DB.Driver != "sqlite3" {
		t.Errorf("incorrect driver name")
	}
	if api.Version.Version != "0.0.1" {
		t.Errorf("inccorect default semver")
	}
	if api.Info.Name != "api" {
		t.Errorf("incorrect default server name")
	}
}
