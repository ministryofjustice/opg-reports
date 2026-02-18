package main

import (
	"context"
	"opg-reports/report/internal/cost/costmigrate"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/logger"
	"os"

	"github.com/spf13/cobra"
)

type cli struct {
	DB            string `json:"db"`             // --db
	Driver        string `json:"driver"`         // --driver
	Params        string `json:"params"`         // --params
	MigrationFile string `json:"migration_file"` // --file
}

var flags = &cli{
	Driver:        "sqlite3",
	DB:            "./database/api.db",
	MigrationFile: "migrations.json",
}

var root *cobra.Command = &cobra.Command{
	Use:   "migrate",
	Short: `run migrations for the database`,
	RunE:  runCMD,
}

func runCMD(cmd *cobra.Command, args []string) (err error) {
	// run the migration command for costs
	err = costmigrate.Migrate(cmd.Context(), &costmigrate.Input{
		DB:            flags.DB,
		Driver:        flags.Driver,
		Params:        flags.Params,
		MigrationFile: flags.MigrationFile,
	})
	return
}

func init() {
	root.PersistentFlags().StringVar(&flags.Driver, "driver", flags.Driver, "Database driver")
	root.PersistentFlags().StringVar(&flags.DB, "db", flags.DB, "Database path")
	root.PersistentFlags().StringVar(&flags.Params, "params", flags.Params, "Database params")
	root.PersistentFlags().StringVar(&flags.MigrationFile, "file", flags.MigrationFile, "migration file")
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
