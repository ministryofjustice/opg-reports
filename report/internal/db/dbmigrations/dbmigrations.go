package dbmigrations

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/db/dbexec"
	"opg-reports/report/internal/db/dbinserts"
	"opg-reports/report/internal/db/dbselects"
	"opg-reports/report/internal/db/dbstmts"

	"github.com/jmoiron/sqlx"
)

// Migration contains migration data to run
type Migration struct {
	ID   string `json:"id,omitempty" db:"id"`
	Name string `json:"name" db:"name"`
	SQL  string `json:"sql" db:"sql"`
}

// empty is a dummy struct for the select
type empty struct{}

const selectStmt string = `SELECT id, name FROM migrations ORDER BY id ASC`
const insertStmt string = `INSERT INTO migrations (name) VALUES (:name) RETURNING id`

// errors
var ErrMigrationExecFailed = errors.New("migration statement failed with error")

// Migrate runs over all the configured statements, compares to what is in the database already and
// runs anythign that has not been handled
//
// Migrations that are excuted are added to the database via defer func call so trigger on
// exiting
//
// If `migrations` param is empty, then the default MIGRATIONS is used
func Migrate(ctx context.Context, log *slog.Logger, db *sqlx.DB, migrations ...*Migration) (err error) {

	var (
		allMigrations   []*Migration                        = MIGRATIONS // use this by default
		migrationsToRun []*Migration                        = allMigrations
		done            []*Migration                        = []*Migration{}
		lg              *slog.Logger                        = log.With("func", "dbmigrations.Migrate")
		selector        *dbstmts.Select[*empty, *Migration] = &dbstmts.Select[*empty, *Migration]{
			Statement: selectStmt,
			Data:      &empty{},
		}
	)

	lg.Debug("starting ...")

	lg.Debug("running migration table create ....")
	_, err = dbexec.Exec(ctx, log, db, dbstmts.Statement(create_migration))
	if err != nil {
		return
	}

	// if migrations are passed, use those
	if len(migrations) > 0 {
		migrationsToRun = migrations
	}
	// first, find all migrations that have happened
	lg.Info("getting existing migrations ...")
	err = dbselects.Select(ctx, log, db, selector)

	// insert migrations when done with the function, should trigger before error returns
	defer func() {
		insertMigrationData(ctx, log, db, done)
	}()

	// get the migrations we need to run
	migrationsToRun = migrationsToExec(selector.Returned, allMigrations)
	// now run the migrations
	lg.With("count", len(migrationsToRun)).Debug("running migrataions ...")
	for _, toRun := range migrationsToRun {
		_, err = dbexec.Exec(ctx, log, db, dbstmts.Statement(toRun.SQL))
		if err != nil {
			lg.Error("error with migration exec", "err", err.Error())
			err = errors.Join(ErrMigrationExecFailed, err)
			return
		} else {
			done = append(done, toRun)
		}
	}
	lg.Debug("complete.")
	return

}

// insertMigrationData insert migration info based on set of done migrations
func insertMigrationData(ctx context.Context, log *slog.Logger, db *sqlx.DB, done []*Migration) (err error) {
	var stmts = []*dbstmts.Insert[*Migration, int]{}
	for _, item := range done {
		stmts = append(stmts, &dbstmts.Insert[*Migration, int]{
			Statement: insertStmt,
			Data:      &Migration{Name: item.Name},
		})
	}
	err = dbinserts.Insert(ctx, log, db, stmts...)
	return
}

// migrationsToExec finds the migrations need to run
func migrationsToExec(returned []*Migration, all []*Migration) (migrationsToRun []*Migration) {
	migrationsToRun = []*Migration{}
	// check each migratrion and see if they are listed in the returned data
	for _, migration := range all {
		var add = true
		for _, row := range returned {
			if row.Name == migration.Name {
				add = false
			}
		}
		if add {
			migrationsToRun = append(migrationsToRun, migration)
		}
	}
	return
}
