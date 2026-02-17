package queryx

import (
	"context"
	"database/sql"
	"log/slog"
	"opg-reports/report/package/bind"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/conn"
)

type scanFunc func(rows *sql.Rows) error

type Input struct {
	DB      string
	Driver  string
	Params  string
	BindMap map[string]interface{}
	ScanF   scanFunc // used in rows.Next() to process a row
}

// Select uses the stmt provided and the args (in.BindMap) to generate
// a bound sql (`:name` => `?` etc) and then uses the ScanF function passed
// to load each row into a struct.
//
// The ScanF function handles the model creation etc to is normally created
// inline within the func for scoping of types etc.
func Select(ctx context.Context, stmt string, in *Input) (err error) {
	var (
		db   *sql.DB
		args []interface{}
		rows *sql.Rows
		log  *slog.Logger = cntxt.GetLogger(ctx).With("package", "queryx", "func", "Select")
	)
	// create the db connection
	db, err = sql.Open(in.Driver, conn.SqlitePath(in.DB, in.Params))
	if err != nil {
		return
	}
	defer db.Close()

	// generate the bound statement
	stmt, args, err = bind.Bind(stmt, in.BindMap)
	if err != nil {
		log.Error("error in bind", "err", err.Error())
		return
	}
	// setup row scanning
	rows, err = db.QueryContext(ctx, stmt, args...)
	if err != nil {
		log.Error("error in query", "err", err.Error())
		return
	}

	defer rows.Close()
	for rows.Next() {
		err = in.ScanF(rows)
		if err != nil {
			log.Error("row scan failed", "err", err.Error())
			return
		}
	}

	return
}
