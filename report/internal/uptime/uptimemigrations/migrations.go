package uptimemigrations

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
CREATE TABLE IF NOT EXISTS aws_uptime (
	id INTEGER PRIMARY KEY,
	created_at TEXT NOT NULL DEFAULT (strftime('%FT%TZ', 'now') ),
	date TEXT NOT NULL,
	aws_account_id TEXT,
	average TEXT NOT NULL,
	granularity TEXT NOT NULL,
	UNIQUE (aws_account_id,date)
) STRICT;
 `

const idx_date string = `CREATE INDEX IF NOT EXISTS aws_uptime_date_idx ON aws_uptime(date);`
const idx_date_account string = `CREATE INDEX IF NOT EXISTS aws_uptime_account_date_idx ON aws_uptime(aws_account_id,date);`

var migrations []string = []string{
	table,
	idx_date,
	idx_date_account,
}

// Migrate runs all the create / index setup calls and executes them against the
// db connections passed
func Migrate(ctx context.Context, log *slog.Logger, db *sqlx.DB) (err error) {

	log = log.With("package", "uptimetime", "func", "Migrate")
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
