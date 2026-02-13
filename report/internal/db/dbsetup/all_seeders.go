package dbsetup

import (
	"fmt"
	"math/rand/v2"
	"opg-reports/report/internal/domain/accounts/accountmodels"
	"opg-reports/report/internal/domain/codebases/codebasemodels"
	"opg-reports/report/internal/domain/codeowners/codeownermodels"
	"opg-reports/report/internal/domain/infracosts/infracostmodels"
	"opg-reports/report/internal/domain/teams/teammodels"
	"opg-reports/report/internal/domain/uptime/uptimemodels"
	"opg-reports/report/internal/utils/times"
	"time"
)

// environments are shared account environment names to use in seeded data
var environments []string = []string{
	"development",
	"pre-production",
	"integrations",
	"production",
}

// fix regions to used for seeds
var regions []string = []string{
	"eu-west-1",
	"eu-west-2",
	"us-east-1",
	"NoRegion",
}

// aws service that we use for seeding
var services []string = []string{
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

// org name to use in mocks
var githubOrg string = "mock-gh-org"

func generateTeams(n int) (data []*teammodels.Team) {

	data = []*teammodels.Team{}
	for i := 0; i < n; i++ {
		data = append(data, &teammodels.Team{
			Name: fmt.Sprintf("TEAM-%03d", i+1),
		})
	}
	return
}

func generateAccounts(n int, teams []*teammodels.Team) (data []*accountmodels.Account) {
	data = []*accountmodels.Account{}

	for i := 0; i < n; i++ {
		var envI = rand.IntN(len(environments))
		var teamI = rand.IntN(len(teams))
		var id = fmt.Sprintf("%04d", i+1)

		data = append(data, &accountmodels.Account{
			ID:          id,
			Name:        fmt.Sprintf("Account %s", id),
			Label:       fmt.Sprintf("%d", i+1),
			Environment: environments[envI],
			TeamName:    teams[teamI].Name,
		})

	}
	return data

}

func generateInfracosts(n int, accounts []*accountmodels.Account) (data []*infracostmodels.Cost) {
	var (
		end    = times.ResetMonth(times.Add(time.Now().UTC(), -1, times.MONTH))
		start  = times.ResetMonth(times.Add(end, -10, times.YEAR))
		months = times.Months(start, end)
	)
	data = []*infracostmodels.Cost{}

	for i := 0; i < n; i++ {
		var accountI = rand.IntN(len(accounts))
		var monthI = rand.IntN(len(months))
		var regionI = rand.IntN(len(regions))
		var serviceI = rand.IntN(len(services))
		var price float64 = (-1000.0) + (rand.Float64() * (1000 - -1000.0)) // 95-100%

		data = append(data, &infracostmodels.Cost{
			Region:    regions[regionI],
			Service:   services[serviceI],
			Date:      times.AsYMDString(months[monthI]),
			Cost:      fmt.Sprintf("%g", price),
			AccountID: accounts[accountI].ID,
		})

	}

	return
}

func generateUptime(n int, accounts []*accountmodels.Account) (data []*uptimemodels.Uptime) {
	var (
		end   = times.ResetMonth(times.Add(time.Now().UTC(), -1, times.MONTH))
		start = times.ResetMonth(times.Add(end, -2, times.YEAR))
		days  = times.Days(start, end)
	)
	data = []*uptimemodels.Uptime{}

	for i := 0; i < n; i++ {
		var accountI = rand.IntN(len(accounts))
		var dayI = rand.IntN(len(days))
		var avg float64 = (95) + (rand.Float64() * (100 - 95)) // 95-100%

		data = append(data, &uptimemodels.Uptime{
			Date:      times.AsYMDString(days[dayI]),
			AccountID: accounts[accountI].ID,
			Average:   fmt.Sprintf("%g", avg),
		})

	}

	return
}

func generateCodebases(n int) (data []*codebasemodels.Codebase) {

	data = []*codebasemodels.Codebase{}
	for i := 0; i < n; i++ {
		var name = fmt.Sprintf("code-repository-%02d", i+1)
		data = append(data, &codebasemodels.Codebase{
			Name:     name,
			FullName: fmt.Sprintf("%s/%s", githubOrg, name),
			Url:      fmt.Sprintf("https://mock-github.local/%s/%s", githubOrg, name),
		})
	}

	return
}

func generateCodeowners(n int, teams []*teammodels.Team, codebases []*codebasemodels.Codebase) (data []*codeownermodels.Codeowner) {
	data = []*codeownermodels.Codeowner{}

	for i := 0; i < n; i++ {
		var teamI = rand.IntN(len(teams))
		var codeI = rand.IntN(len(codebases))

		data = append(data, &codeownermodels.Codeowner{
			TeamName:         teams[teamI].Name,
			CodebaseFullName: codebases[codeI].FullName,
			Name:             fmt.Sprintf("codeowner-%03d", i+1),
		})

	}
	return
}
