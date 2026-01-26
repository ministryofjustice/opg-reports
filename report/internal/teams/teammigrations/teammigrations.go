package teammigrations

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/db/dbexec"
	"opg-reports/report/internal/db/dbstatements"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const table string = `
CREATE TABLE IF NOT EXISTS teams (
	name TEXT PRIMARY KEY,
	created_at TEXT NOT NULL DEFAULT (strftime('%FT%TZ', 'now') )
) STRICT;`

var migrations []string = []string{
	table,
}

// Migrate runs all the create / index setup calls and executes them against the
// db connections passed
func Migrate(ctx context.Context, log *slog.Logger, db *sqlx.DB) (err error) {

	log = log.With("package", "teams", "func", "Migrate")
	log.Debug("starting ...")

	for _, stmt := range migrations {
		_, err = dbexec.Exec(ctx, log, db, dbstatements.Statement(stmt))
		if err != nil {
			log.Error("error with migration exec", "err", err.Error())
			err = errors.Join(ErrMigrationExecFailed, err)
			return
		}
	}

	log.Debug("complete")
	return

}
