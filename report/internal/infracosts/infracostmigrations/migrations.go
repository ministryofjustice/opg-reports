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
CREATE TABLE IF NOT EXISTS aws_costs (
	id INTEGER PRIMARY KEY,
	created_at TEXT NOT NULL DEFAULT (strftime('%FT%TZ', 'now') ),
	region TEXT DEFAULT "NoRegion" NOT NULL,
	service TEXT NOT NULL,
	date TEXT NOT NULL,
	cost TEXT NOT NULL,
	aws_account_id TEXT,
	UNIQUE (aws_account_id,date,region,service)
) STRICT;`

const idx_date string = `CREATE INDEX IF NOT EXISTS aws_costs_date_idx ON aws_costs(date);`
const idx_date_account string = `CREATE INDEX IF NOT EXISTS aws_costs_date_account_idx ON aws_costs(date, aws_account_id);`
const idx_merged string = `CREATE INDEX IF NOT EXISTS aws_costs_unique_idx ON aws_costs(aws_account_id,date,region,service);`

var migrations []string = []string{
	table,
	idx_date,
	idx_date_account,
	idx_merged,
}

// Migrate runs all the create / index setup calls and executes them against the
// db connections passed
func Migrate(ctx context.Context, log *slog.Logger, db *sqlx.DB) (err error) {

	log = log.With("package", "infracosts.ce", "func", "Migrate")
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
