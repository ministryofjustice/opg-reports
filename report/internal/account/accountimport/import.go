package accountimport

import (
	"context"
	"log/slog"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/dbx"
	"opg-reports/report/package/files"

	_ "github.com/mattn/go-sqlite3"
)

const InsertStatement string = `
INSERT INTO accounts (
	id,
	name,
	label,
	environment,
	team_name
) VALUES (
	:id,
	:name,
	:label,
	:environment,
	:billing_unit
) ON CONFLICT (id) DO UPDATE SET
	name=excluded.name,
	team_name=excluded.team_name,
	label=excluded.label,
	environment=excluded.environment
RETURNING id
;
`

// Model represents a simple, joinless, db row in the team table; used by imports and seeding commands
type Model struct {
	ID          string `json:"id,omitempty"`            // This is the Account ID as a string - they can have leading 0
	Name        string `json:"name,omitempty" `         // account name as used internally
	Label       string `json:"label,omitempty" `        // internal label
	Environment string `json:"environment,omitempty" `  // environment type
	TeamName    string `json:"billing_unit,omitempty" ` // team associated with the account; uses builling_unit due to the source data in opg-metadata
}

type Args struct {
	DB      string `json:"db"`       // database path
	Driver  string `json:"driver"`   // database driver
	Params  string `json:"params"`   // database connection params
	SrcFile string `json:"src-file"` // src file to import from
}

func Import(ctx context.Context, in *Args) (err error) {
	var (
		accounts []*Model     = []*Model{}
		log      *slog.Logger = cntxt.GetLogger(ctx).With("package", "accountimport", "func", "Import")
	)
	log.Info("starting ...", "db", in.DB, "file", in.SrcFile)

	err = files.ReadJSON(ctx, in.SrcFile, &accounts)
	if err != nil {
		log.Error("failed to read in source file", "err", err.Error())
		return
	}

	// now write to db
	err = dbx.Insert(ctx, InsertStatement, accounts, &dbx.InsertArgs{
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
