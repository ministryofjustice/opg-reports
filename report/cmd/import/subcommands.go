package main

import (
	"opg-reports/report/internal/account/accountimport"
	"opg-reports/report/internal/codebases/codebasesimport"
	"opg-reports/report/internal/codebasestats/codebasestatsimport"
	"opg-reports/report/internal/codeowners/codeownersimport"
	"opg-reports/report/internal/cost/costimport"
	"opg-reports/report/internal/global/migrations"
	"opg-reports/report/internal/team/teamimport"
	"opg-reports/report/internal/uptime/uptimeimport"
	"opg-reports/report/package/awsclients"
	"opg-reports/report/package/awsid"
	"opg-reports/report/package/env"
	"opg-reports/report/package/ghclients"
	"opg-reports/report/package/times"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/google/go-github/v84/github"
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

// uptime import command
var uptimeCmd = &cobra.Command{
	Use:   `uptime`,
	Short: `import uptime`,
	RunE:  runUptimeImport,
}

// codebase import command
var codebasesCmd = &cobra.Command{
	Use:   `codebases`,
	Short: `import codebases`,
	RunE:  runCodebaseImport,
}

// codeowner import command
var codeownersCmd = &cobra.Command{
	Use:   `codeowners`,
	Short: `import codeowners`,
	RunE:  runCodeownersImport,
}

// codebase stats import command
var codebaseStatsCmd = &cobra.Command{
	Use:   `codebase-stats`,
	Short: `import codebase stats`,
	RunE:  runCodebaseStatsImport,
}

// runTeamsImport runs the teams
func runTeamsImport(cmd *cobra.Command, args []string) (err error) {
	var ctx = cmd.Context()
	// overwrite arg flags from env values
	if e := env.OverwriteStruct(&flags); e != nil {
		return
	}
	// run the migrations
	err = migrations.Migrate(ctx, &migrations.Args{
		DB:     flags.DB,
		Driver: flags.Driver,
		Params: flags.Params,
	})
	if err != nil {
		return
	}
	// run the import
	err = teamimport.Import(ctx, &teamimport.Args{
		DB:      flags.DB,
		Driver:  flags.Driver,
		Params:  flags.Params,
		SrcFile: flags.SrcFile,
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
	err = migrations.Migrate(ctx, &migrations.Args{
		DB:     flags.DB,
		Driver: flags.Driver,
		Params: flags.Params,
	})
	if err != nil {
		return
	}
	// run the import
	err = accountimport.Import(ctx, &accountimport.Args{
		DB:      flags.DB,
		Driver:  flags.Driver,
		Params:  flags.Params,
		SrcFile: flags.SrcFile,
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
	err = migrations.Migrate(ctx, &migrations.Args{
		DB:     flags.DB,
		Driver: flags.Driver,
		Params: flags.Params,
	})
	if err != nil {
		return
	}

	err = costimport.Import(ctx, client, &costimport.Args{
		DB:        flags.DB,
		Driver:    flags.Driver,
		Params:    flags.Params,
		DateStart: times.MustFromString(flags.DateStartCosts), // use the other start date thats further back in time
		DateEnd:   times.MustFromString(flags.DateEnd),
		AccountID: awsid.AccountID(ctx, flags.Region),
	})
	return
}

// runUptimeImport runs the uptime import
func runUptimeImport(cmd *cobra.Command, args []string) (err error) {
	var client *cloudwatch.Client
	var region = "us-east-1" // forced region
	var ctx = cmd.Context()
	// overwrite arg flags from env values
	if e := env.OverwriteStruct(&flags); e != nil {
		return
	}
	client, err = awsclients.New[*cloudwatch.Client](ctx, region)
	if err != nil {
		return
	}
	// run the migrations
	err = migrations.Migrate(ctx, &migrations.Args{
		DB:     flags.DB,
		Driver: flags.Driver,
		Params: flags.Params,
	})
	if err != nil {
		return
	}

	err = uptimeimport.Import(ctx, client, &uptimeimport.Args{
		DB:        flags.DB,
		Driver:    flags.Driver,
		Params:    flags.Params,
		DateStart: times.MustFromString(flags.DateStart),
		DateEnd:   times.MustFromString(flags.DateEnd),
		AccountID: awsid.AccountID(ctx, flags.Region),
	})
	return
}

// runCodebaseImport runs the codebase import with stats
func runCodebaseImport(cmd *cobra.Command, arglist []string) (err error) {
	var client *github.Client
	var tk = os.Getenv("GITHUB_TOKEN")
	var ctx = cmd.Context()
	// overwrite arg flags from env values
	if e := env.OverwriteStruct(&flags); e != nil {
		return
	}
	client, err = ghclients.New(ctx, tk)
	if err != nil {
		return
	}
	// run the migrations
	err = migrations.Migrate(ctx, &migrations.Args{
		DB:     flags.DB,
		Driver: flags.Driver,
		Params: flags.Params,
	})
	if err != nil {
		return
	}

	err = codebasesimport.Import(ctx, client.Teams, &codebasesimport.Args{
		DB:           flags.DB,
		Driver:       flags.Driver,
		Params:       flags.Params,
		OrgSlug:      flags.OrgSlug,
		ParentSlug:   flags.ParentSlug,
		FilterByName: flags.Filter,
	})
	return
}

// runCodeownersImport runs the codebase and codeowner import - this is a bit slower due to fetching files
func runCodeownersImport(cmd *cobra.Command, arglist []string) (err error) {
	var client *github.Client
	var tk = os.Getenv("GITHUB_TOKEN")
	var ctx = cmd.Context()
	// overwrite arg flags from env values
	if e := env.OverwriteStruct(&flags); e != nil {
		return
	}
	client, err = ghclients.New(ctx, tk)
	if err != nil {
		return
	}
	// run the migrations
	err = migrations.Migrate(ctx, &migrations.Args{
		DB:     flags.DB,
		Driver: flags.Driver,
		Params: flags.Params,
	})
	if err != nil {
		return
	}

	clients := &codeownersimport.Clients{
		Teams: client.Teams,
		Repos: client.Repositories,
	}

	err = codeownersimport.Import(ctx, clients, &codeownersimport.Args{
		DB:           flags.DB,
		Driver:       flags.Driver,
		Params:       flags.Params,
		OrgSlug:      flags.OrgSlug,
		ParentSlug:   flags.ParentSlug,
		FilterByName: flags.Filter,
	})
	return
}

// runCodebaseStatsImport runs code base import with stats data
func runCodebaseStatsImport(cmd *cobra.Command, arglist []string) (err error) {
	var client *github.Client
	var tk = os.Getenv("GITHUB_TOKEN")
	var ctx = cmd.Context()
	// overwrite arg flags from env values
	if e := env.OverwriteStruct(&flags); e != nil {
		return
	}
	client, err = ghclients.New(ctx, tk)
	if err != nil {
		return
	}
	// run the migrations
	err = migrations.Migrate(ctx, &migrations.Args{
		DB:     flags.DB,
		Driver: flags.Driver,
		Params: flags.Params,
	})
	if err != nil {
		return
	}

	clients := &codebasestatsimport.Clients{
		Teams: client.Teams,
		Repos: client.Repositories,
	}

	err = codebasestatsimport.Import(ctx, clients, &codebasestatsimport.Args{
		DB:           flags.DB,
		Driver:       flags.Driver,
		Params:       flags.Params,
		OrgSlug:      flags.OrgSlug,
		ParentSlug:   flags.ParentSlug,
		FilterByName: flags.Filter,
	})
	return
}
