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
		lg      = logger.New("debug", "text")
		driver  = "sqlite3"
		connStr = fmt.Sprintf("%s/%s", dir, "test.db")
	)

	_, err = Connection(ctx, lg, driver, connStr)
	if err != nil {
		t.Errorf("unexpected connection issue:\n%v", err.Error())
	}

}
