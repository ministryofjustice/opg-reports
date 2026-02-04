/*
db used to download / upload database to configured locations.

Used by the reporting workflows to fetch the database from s3.

# Configure the location of datbase via environment varaibles

Usage:

	db [commands]

Available commands:

	download
	upload

# Example

`aws-vault exec <profile> -- env DATABASE_PATH="./data/api.db" db download`
*/
package main

import (
	"context"
	"log/slog"
	"opg-reports/report/conf"
	"opg-reports/report/internal/utils/logger"
	"os"

	"github.com/spf13/cobra"
)

const (
	cmdName   string = "db" // root command name
	shortDesc string = `db downloads or uploads sqlite database from s3`
	longDesc  string = `
db downloads or uploads the sqlite database from s3 bucket configured via argments with defaults.
`
)

// config items
var (
	cfg *conf.Config    // default config
	ctx context.Context // default context
	log *slog.Logger    // default logger
)

var (
	rootCmd *cobra.Command = &cobra.Command{
		Use:   cmdName,
		Short: shortDesc,
		Long:  longDesc,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}
)

func setup() {
	cfg = conf.New()
	ctx = context.Background()
	log = logger.New(cfg.Log.Level, cfg.Log.Type)
}

// setup default values for config and logging & add options
func init() {
	setup()
}

func main() {
	var err error

	err = rootCmd.Execute()
	if err != nil {
		log.Error("error running command", "err", err.Error())
		os.Exit(1)
	}
}
