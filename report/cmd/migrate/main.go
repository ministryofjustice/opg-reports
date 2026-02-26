package main

import (
	"context"
	"opg-reports/report/internal/global/migrations"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/logger"
	"os"

	"github.com/spf13/cobra"
)

var flags = &migrations.Args{
	Driver: "sqlite3",
	DB:     "./database/api.db",
}

var runConversion bool = false

var root *cobra.Command = &cobra.Command{
	Use:   "migrate",
	Short: `run migrations for the database`,
	RunE:  runCMD,
}

func runCMD(cmd *cobra.Command, args []string) (err error) {
	var ctx = cmd.Context()
	err = migrations.Migrate(ctx, flags)
	if err != nil {
		return
	}
	if runConversion {
		migrations.Convert(ctx, flags)
	}
	return
}

func init() {
	root.PersistentFlags().StringVar(&flags.Driver, "driver", flags.Driver, "Database driver")
	root.PersistentFlags().StringVar(&flags.DB, "db", flags.DB, "Database path")
	root.PersistentFlags().StringVar(&flags.Params, "params", flags.Params, "Database params")
	root.PersistentFlags().BoolVar(&runConversion, "convert", runConversion, "Run DB conversion to upgrade from older structure")

}

func main() {
	var err error
	var log = logger.New(os.Getenv("LOG_LEVEL"), os.Getenv("LOG_TYPE"))
	var ctx = cntxt.AddLogger(context.Background(), log)

	err = root.ExecuteContext(ctx)
	if err != nil {
		log.Error("error running command", "err", err.Error())
		os.Exit(1)
	}
}
