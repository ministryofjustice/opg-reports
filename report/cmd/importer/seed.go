package main

import (
	"github.com/ministryofjustice/opg-reports/report/internal/repository/sqlr"
	"github.com/ministryofjustice/opg-reports/report/internal/service/seed"
	"github.com/spf13/cobra"
)

// seedCmd uses fixture / seed data to populate a fresh database which can then
// be used for local dev / testing
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "seed inserts known test data.",
	Long: `
seed inserts known test data for use in development.

env variables used that can be adjusted:

	DATABASE_PATH
		The file path to the sqlite database that will be used
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var (
			sqlStore    *sqlr.Repository = sqlr.Default(ctx, log, conf)
			seedService *seed.Service    = seed.Default(ctx, log, conf)
		)
		_, err = seedService.All(sqlStore)
		return
	},
}
