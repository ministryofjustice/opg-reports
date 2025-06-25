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
	log.Info("running [awsaccounts] imports ...")

	var now = time.Now().UTC().Format(time.RFC3339)
	// fetch the gh repository first and then create the opgmeta data service
	gh, err := gh.New(ctx, log, conf)
	if err != nil {
		return
	}
	metaService, err := opgmetadata.NewService(ctx, log, conf, gh)
	if err != nil {
		return
	}

	// get just the aws accounts
	rawAccounts, err := metaService.GetAllAwsAccounts()
	if err != nil {
		return
	}

	// convert to db model
	list := []*AwsAccountImport{}
	err = utils.Convert(rawAccounts, &list)
	if err != nil {
		return
	}
	// before we insert, set environment to production if empty
	for _, acc := range list {
		if acc.Environment == "" {
			acc.Environment = "production"
		}
		if acc.CreatedAt == "" {
			acc.CreatedAt = now
		}
	}
	// sqldb
	store, err := sqldb.New[*AwsAccount](ctx, log, conf)
	if err != nil {
		return
	}
	// service
	srv, err := NewService(ctx, log, conf, store)
	if err != nil {
		return
	}

	_, err = srv.Import(list)

	return
}
