package awsaccounts

import (
	"context"
	"log/slog"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/awsaccount"
	"github.com/ministryofjustice/opg-reports/report/internal/gh"
	"github.com/ministryofjustice/opg-reports/report/internal/opgmetadata"
	"github.com/ministryofjustice/opg-reports/report/internal/sqldb"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Cmd returns the cobra command and handles binding of cli args into the
// config setup
//
// Example:
//
//	importer awsaccounts --db=$path --org=githubOrgansiation
//
// interface: ImporterCLICommand
func Cmd(conf *config.Config, viperConf *viper.Viper) (cmd *cobra.Command) {
	// handles importing only team data
	cmd = &cobra.Command{
		Use:   "awsaccounts",
		Short: "AwsAccounts command handles importing 'awsaccounts' list.",
		Long:  `awsaccounts command imports accounts from the opg-metadata repository. In order to use this command you will need a GITHUB_TOKEN set with correct permissions.`,
		Example: `
  importer awsaccounts
	--org=$githubOrganisationName`,
		Run: func(cmd *cobra.Command, args []string) {
			var (
				ctx context.Context = context.Background()
				log *slog.Logger    = utils.Logger(conf.Log.Level, conf.Log.Type)
			)
			// import teams first
			Import(ctx, log, conf)
		},
	}
	// bind the github.organisation config item to the shorter --org
	cmd.Flags().StringVar(&conf.Github.Organisation, "org", conf.Github.Organisation, "GitHub organisation name")
	viperConf.BindPFlag("github.organisation", cmd.Flags().Lookup("org"))
	return
}

// Import generates new aws account from the accounts.json information within the
// opg-metadata published data. The accounts.aws.json is parsed and converted to db entries
//
// That is a private repository so you need permissions to read and fetch data to be
// able to download the release asset.
//
// interface: ImporterImportCommand
func Import(ctx context.Context, log *slog.Logger, conf *config.Config) (err error) {
	log.Info("running [team] imports ...")

	// to import teams, we create an opgmetadata service and call the getTeams
	// so fetch the gh repository first and then create the opgmeta data service
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
	list := []*awsaccount.AwsAccount{}
	err = utils.Convert(rawAccounts, &list)
	if err != nil {
		return
	}
	// before we insert, set environment to production if empty
	for _, acc := range list {
		if acc.Environment == "" {
			acc.Environment = "production"
		}
	}
	// sqldb
	store, err := sqldb.New[*awsaccount.AwsAccount](ctx, log, conf)
	if err != nil {
		return
	}
	// service
	srv, err := awsaccount.NewService(ctx, log, conf, store)
	if err != nil {
		return
	}
	_, err = srv.Import(list)

	return
}
