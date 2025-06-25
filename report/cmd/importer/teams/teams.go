package teams

import (
	"context"
	"log/slog"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/gh"
	"github.com/ministryofjustice/opg-reports/report/internal/opgmetadata"
	"github.com/ministryofjustice/opg-reports/report/internal/sqldb"
	"github.com/ministryofjustice/opg-reports/report/internal/team"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Cmd returns the cobra command and handles binding of cli args into the
// config setup
//
// Example:
//
//	importer teams --db=$path --org=githubOrgansiation
func Cmd(conf *config.Config, viperConf *viper.Viper) (cmd *cobra.Command) {
	// handles importing only team data
	cmd = &cobra.Command{
		Use:   "teams",
		Short: "Teams command handles importing 'teams' grouping data.",
		Long:  `teams command imports organsiational teams from the opg-metadata repository. In order to use this command you will need a GITHUB_TOKEN set with correct permissions.`,
		Example: `
  importer teams
	--db=$pathToFile
	--gh-org=$githubOrganisationName`,
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

// Import generates new team data from the billing_unit information within the
// opg-metadata published list of accounts.
//
// That is a private repository so you need permissions to read and fetch data to be
// able to download the release asset.
//
// The account.json is parsed and all unique billing_units are converted into team.Team
// entries and inserted into the databse using the team.Service.Import method
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

	rawTeams, err := metaService.GetAllTeams()
	if err != nil {
		return
	}
	// now we have raw team data, we need to create a team store and service
	// convert the maps into structs and import to the database

	// convert raw to teams
	list := []*team.Team{}
	err = utils.Convert(rawTeams, &list)
	// sqldb
	store, err := sqldb.New[*team.Team](ctx, log, conf)
	if err != nil {
		return
	}
	// service
	teamService, err := team.NewService(ctx, log, conf, store)
	if err != nil {
		return
	}
	_, err = teamService.Import(list)

	return
}
