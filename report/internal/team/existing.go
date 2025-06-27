package team

import (
	"context"
	"log/slog"
	"slices"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/interfaces"
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
// T = *Team
//
// interface: ImporterExistingCommand
func Existing[T interfaces.Model](ctx context.Context, log *slog.Logger, conf *config.Config, service interfaces.MetadataService[T]) (err error) {
	var (
		data     []T
		store    *sqldb.Repository[*TeamImport]
		inserts  []*sqldb.BoundStatement = []*sqldb.BoundStatement{}
		names    []string                = []string{}
		teams    []*TeamImport           = []*TeamImport{}
		owner    string                  = conf.Github.Organisation
		repo     string                  = conf.Github.Metadata.Repository
		asset    string                  = conf.Github.Metadata.Asset
		dataFile string                  = "accounts.json"
		sw                               = utils.Stopwatch()
	)
	defer service.Close()
	// timer
	sw.Start()
	log = log.With("operation", "Existing", "service", "team")
	log.Info("[team] starting existing records import ...")

	data, err = service.DownloadAndReturn(owner, repo, asset, dataFile)
	if err != nil || len(data) <= 0 {
		return
	}
	log.Debug("[team] converting to local format ...")
	// convert data list
	err = utils.Convert(data, &teams)
	if err != nil || len(teams) <= 0 {
		return
	}
	log.Debug("[team] generating unique list of team names ...")
	// now filter this down to unique billing_unit values
	for _, item := range teams {
		names = append(names, item.Name)
	}
	slices.Sort(names)
	names = slices.Compact(names)

	for _, nm := range names {
		inserts = append(inserts, &sqldb.BoundStatement{Statement: stmtImport, Data: &TeamImport{Name: nm}})
	}
	log.With("count", len(inserts)).Debug("[team] records to insert ...")

	log.Debug("[team] creating writer store for insert ...")
	store, err = sqldb.New[*TeamImport](ctx, log, conf)
	if err != nil {
		return
	}

	log.Debug("[team] running insert ...")
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
