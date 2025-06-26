package team

import (
	"context"
	"log/slog"
	"time"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/gh"
	"github.com/ministryofjustice/opg-reports/report/internal/opgmetadata"
	"github.com/ministryofjustice/opg-reports/report/internal/sqldb"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

// Import generates new team data from the billing_unit information within the
// opg-metadata published list of accounts.
//
// That is a private repository so you need permissions to read and fetch data to be
// able to download the release asset.
//
// The account.json is parsed and all unique billing_units are converted into team.Team
// entries and inserted into the databse using the team.Service.Import method
//
// // interface: ImporterImportCommand
func Import(ctx context.Context, log *slog.Logger, conf *config.Config) (err error) {
	var (
		gr          *gh.Repository
		metaService *opgmetadata.Service
		raw         []map[string]interface{}
		list        []*Team
		store       *sqldb.Repository[*Team]
		now         string                  = time.Now().UTC().Format(time.RFC3339)
		inserts     []*sqldb.BoundStatement = []*sqldb.BoundStatement{}
	)
	log = log.With("operation", "Import", "service", "team")
	log.Info("running [team] imports ...")

	// to import teams, we create an opgmetadata service and call the getTeams
	// so fetch the gh repository first and then create the opgmeta data service
	gr, err = gh.New(ctx, log, conf)
	if err != nil {
		return
	}
	metaService, err = opgmetadata.NewService(ctx, log, conf, gr)
	if err != nil {
		return
	}
	defer metaService.Close()

	raw, err = metaService.GetAllTeams()
	if err != nil {
		return
	}
	// now we have raw team data, we need to create a team store and service
	// convert the maps into structs and import to the database

	// convert raw to teams
	list = []*Team{}
	err = utils.Convert(raw, &list)
	// sqldb
	store, err = sqldb.New[*Team](ctx, log, conf)
	if err != nil {
		return
	}
	log.Info("importing ...")
	log.Debug("generating bound statements ...")
	for _, row := range list {
		row.CreatedAt = now
		inserts = append(inserts, &sqldb.BoundStatement{Statement: stmtImport, Data: row})
	}
	log.Debug("running insert ...")
	err = store.Insert(inserts...)

	return
}
