package dbconnection

import (
	"fmt"
	"opg-reports/report/internal/utils/logger"
	"testing"
)

func TestDBDBConnectionWorking(t *testing.T) {
	var (
		err     error
		dir     = t.TempDir()
		ctx     = t.Context()
		log     = logger.New("error", "text")
		driver  = "sqlite3"
		connStr = fmt.Sprintf("%s/%s", dir, "connection-working.db")
	)

	db, err := Connection(ctx, log, driver, connStr)
	if err != nil {
		t.Errorf("unexpected connection issue:\n%v", err.Error())
	}
	defer db.Close()

}
