package awsaccount

import (
	"context"
	"log/slog"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/interfaces"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/sqlr"
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
func Existing[T interfaces.Model](ctx context.Context, log *slog.Logger, conf *config.Config, service interfaces.MetadataService[T]) (err error) {
	var (
		data     []T
		store    *sqlr.RepositoryWithSelect[*AwsAccountImport]
		inserts  []*sqlr.BoundStatement = []*sqlr.BoundStatement{}
		owner    string                 = conf.Github.Organisation
		repo     string                 = conf.Github.Metadata.Repository
		asset    string                 = conf.Github.Metadata.Asset
		dataFile string                 = "accounts.aws.json"
		sw                              = utils.Stopwatch()
	)
	defer func() {
		service.Close()
		log.With("seconds", sw.Stop().Seconds(), "inserted", len(inserts)).Info("[awsaccount] existing records end.")
	}()

	// timer
	sw.Start()

	log = log.With("operation", "Existing", "service", "awsaccount")
	log.Info("[awsaccount] starting existing records import ...")

	data, err = service.DownloadAndReturn(owner, repo, asset, false, dataFile)

	for _, row := range data {
		inserts = append(inserts, &sqlr.BoundStatement{Statement: stmtImport, Data: row})
	}
	log.With("count", len(inserts)).Debug("[awsaccount] records to insert ...")

	log.Debug("[awsaccount] creating writer store for insert...")
	store, err = sqlr.NewWithSelect[*AwsAccountImport](ctx, log, conf)
	if err != nil {
		return
	}

	log.Debug("[awsaccount] running insert ...")
	err = store.Insert(inserts...)
	if err != nil {
		return
	}

	log.Debug("[awsaccount] set empty environment values ...")
	_, err = store.Exec(stmtUpdateEmptyEnvironments)
	if err != nil {
		return
	}
	log.Info("[awsaccount] existing records successful")

	return
}
