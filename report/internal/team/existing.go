package team

import (
	"context"
	"log/slog"

	"github.com/ministryofjustice/opg-reports/report/config"
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
// // interface: ImporterExistingCommand
func Existing(ctx context.Context, log *slog.Logger, conf *config.Config, service *opgmetadata.Service) (err error) {
	var (
		raw     []map[string]interface{}
		list    []*Team
		store   *sqldb.Repository[*Team]
		inserts []*sqldb.BoundStatement = []*sqldb.BoundStatement{}
		sw                              = utils.Stopwatch()
	)
	defer service.Close()
	// timer
	sw.Start()
	log = log.With("operation", "Existing", "service", "team")
	log.Info("[team] starting existing records import ...")

	log.Debug("getting teams ...")
	raw, err = service.GetAllTeams()
	if err != nil {
		return
	}
	// now we have raw team data, we need to create a team store and service
	// convert the maps into structs and import to the database
	log.Debug("covnerting to Team ...")
	// convert raw to teams
	list = []*Team{}
	err = utils.Convert(raw, &list)

	// sqldb
	log.Debug("creating datastore ...")
	store, err = sqldb.New[*Team](ctx, log, conf)
	if err != nil {
		return
	}

	log.Debug("generating bound statements ...")
	for _, row := range list {
		inserts = append(inserts, &sqldb.BoundStatement{Statement: stmtImport, Data: row})
	}
	log.Debug("running inserts ...")
	err = store.Insert(inserts...)
	if err != nil {
		return
	}
	log.With(
		"seconds", sw.Stop().Seconds(),
		"inserted", len(inserts)).
		Info("[team] existing records imported.")

	return
}
