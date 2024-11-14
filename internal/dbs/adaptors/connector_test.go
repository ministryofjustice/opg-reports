package adaptors_test

import (
	"testing"

	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/dbs/adaptors"
)

// check interface
var _ dbs.Connector = &adaptors.Connection{}

func TestAdaptorsConnector(t *testing.T) {
	var conn = &adaptors.Connection{Driver: "test", Path: "test.db", Parameters: "?yes=1"}
	var connStr = "test.db?yes=1"
	if conn.String() != connStr {
		t.Errorf("GetConnectionString returned unexpected result: [%s]", conn.String())
	}
	if conn.DriverName() != "test" {
		t.Errorf("GetDriverName returned unexpected result: [%s]", conn.DriverName())
	}
}
