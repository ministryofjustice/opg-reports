package seeds

import (
	"context"
	"fmt"
	"math/rand/v2"
	"opg-reports/report/internal/account/accountimport"
	"opg-reports/report/internal/codebases/codebasesimport"
	"opg-reports/report/internal/cost/costimport"
	"opg-reports/report/internal/global/migrations"
	"opg-reports/report/internal/team/teamimport"
	"opg-reports/report/internal/uptime/uptimeimport"
	"opg-reports/report/package/dbx"
	"opg-reports/report/package/times"
)

var teamList []string = []string{
	"team-a",
	"team-b",
	"team-c",
	"team-d",
	"team-e",
	"team-f",
}

// environments are shared account environment names to use in seeded data
var environmentList []string = []string{
	"development",
	"pre-production",
	"integrations",
	"production",
}

// fix regions to used for seeds
var regionList []string = []string{
	"eu-west-1",
	"eu-west-2",
	"us-east-1",
	"NoRegion",
}

// aws service that we use for seeding
var serviceList []string = []string{
	"Amazon Relational Database Service",
	"Amazon Simple Storage Service",
	"AmazonCloudWatch",
	"Amazon Elastic Load Balancing",
	"AWS Shield",
	"AWS Config",
	"AWS CloudTrail",
	"AWS Key Management Service",
	"Amazon Virtual Private Cloud",
	"Amazon Elastic Container Service",
	"EC2 - Other",
}

var codebaseCompliance []string = []string{
	"unknown",
	"baseline",
	"standard",
	"exemplar",
}

// Results contains all the seed data that was inserted
// including any that may have failed
type Results struct {
	Teams     []*teamimport.Model      `json:"teams"`
	Accounts  []*accountimport.Model   `json:"accounts"`
	Costs     []*costimport.Model      `json:"costs"`
	Uptime    []*uptimeimport.Model    `json:"uptime"`
	Codebases []*codebasesimport.Model `json:"codebases"`
}

// Args
type Args struct {
	DB     string `json:"db"`     // --db
	Driver string `json:"driver"` // --driver
	Params string `json:"params"` // --params
}

// SeedAll
func SeedAll(ctx context.Context, in *Args) (results *Results, err error) {
	var (
		numAccounts  = 25
		numCosts     = 13000
		numUptime    = 1200
		numCodebases = 50
	)

	var args = &dbx.InsertArgs{
		DB:     in.DB,
		Driver: in.Driver,
		Params: in.Params,
	}
	results = &Results{}
	// run the migrations
	err = migrations.Migrate(ctx, &migrations.Args{
		DB:     in.DB,
		Driver: in.Driver,
		Params: in.Params,
	})
	if err != nil {
		return
	}
	// seed teams
	results.Teams, err = seedTeams(ctx, args)
	if err != nil {
		return
	}
	// seed accounts
	results.Accounts, err = seedAccounts(ctx, args, numAccounts, results.Teams)
	if err != nil {
		return
	}
	// seed costs
	results.Costs, err = seedCosts(ctx, args, numCosts, results.Accounts)
	if err != nil {
		return
	}
	// seed uptime
	results.Uptime, err = seedUptime(ctx, args, numUptime, results.Accounts)
	if err != nil {
		return
	}
	// seed codebases
	results.Codebases, err = seedCodebases(ctx, args, numCodebases)
	if err != nil {
		return
	}

	return
}

func seedCodebases(ctx context.Context, in *dbx.InsertArgs, n int) (insert []*codebasesimport.Model, err error) {
	var githubOrg = "mock-org"
	insert = []*codebasesimport.Model{}
	for i := 0; i < n; i++ {
		var name = fmt.Sprintf("codebase-%02d", i+1)
		var compI = rand.IntN(len(codebaseCompliance))
		insert = append(insert, &codebasesimport.Model{
			Name:                name,
			FullName:            fmt.Sprintf("%s/%s", githubOrg, name),
			Url:                 fmt.Sprintf("https://mock-github.local/%s/%s", githubOrg, name),
			ComplianceReportUrl: fmt.Sprintf("https://mock-compliance-report.local/%s", name),
			ComplianceBadge:     fmt.Sprintf("https://mock-compliance-report.local/%s/badge", name),
			ComplianceLevel:     codebaseCompliance[compI],
		})
	}
	err = dbx.Insert(ctx, codebasesimport.InsertStatement, insert, in)
	return
}

// seedUptime generates and inserts uptime data
func seedUptime(ctx context.Context, in *dbx.InsertArgs, n int, accounts []*accountimport.Model) (insert []*uptimeimport.Model, err error) {
	var (
		end    = times.ResetMonth(times.Today())
		start  = times.ResetMonth(times.Add(end, -3, times.YEAR))
		months = times.Months(start, end)
	)
	insert = []*uptimeimport.Model{}

	for i := 0; i < n; i++ {
		var accountI = rand.IntN(len(accounts))
		var monthI = rand.IntN(len(months))
		var avg float64 = (95) + (rand.Float64() * (100 - 95)) // 95-100%

		insert = append(insert, &uptimeimport.Model{
			Month:       times.AsYMString(months[monthI]),
			AccountID:   accounts[accountI].ID,
			Granularity: "3600",
			Average:     fmt.Sprintf("%g", avg),
		})
	}
	err = dbx.Insert(ctx, uptimeimport.InsertStatement, insert, in)

	return
}

// seedCosts generates and inserts cost data similar to real life values
func seedCosts(ctx context.Context, in *dbx.InsertArgs, n int, accounts []*accountimport.Model) (insert []*costimport.Model, err error) {
	var (
		end    = times.ResetMonth(times.Today())
		start  = times.ResetMonth(times.Add(end, -3, times.YEAR))
		months = times.Months(start, end)
	)
	insert = []*costimport.Model{}

	for i := 0; i < n; i++ {
		var accountI = rand.IntN(len(accounts))
		var monthI = rand.IntN(len(months))
		var regionI = rand.IntN(len(regionList))
		var serviceI = rand.IntN(len(serviceList))
		var price float64 = (-1000.0) + (rand.Float64() * (1000 - -1000.0)) // 95-100%

		insert = append(insert, &costimport.Model{
			Region:    regionList[regionI],
			Service:   serviceList[serviceI],
			Month:     times.AsYMString(months[monthI]),
			Cost:      fmt.Sprintf("%g", price),
			AccountID: accounts[accountI].ID,
		})
	}
	err = dbx.Insert(ctx, costimport.InsertStatement, insert, in)

	return
}

// seedAccounts generates and inserts cost data similar to real life values
func seedAccounts(ctx context.Context, in *dbx.InsertArgs, n int, teams []*teamimport.Model) (insert []*accountimport.Model, err error) {
	insert = []*accountimport.Model{}
	// generate a random set of accounts
	for i := 0; i < n; i++ {
		var envI = rand.IntN(len(environmentList))
		var teamI = rand.IntN(len(teams))
		var id = fmt.Sprintf("%04d", i+1)
		insert = append(insert, &accountimport.Model{
			ID:          id,
			Name:        fmt.Sprintf("Account %s", id),
			Label:       fmt.Sprintf("%d", i+1),
			Environment: environmentList[envI],
			TeamName:    teams[teamI].Name,
		})
	}

	err = dbx.Insert(ctx, accountimport.InsertStatement, insert, in)
	return

}

// seedTeams generates and inserts cost data similar to real life values
func seedTeams(ctx context.Context, in *dbx.InsertArgs) (insert []*teamimport.Model, err error) {
	insert = []*teamimport.Model{}
	for _, team := range teamList {
		insert = append(insert, &teamimport.Model{Name: team})
	}
	err = dbx.Insert(ctx, teamimport.InsertStatement, insert, in)

	return
}
