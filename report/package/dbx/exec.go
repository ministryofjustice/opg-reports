package dbx

import (
	"context"
	"database/sql"
	"log/slog"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/conn"
)

type ExecArgs struct {
	DB     string `json:"db"`     // database path
	Driver string `json:"driver"` // database driver
	Params string `json:"params"` // database connection params
}

func Exec(ctx context.Context, stmt string, in *ExecArgs, args ...any) (err error) {
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

	_, err = db.ExecContext(ctx, stmt, args...)
	if err != nil {
		log.Error("error running statement", "err", err.Error())
		return
	}

	return
}
