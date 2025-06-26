package awsaccount

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

// Import generates new aws account data from the accounts.json information within the
// opg-metadata published data. The accounts.aws.json is parsed and converted to db entries
//
// That is a private repository so you need permissions to read and fetch data to be
// able to download the release asset.
//
// interface: ImporterImportCommand
func Import(ctx context.Context, log *slog.Logger, conf *config.Config) (err error) {
	var (
		gr          *gh.Repository
		metaService *opgmetadata.Service
		rawAccounts []map[string]interface{}
		list        []*awsAccountImport
		store       *sqldb.Repository[*AwsAccount]
		now         string                  = time.Now().UTC().Format(time.RFC3339)
		inserts     []*sqldb.BoundStatement = []*sqldb.BoundStatement{}
	)
	log = log.With("operation", "Import", "service", "awsaccount")
	log.Info("running [awsaccounts] imports ...")

	// fetch the gh repository first and then create the opgmeta data service
	gr, err = gh.New(ctx, log, conf)
	if err != nil {
		return
	}
	// get the service
	metaService, err = opgmetadata.NewService(ctx, log, conf, gr)
	if err != nil {
		return
	}
	defer metaService.Close()

	// get just the aws accounts
	rawAccounts, err = metaService.GetAllAwsAccounts()
	if err != nil {
		return
	}
	// convert to db model
	list = []*awsAccountImport{}
	err = utils.Convert(rawAccounts, &list)
	if err != nil {
		return
	}
	// sqldb setup
	store, err = sqldb.New[*AwsAccount](ctx, log, conf)
	if err != nil {
		return
	}
	log.Info("importing ...")
	log.Debug("generating bound statements ...")
	for _, row := range list {
		if row.Environment == "" {
			row.Environment = "production"
		}
		row.CreatedAt = now
		inserts = append(inserts, &sqldb.BoundStatement{Statement: stmtImport, Data: row})
	}
	log.Debug("running insert ...")
	err = store.Insert(inserts...)

	return
}
