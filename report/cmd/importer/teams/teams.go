package teams

import (
	"context"
	"flag"
	"log/slog"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/gh"
	"github.com/ministryofjustice/opg-reports/report/internal/opgmetadata"
	"github.com/ministryofjustice/opg-reports/report/internal/sqldb"
	"github.com/ministryofjustice/opg-reports/report/internal/team"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

const githubOrg string = "ministryofjustice"

// options this import function will accept from the cli
type cliArgs struct {
	Organisation string // Overwrites the github organisation being called
}

var (
	options = &cliArgs{}
	flagSet = flag.NewFlagSet("teams", flag.ExitOnError)
)

// Import generates new team data from the billing_unit information within the
// opg-metadata published list of accounts.
//
// That is a private repository so you need permissions to read and fetch data to be
// able to download the release asset.
//
// The account.json is parsed and all unique billing_units are converted into team.Team
// entries and inserted into the databse using the team.Service.Import method
func Import(ctx context.Context, log *slog.Logger, conf *config.Config, args []string) (err error) {
	log.Info("running [team] imports ...")

	// arg handling
	flagSet.StringVar(&options.Organisation, "org", githubOrg, "name of the github org.")
	flagSet.Parse(args)
	if options.Organisation == "" {
		options.Organisation = githubOrg
	}

	// setup the config
	conf.Github.Organisation = options.Organisation

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
