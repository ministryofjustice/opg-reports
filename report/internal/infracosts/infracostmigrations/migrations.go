package infracostmigrations

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
CREATE TABLE IF NOT EXISTS infracosts (
	id INTEGER PRIMARY KEY,
	created_at TEXT NOT NULL DEFAULT (strftime('%FT%TZ', 'now') ),
	vendor TEXT NOT NULL DEFAULT 'aws',
	region TEXT DEFAULT "NoRegion" NOT NULL,
	service TEXT NOT NULL,
	date TEXT NOT NULL,
	cost TEXT NOT NULL,
	account_id TEXT,
	UNIQUE (account_id,date,region,service)
) STRICT;`

const idx_date string = `CREATE INDEX IF NOT EXISTS infracosts_date_idx ON infracosts(date);`
const idx_date_account string = `CREATE INDEX IF NOT EXISTS infracosts_date_account_idx ON infracosts(date, aws_account_id);`
const idx_merged string = `CREATE INDEX IF NOT EXISTS infracosts_unique_idx ON infracosts(aws_account_id,date,region,service);`

var migrations []string = []string{
	table,
	idx_date,
	idx_date_account,
	idx_merged,
}

// Migrate runs all the create / index setup calls and executes them against the
// db connections passed
func Migrate(ctx context.Context, log *slog.Logger, db *sqlx.DB) (err error) {

	log = log.With("package", "infracosts", "func", "Migrate")
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
