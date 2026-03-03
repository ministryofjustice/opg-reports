package dbx

import (
	"context"
	"database/sql"
	"log/slog"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/cnv"
	"opg-reports/report/package/conn"
)

type InsertArgs struct {
	DB     string `json:"db"`     // database path
	Driver string `json:"driver"` // database driver
	Params string `json:"params"` // database connection params
}

func Insert[T any](ctx context.Context, stmt string, records []T, in *InsertArgs) (err error) {
	var (
		db  *sql.DB
		log *slog.Logger = cntxt.GetLogger(ctx).With("package", "dbx", "func", "Insert")
	)
	db, err = sql.Open(in.Driver, conn.SqlitePath(in.DB, in.Params))
	if err != nil {
		log.Error("error connecting to database", "err", err.Error())
		return
	}
	defer db.Close()

	for _, model := range records {
		// convert to map
		row := map[string]interface{}{}
		if err = cnv.Convert(model, &row); err != nil {
			return
		}
		// use bind
		bound, args, e := Bind(ctx, stmt, row)
		if e != nil {
			return e
		}
		_, err = db.ExecContext(ctx, bound, args...)
		if err != nil {
			return
		}

	}

	return
}
