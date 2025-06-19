/*
Import data into a database

	importer
*/
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/gh"
	"github.com/ministryofjustice/opg-reports/report/internal/opgmetadata"
	"github.com/ministryofjustice/opg-reports/report/internal/sqldb"
	"github.com/ministryofjustice/opg-reports/report/internal/team"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

const databasePath string = "./__database/api.db"

// importFunc is a short hand to express the pattern of functions used so the importers var
// is cleaner
type importFunc func(ctx context.Context, log *slog.Logger, conf *config.Config) (err error)

var importers map[string][]importFunc = map[string][]importFunc{
	"all":   {ImportTeams},
	"teams": {ImportTeams},
}

// ImportTeams generates new team data from the billing_unit information within the
// opg-metadata published list of accounts.
//
// That is a private repository so you need permissions to read and fetch data to be
// able to download the release asset.
//
// The account.json os parsed and all unique billing_units are converted into team.Team
// entries and inserted into the databse using the team.Service.Import method
func ImportTeams(ctx context.Context, log *slog.Logger, conf *config.Config) (err error) {
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

// runner checks the choice against the known list of import options
// and will run each function in order thats attached
func runner(ctx context.Context, log *slog.Logger, choice string) (err error) {

	var conf = config.NewConfig()
	// setup some config overrides
	conf.Github.Organisation = "ministryofjustice"
	conf.Database.Path = databasePath

	if funcList, ok := importers[choice]; ok {
		for _, lambdaF := range funcList {
			// run the function and if there is an, flag that
			// to bre returned
			e := lambdaF(ctx, log, conf)
			if e != nil {
				err = errors.Join(e, err)
				return
			}
		}
	} else {
		err = fmt.Errorf("unsupported import [%s]", choice)
	}

	return
}

func main() {
	var (
		err error
		log *slog.Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
		ctx              = context.Background()
	)

	if len(os.Args) <= 1 {
		log.Error("no arguments passed")
		os.Exit(1)
	}

	err = runner(ctx, log, os.Args[1])

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}
