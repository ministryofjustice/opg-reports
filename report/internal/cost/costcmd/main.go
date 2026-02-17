package main

import (
	"context"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/logger"
	"os"

	"github.com/spf13/cobra"
)

var root *cobra.Command = &cobra.Command{
	Use:               "costcmd",
	Short:             `costcmd local wrapper to test capabilities`,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

func main() {
	var err error
	var log = logger.New(os.Getenv("LOG_LEVEL"), os.Getenv("LOG_TYPE"))
	var ctx = cntxt.AddLogger(context.Background(), log)

	root.AddCommand(
		migrationCmd,
		importCmd,
	)

	err = root.ExecuteContext(ctx)
	if err != nil {
		log.Error("error running command", "err", err.Error())
		os.Exit(1)
	}
}
