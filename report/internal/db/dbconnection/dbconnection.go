package dbconnection

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var ErrFailedToConnect = errors.New("db connection failed with error.")

var driverParams map[string]string = map[string]string{
	"sqlite3": "?_journal=WAL&_busy_timeout=5000&_vacuum=incremental&_synchronous=NORMAL&_cache_size=1000000000",
}

// Connection generates the db connection for the path and driver; returns db object
//
// If the folder path does not exist, it will be created.
func Connection(ctx context.Context, log *slog.Logger, driver string, connectionStr string) (db *sqlx.DB, err error) {
	var (
		parentDir string       = filepath.Dir(connectionStr)
		param     string       = driverParams[driver]
		lg        *slog.Logger = log.With("func", "dbconnection.Connection", "parentDir", parentDir, "driver", driver, "connectionStr", connectionStr)
	)

	// remove and re-add the params
	connectionStr = strings.ReplaceAll(connectionStr, param, "") + param

	lg.Debug("starting ...")
	os.MkdirAll(parentDir, os.ModePerm)
	db, err = sqlx.ConnectContext(ctx, driver, connectionStr)
	if err != nil {
		lg.Error("error with connection", "err", err.Error())
		err = errors.Join(ErrFailedToConnect, err)
		return
	}
	lg.Debug("complete.")
	return
}
