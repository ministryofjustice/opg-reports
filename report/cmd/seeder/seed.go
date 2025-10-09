package main

import (
	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/service/seed"

	"github.com/spf13/cobra"
)

// allCmd uses fixture / seed data to populate a fresh database which can then
// be used for local dev / testing
var allCmd = &cobra.Command{
	Use:   "all",
	Short: "all inserts all seeds in order",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var (
			sqlStore    *sqlr.Repository = sqlr.Default(ctx, log, conf)
			seedService *seed.Service    = seed.Default(ctx, log, conf)
		)
		_, err = seedService.All(sqlStore)
		return
	},
}
