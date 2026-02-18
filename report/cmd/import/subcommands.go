package main

import (
	"opg-reports/report/internal/account/accountimport"
	"opg-reports/report/internal/cost/costimport"
	"opg-reports/report/internal/global/migrations"
	"opg-reports/report/internal/team/teamimport"
	"opg-reports/report/package/awsclients"
	"opg-reports/report/package/awsid"
	"opg-reports/report/package/env"
	"opg-reports/report/package/times"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/spf13/cobra"
)

// team import command
var teamsCmd = &cobra.Command{
	Use:   `teams`,
	Short: `import teams`,
	RunE:  runTeamsImport,
}

// accounts import command
var accountsCmd = &cobra.Command{
	Use:   `accounts`,
	Short: `import accounts`,
	RunE:  runAccountsImport,
}

// costs import command
var costsCmd = &cobra.Command{
	Use:   `costs`,
	Short: `import costs`,
	RunE:  runCostsImport,
}

// runTeamsImport runs the teams
func runTeamsImport(cmd *cobra.Command, args []string) (err error) {
	var ctx = cmd.Context()
	// overwrite arg flags from env values
	if e := env.OverwriteStruct(&flags); e != nil {
		return
	}
	// run the migrations
	err = migrations.MigrateAll(ctx, &migrations.Args{
		DB:            flags.DB,
		Driver:        flags.Driver,
		Params:        flags.Params,
		MigrationFile: flags.MigrationFile,
	})
	if err != nil {
		return
	}
	// run the import
	err = teamimport.Import(ctx, &teamimport.Args{
		DB:            flags.DB,
		Driver:        flags.Driver,
		Params:        flags.Params,
		MigrationFile: flags.MigrationFile,
		SrcFile:       flags.SrcFile,
	})
	return
}

// runAccountsImport runs the accounts import
func runAccountsImport(cmd *cobra.Command, args []string) (err error) {
	var ctx = cmd.Context()
	// overwrite arg flags from env values
	if e := env.OverwriteStruct(&flags); e != nil {
		return
	}
	// run the migrations
	err = migrations.MigrateAll(ctx, &migrations.Args{
		DB:            flags.DB,
		Driver:        flags.Driver,
		Params:        flags.Params,
		MigrationFile: flags.MigrationFile,
	})
	if err != nil {
		return
	}
	// run the import
	err = accountimport.Import(ctx, &accountimport.Args{
		DB:            flags.DB,
		Driver:        flags.Driver,
		Params:        flags.Params,
		MigrationFile: flags.MigrationFile,
		SrcFile:       flags.SrcFile,
	})
	return
}

// runCostsImport runs the costs
func runCostsImport(cmd *cobra.Command, args []string) (err error) {
	var client *costexplorer.Client
	var ctx = cmd.Context()
	// overwrite arg flags from env values
	if e := env.OverwriteStruct(&flags); e != nil {
		return
	}
	client, err = awsclients.New[*costexplorer.Client](ctx, flags.Region)
	if err != nil {
		return
	}
	// run the migrations
	err = migrations.MigrateAll(ctx, &migrations.Args{
		DB:            flags.DB,
		Driver:        flags.Driver,
		Params:        flags.Params,
		MigrationFile: flags.MigrationFile,
	})
	if err != nil {
		return
	}

	err = costimport.Import(ctx, client, &costimport.Args{
		DB:            flags.DB,
		Driver:        flags.Driver,
		Params:        flags.Params,
		MigrationFile: flags.MigrationFile,
		DateStart:     times.MustFromString(flags.DateStart),
		DateEnd:       times.MustFromString(flags.DateEnd),
		AccountID:     awsid.AccountID(ctx, flags.Region),
	})
	return
}
