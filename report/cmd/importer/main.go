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

	"github.com/ministryofjustice/opg-reports/report/cmd/importer/teams"
	"github.com/ministryofjustice/opg-reports/report/config"
)

const databasePath string = "./__database/api.db"

// importFunc is a short hand to express the pattern of functions used so the importers var
// is cleaner
type importFunc func(ctx context.Context, log *slog.Logger, conf *config.Config, args []string) (err error)

var importers map[string][]importFunc = map[string][]importFunc{
	"all":   {teams.Import},
	"teams": {teams.Import},
}

// runner checks the choice against the known list of import options
// and will run each function in order thats attached
func runner(ctx context.Context, log *slog.Logger, choice string, args []string) (err error) {

	var conf = config.NewConfig()
	// setup some config overrides
	conf.Database.Path = databasePath

	// look for an run functions to import
	if funcList, ok := importers[choice]; ok {
		for _, lambdaF := range funcList {
			// run the function and if there is an, flag that
			// to bre returned
			e := lambdaF(ctx, log, conf, args)
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

	err = runner(ctx, log, os.Args[1], os.Args[2:])

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}
