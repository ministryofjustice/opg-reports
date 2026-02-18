package teamimport

import (
	"context"
	"log/slog"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/dbx"
	"opg-reports/report/package/files"

	_ "github.com/mattn/go-sqlite3"
)

// TeamModel represents a simple, joinless, db row in the team table; used by imports and seeding commands
type TeamModel struct {
	Name string `json:"name" db:"name"`
}

type Args struct {
	DB     string `json:"db"`     // database path
	Driver string `json:"driver"` // database driver
	Params string `json:"params"` // database connection params
	File   string `json:""`
}

func Import(ctx context.Context, in *Args) (err error) {
	var (
		teams []*TeamModel = []*TeamModel{}
		log   *slog.Logger = cntxt.GetLogger(ctx).With("package", "teamimport", "func", "Import")
	)
	log.Info("starting ...", "db", in.DB, "file", in.File)

	err = files.ReadJSON(ctx, in.File, &teams)
	if err != nil {
		log.Error("failed to read in source file", "err", err.Error())
		return
	}

	// now write to db
	err = dbx.Insert(ctx, importStatement, teams, &dbx.InsertArgs{
		DB:     in.DB,
		Driver: in.Driver,
		Params: in.Params,
	})
	if err != nil {
		log.Error("error write data during import", "err", err.Error())
		return
	}

	log.Info("complete.")
	return
}
