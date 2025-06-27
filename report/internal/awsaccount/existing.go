package awsaccount

import (
	"context"
	"log/slog"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/interfaces"
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
// T =  *AwsAccountImport
//
// interface: ImporterExistingCommand
func Existing[T interfaces.Model](ctx context.Context, log *slog.Logger, conf *config.Config, service *opgmetadata.Service[T]) (err error) {
	var (
		data     []T
		store    *sqldb.Repository[*AwsAccountImport]
		inserts  []*sqldb.BoundStatement = []*sqldb.BoundStatement{}
		owner    string                  = conf.Github.Organisation
		repo     string                  = conf.Github.Metadata.Repository
		asset    string                  = conf.Github.Metadata.Asset
		dataFile string                  = "accounts.aws.json"
		sw                               = utils.Stopwatch()
	)
	defer service.Close()
	// timer
	sw.Start()

	log = log.With("operation", "Existing", "service", "awsaccount")
	log.Info("[awsaccount] starting existing records import ...")

	data, err = service.DownloadAndReturn(owner, repo, asset, dataFile)

	for _, row := range data {
		inserts = append(inserts, &sqldb.BoundStatement{Statement: stmtImport, Data: row})
	}
	log.With("count", len(inserts)).Debug("records to insert ...")

	log.Debug("creating writer store for insert...")
	store, err = sqldb.New[*AwsAccountImport](ctx, log, conf)
	if err != nil {
		return
	}

	log.Debug("running insert ...")
	err = store.Insert(inserts...)
	if err != nil {
		return
	}

	log.Debug("set empty environment values ...")
	_, err = store.Exec(stmtUpdateEmptyEnvironments)
	if err != nil {
		return
	}

	log.With(
		"seconds", sw.Stop().Seconds(),
		"inserted", len(inserts)).
		Info("[awsaccount] existing records imported.")

	return
}
