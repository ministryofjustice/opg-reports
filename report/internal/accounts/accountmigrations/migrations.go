package accountmigrations

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
CREATE TABLE IF NOT EXISTS accounts (
	id TEXT PRIMARY KEY,
	created_at TEXT NOT NULL DEFAULT (strftime('%FT%TZ', 'now') ),
	vendor TEXT NOT NULL DEFAULT 'aws',
	name TEXT NOT NULL,
	label TEXT NOT NULL,
	environment TEXT NOT NULL DEFAULT "production",
	uptime_tracking TEXT NOT NULL DEFAULT "false",
	team_name TEXT NOT NULL DEFAULT "ORG"
) WITHOUT ROWID;
 `
const idx_account_id string = `CREATE INDEX IF NOT EXISTS aws_accounts_id_idx ON accounts(id);`

var migrations []string = []string{
	table,
	idx_account_id,
}

// Migrate runs all the create / index setup calls and executes them against the
// db connections passed
func Migrate(ctx context.Context, log *slog.Logger, db *sqlx.DB) (err error) {

	log = log.With("package", "accounts", "func", "Migrate")
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
