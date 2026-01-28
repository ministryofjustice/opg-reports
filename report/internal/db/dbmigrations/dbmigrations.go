package dbmigrations

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/db/dbexec"
	"opg-reports/report/internal/db/dbinserts"
	"opg-reports/report/internal/db/dbselects"
	"opg-reports/report/internal/db/dbstatements"

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
func Migrate(ctx context.Context, log *slog.Logger, db *sqlx.DB, migrations ...Migration) (err error) {

	var (
		allMigrations   []*Migration                                      = MIGRATIONS // use this by default
		migrationsToRun []*Migration                                      = []*Migration{}
		done            []*Migration                                      = []*Migration{}
		selector        *dbstatements.SelectStatement[*empty, *Migration] = &dbstatements.SelectStatement[*empty, *Migration]{
			Statement: selectStmt,
			Data:      &empty{},
		}
	)

	log = log.With("package", "dbmigration", "func", "Migrate")
	log.Debug("starting ...")

	// if migrations are passed, use those
	if len(migrations) > 0 {
		allMigrations = migrationsToRun
	}

	// first, find all migrations that have happened
	log.Info("getting existing migrations ...")
	err = dbselects.Select(ctx, log, db, selector)
	// if we get an error about missing table, that means no migration table
	// is present so we need to run all migrations
	if errors.Is(err, dbselects.ErrMissingTable) {
		log.Debug("no migration table found, so will run all migrations.")
		err = nil
	} else if err != nil {
		log.Error("migration selection failed", "err", err.Error())
		return
	}

	// insert migrations when done with the function, should trigger before error returns
	defer func() {
		insertMigrationData(ctx, log, db, done)
	}()
	// get the migrations we need to run
	migrationsToRun = migrationsToExec(selector.Returned, allMigrations)
	// now run the migrations
	log.With("count", len(migrationsToRun)).Info("running migrataions ...")
	for _, toRun := range migrationsToRun {
		_, err = dbexec.Exec(ctx, log, db, dbstatements.Statement(toRun.SQL))
		if err != nil {
			log.Error("error with migration exec", "err", err.Error())
			err = errors.Join(ErrMigrationExecFailed, err)
			return
		} else {
			done = append(done, toRun)
		}
	}

	log.Debug("complete")
	return

}

// insertMigrationData insert migration info based on set of done migrations
func insertMigrationData(ctx context.Context, log *slog.Logger, db *sqlx.DB, done []*Migration) (err error) {
	var stmts = []*dbstatements.InsertStatement[*Migration, int]{}
	for _, item := range done {
		stmts = append(stmts, &dbstatements.InsertStatement[*Migration, int]{
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
