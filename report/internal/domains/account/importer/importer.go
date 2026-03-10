package importer

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/domains/account/types"
	"opg-reports/report/packages/args"
	"opg-reports/report/packages/files"
	"opg-reports/report/packages/logger"
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
	lower(:billing_unit)
) ON CONFLICT (id) DO UPDATE SET
	name=excluded.name,
	team_name=excluded.team_name,
	label=excluded.label,
	environment=excluded.environment
RETURNING id
;
`

// dummy used to comply with interface
type client interface{}

func Get(ctx context.Context, client *client, opts *args.Import, previous ...*types.ImportAccount) (found []*types.ImportAccount, err error) {
	var log *slog.Logger
	ctx, log = logger.Get(ctx)

	log.Info("getting account data from local file ...")
	found = []*types.ImportAccount{}
	err = files.ReadJSON(ctx, opts.File.Path, &found)
	if err != nil {
		log.Error("failed to read in source file", "err", err.Error())
		return
	}

	log.Info("account data completed.")

	return
}

// Filter doesnt apply for this data set
func Filter(ctx context.Context, items []*types.ImportAccount, filters *args.Filters) (included []*types.ImportAccount) {
	return items
}

// Transform converts the original data into record for local database insertion
func Transform(ctx context.Context, data []*types.ImportAccount, opts *args.Import) (results []*types.ImportAccount, err error) {
	results = data
	return
}
