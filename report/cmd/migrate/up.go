package main

import (
	"opg-reports/report/internal/repository/sqlr"

	"github.com/spf13/cobra"
)

// upCmd runs the schema sql agains the DB
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "up runs exec against the DB of the existing SCHEMA COMMAND",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var dbR *sqlr.Repository
		dbR = sqlr.Default(ctx, log, conf)

		err = dbR.Ping()
		if err != nil {
			return
		}

		err = sqlr.MigrateUp(dbR)

		return
	},
}
