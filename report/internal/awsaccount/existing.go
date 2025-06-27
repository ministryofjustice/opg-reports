package awsaccount

import (
	"context"
	"log/slog"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/opgmetadata"
	"github.com/ministryofjustice/opg-reports/report/internal/sqldb"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

// Existing generates new aws account data from the accounts.json information within the
// opg-metadata published data. The accounts.aws.json is parsed and converted to db entries
//
// That is a private repository so you need permissions to read and fetch data to be
// able to download the release asset.
//
// Example account from the opg-metadata source file:
//
//	[{
//		"id": "500000067891",
//		"name": "My production",
//		"billing_unit": "Team A",
//		"label": "prod",
//		"environment": "production",
//		"type": "aws",
//		"uptime_tracking": true
//	}]
//
// interface: ImporterExistingCommand
func Existing(ctx context.Context, log *slog.Logger, conf *config.Config, service *opgmetadata.Service) (err error) {
	var (
		rawAccounts []map[string]interface{}
		list        []*awsAccountImport
		store       *sqldb.Repository[*awsAccountImport]
		inserts     []*sqldb.BoundStatement = []*sqldb.BoundStatement{}
		sw                                  = utils.Stopwatch()
	)
	defer service.Close()
	// timer
	sw.Start()

	log = log.With("operation", "Existing", "service", "awsaccount")
	log.Info("[awsaccount] starting existing records import ...")

	// get just the aws accounts
	log.Debug("getting accounts ...")
	rawAccounts, err = service.GetAllAwsAccounts()
	if err != nil {
		return
	}
	// convert to db model
	log.Debug("converting to model ...")
	list = []*awsAccountImport{}
	err = utils.Convert(rawAccounts, &list)
	if err != nil {
		return
	}
	// sqldb setup
	log.Debug("creating datastore ...")
	store, err = sqldb.New[*awsAccountImport](ctx, log, conf)
	if err != nil {
		return
	}
	log.Debug("generating bound statements ...")
	for _, row := range list {
		if row.Environment == "" {
			row.Environment = "production"
		}
		inserts = append(inserts, &sqldb.BoundStatement{Statement: stmtImport, Data: row})
	}

	log.Debug("running insert ...")
	err = store.Insert(inserts...)
	if err != nil {
		return
	}

	log.With(
		"seconds", sw.Stop().Seconds(),
		"inserted", len(inserts)).
		Info("[awsaccount] existing records imported.")

	return
}
