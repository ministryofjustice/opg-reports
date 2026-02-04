package dbconnection

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var ErrFailedToConnect = errors.New("db connection failed with error.")

// Connection generates the db connection for the path and driver; returns db object
//
// If the folder path does not exist, it will be created.
func Connection(ctx context.Context, log *slog.Logger, driver string, connectionStr string) (db *sqlx.DB, err error) {
	var parentDir string = filepath.Dir(connectionStr)
	var lg *slog.Logger = log.With("func", "dbconnection.Connection", "parentDir", parentDir, "driver", driver, "connectionStr", connectionStr)

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
