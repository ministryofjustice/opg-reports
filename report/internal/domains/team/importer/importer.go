package importer

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/domains/team/types"
	"opg-reports/report/packages/args"
	"opg-reports/report/packages/files"
	"opg-reports/report/packages/logger"
)

const InsertStatement string = `
INSERT INTO teams (
	name
) VALUES (
	lower(:name)
) ON CONFLICT (name)
 	DO UPDATE SET name=excluded.name
RETURNING name
;
`

// dummy used to comply with interface
type client interface{}

func Get(ctx context.Context, client *client, opts *args.Import, previous ...*types.Team) (found []*types.Team, err error) {
	var log *slog.Logger
	ctx, log = logger.Get(ctx)

	log.Info("getting team data from local file ...")
	found = []*types.Team{}
	err = files.ReadJSON(ctx, opts.File.Path, &found)
	if err != nil {
		log.Error("failed to read in source file", "err", err.Error())
		return
	}

	log.Info("team data completed.")

	return
}

// Filter doesnt apply for this data set
func Filter(ctx context.Context, items []*types.Team, filters *args.Filters) (included []*types.Team) {
	return items
}

// Transform converts the original data into record for local database insertion
func Transform(ctx context.Context, data []*types.Team, opts *args.Import) (results []*types.Team, err error) {
	results = data
	return
}
