package main

import (
	"context"
	"opg-reports/report/internal/cost/costmigrate"
)

func migrate(ctx context.Context, flags *cli) (err error) {
	// run the migration command
	err = costmigrate.Migrate(ctx, &costmigrate.Input{
		DB:            flags.DB,
		Driver:        flags.Driver,
		Params:        flags.Params,
		MigrationFile: flags.MigrationFile,
	})
	return
}
