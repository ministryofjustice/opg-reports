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

// Connection generates the db connection for the path and driver; returns db object
//
// If the folder path does not exist, it will be created.
func Connection(ctx context.Context, log *slog.Logger, driver string, connectionStr string) (db *sqlx.DB, err error) {
	var parentDir string = filepath.Dir(connectionStr)
	// add log info
	log = log.With("package", "db.dbconnection", "func", "Connection", "parentDir", parentDir, "driver", driver, "connectionStr", connectionStr)

	log.Debug("starting ...")
	os.MkdirAll(parentDir, os.ModePerm)

	db, err = sqlx.ConnectContext(ctx, driver, connectionStr)
	if err != nil {
		log.Error("error with connection", "err", err.Error())
		err = errors.Join(ErrFailedToConnect, err)
		return
	}
	log.Debug("completed.")
	return
}
