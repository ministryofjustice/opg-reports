package main

import (
	"opg-reports/report/internal/global/imports"
	"opg-reports/report/package/env"

	"github.com/spf13/cobra"
)

// costs import command
var costsCmd = &cobra.Command{
	Use:   `costs`,
	Short: `import costs`,
	RunE:  runCostsImport,
}

// runCostsImport runs the costs - matches the costcmd version
func runCostsImport(cmd *cobra.Command, args []string) (err error) {
	var ctx = cmd.Context()
	// overwrite arg flags from env values
	if e := env.OverwriteStruct(&flags); e != nil {
		return
	}
	err = imports.ImportCosts(ctx, flags)
	return
}
