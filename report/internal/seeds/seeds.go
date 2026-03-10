package seeds

import (
	"context"
	"fmt"
	"math/rand/v2"
	accountsimport "opg-reports/report/internal/domains/account/importer"
	accounts "opg-reports/report/internal/domains/account/types"
	codebasesimport "opg-reports/report/internal/domains/code/codebases/importer"
	codebases "opg-reports/report/internal/domains/code/types"
	teamsimport "opg-reports/report/internal/domains/team/importer"
	teams "opg-reports/report/internal/domains/team/types"
	"opg-reports/report/internal/migrations"
	"opg-reports/report/packages/args"
	"opg-reports/report/packages/dbx"
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

// Seeded contains all the seed data that was inserted
// including any that may have failed
type Seeded struct {
	Teams     []*teams.Team             `json:"teams"`
	Accounts  []*accounts.ImportAccount `json:"accounts"`
	Codebases []*codebases.Codebase     `json:"codebases"`
}

// quantities
var (
	numAccounts  int = 25
	numCodebases int = 50
)

func Seed(ctx context.Context, opts *args.DB) (seeded *Seeded, err error) {

	seeded = &Seeded{}
	// run the migrations
	err = migrations.Migrate(ctx, opts)
	if err != nil {
		return
	}
	// seed teams
	seeded.Teams, err = seedTeams(ctx, opts)
	if err != nil {
		return
	}
	// seed accounts
	seeded.Accounts, err = seedAccounts(ctx, numAccounts, seeded.Teams, opts)
	if err != nil {
		return
	}
	// seed codebases
	seeded.Codebases, err = seedCodebases(ctx, numCodebases, opts)
	if err != nil {
		return
	}

	return
}

func seedCodebases(ctx context.Context, n int, opts *args.DB) (models []*codebases.Codebase, err error) {
	var org = "mock-org"
	models = []*codebases.Codebase{}

	for i := 0; i < n; i++ {
		var name = fmt.Sprintf("codebase-%02d", i+1)

		models = append(models, &codebases.Codebase{
			Name:     name,
			FullName: fmt.Sprintf("%s/%s", org, name),
			Url:      fmt.Sprintf("https://mock-github.local/%s/%s", org, name),
			Archived: 0,
		})
	}
	err = dbx.Insert(ctx, codebasesimport.InsertStatement, models, opts)

	return
}

func seedAccounts(ctx context.Context, n int, list []*teams.Team, opts *args.DB) (models []*accounts.ImportAccount, err error) {
	models = []*accounts.ImportAccount{}
	// generate a random set of accounts
	for i := 0; i < n; i++ {
		var envI = rand.IntN(len(environmentList))
		var teamI = rand.IntN(len(list))
		var id = fmt.Sprintf("%04d", i+1)
		models = append(models, &accounts.ImportAccount{
			ID:          id,
			Name:        fmt.Sprintf("Account %s", id),
			Label:       fmt.Sprintf("%d", i+1),
			Environment: environmentList[envI],
			TeamName:    list[teamI].Name,
		})
	}

	err = dbx.Insert(ctx, accountsimport.InsertStatement, models, opts)
	return
}

func seedTeams(ctx context.Context, opts *args.DB) (models []*teams.Team, err error) {
	models = []*teams.Team{}
	for _, name := range teamList {
		models = append(models, &teams.Team{Name: name})
	}
	err = dbx.Insert(ctx, teamsimport.InsertStatement, models, opts)
	return
}
