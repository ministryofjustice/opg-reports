package lib

import (
	"github.com/ministryofjustice/opg-reports/internal/endpoints"
	"github.com/ministryofjustice/opg-reports/internal/navigation"
	"github.com/ministryofjustice/opg-reports/servers/inout"
)

// AwsAccounts
const AwsAccountsList endpoints.ApiEndpoint = "/{version}/aws/accounts/list"

const (
	AwsCostsTotal endpoints.ApiEndpoint = "/{version}/aws/costs/total/{billing_date:-5}/{billing_date:1}" // Fetches the overal total between the dates (ex. Tax)
	AwsCostsList  endpoints.ApiEndpoint = "/{version}/aws/costs/list/{billing_date:-5}/{billing_date:1}"  // Returns all database entries between the dates (ex. Tax)

	AwsCostsMonthTaxes       endpoints.ApiEndpoint = "/{version}/aws/costs/tax/month/{billing_date:-5}/{billing_date:1}"              // Returns monthly totals with and without tax - can be filtered by unit
	AwsCostsMonthSum         endpoints.ApiEndpoint = "/{version}/aws/costs/sum/month/{billing_date:-5}/{billing_date:1}"              // Returns monthly totals without tax - can be filtered by unit
	AwsCostsMonthSumUnit     endpoints.ApiEndpoint = "/{version}/aws/costs/sum-per-unit/month/{billing_date:-5}/{billing_date:1}"     // Returns monthly totals group by unit
	AwsCostsMonthSumUnitEnv  endpoints.ApiEndpoint = "/{version}/aws/costs/sum-per-unit-env/month/{billing_date:-5}/{billing_date:1}" // Returns monthly totals grouped by the unit and account environment
	AwsCostsMonthSumDetailed endpoints.ApiEndpoint = "/{version}/aws/costs/sum-detailed/month/{billing_date:-5}/{billing_date:1}"     // Returns costs grouped by month, unit, account number, account environment, aws service and aws region - can be filtered by unit

	AwsCostsDayTaxes       endpoints.ApiEndpoint = "/{version}/aws/costs/tax/day/{billing_date:0}/{billing_date:1}"              // Returns daily totals with and without tax - can be filtered by unit
	AwsCostsDaySum         endpoints.ApiEndpoint = "/{version}/aws/costs/sum/day/{billing_date:0}/{billing_date:1}"              // Returns daily totals without tax - can be filtered by unit
	AwsCostsDaySumUnit     endpoints.ApiEndpoint = "/{version}/aws/costs/sum-per-unit/day/{billing_date:0}/{billing_date:1}"     // Returns daily totals group by unit
	AwsCostsDaySumUnitEnv  endpoints.ApiEndpoint = "/{version}/aws/costs/sum-per-unit-env/day/{billing_date:0}/{billing_date:1}" // Returns daily totals grouped by the unit and account environment
	AwsCostsDaySumDetailed endpoints.ApiEndpoint = "/{version}/aws/costs/sum-detailed/day/{billing_date:0}/{billing_date:1}"     // Returns costs grouped by day, unit, account number, account environment, aws service and aws region - can be filtered by unit

)

const (
	AwsUptimeList endpoints.ApiEndpoint = "/{version}/aws/uptime/list/{month:-6}/{month:0}" // Returns all uptime data between dates

	AwsUptimeMonthAverage     endpoints.ApiEndpoint = "/{version}/aws/uptime/average/month/{month:-6}/{month:0}"          // Returns average uptime % between dates
	AwsUptimeMonthAverageUnit endpoints.ApiEndpoint = "/{version}/aws/uptime/average-per-unit/month/{month:-6}/{month:0}" // Returns average uptime % between dates grouped by unit

	AwsUptimeDayAverage     endpoints.ApiEndpoint = "/{version}/aws/uptime/average/day/{day:-14}/{day:0}"          // Returns average uptime % between dates
	AwsUptimeDayAverageUnit endpoints.ApiEndpoint = "/{version}/aws/uptime/average-per-unit/day/{day:-14}/{day:0}" // Returns average uptime % between dates grouped by unit
)

const (
	GitHubReleaseList endpoints.ApiEndpoint = "/{version}/github/release/list/{month:-6}/{month:0}" // Return all releases between the dates

	GitHubReleaseMonthCount     endpoints.ApiEndpoint = "/{version}/github/release/count/month/{month:-6}/{month:0}"          // Return count of all releases grouped by month
	GitHubReleaseMonthCountUnit endpoints.ApiEndpoint = "/{version}/github/release/count-per-unit/month/{month:-6}/{month:0}" // Return count of releases grouped by month and unit

	GitHubReleaseDayCount     endpoints.ApiEndpoint = "/{version}/github/release/count/day/{day:-14}/{day:0}"          // Return count of all releases grouped by day
	GitHubReleaseDayCountUnit endpoints.ApiEndpoint = "/{version}/github/release/count-per-unit/day/{day:-14}/{day:0}" // Return count of releases grouped by day and unit
)

const GitHubRepositoryList endpoints.ApiEndpoint = "/{version}/github/repository/list" // Return list of all repositories

const GitHubRespoitoryStandardList endpoints.ApiEndpoint = "/{version}/github/standard/list" // Return list of all repository standards

const GitHubTeamList endpoints.ApiEndpoint = "/{version}/github/team/list" // Return all github teams

const UnitList endpoints.ApiEndpoint = "/{version}/unit/list" // Return all units

// -- Costs navigation items

// costsTaxOverview config
var costsTaxOverview = navigation.New(
	"Tax overview",
	"/costs/tax-overview",
	&navigation.Display{PageTemplate: "costs-tax"},
	&navigation.Data{
		Source:      AwsCostsMonthTaxes,
		Namespace:   "CostsTax",
		Body:        &inout.AwsCostsTaxesBody{},
		Transformer: inout.Transform,
	})

// costsPerTeam config
var costsPerTeam = navigation.New(
	"Costs per team",
	"/costs/unit",
	&navigation.Display{PageTemplate: "costs-unit"},
	&navigation.Data{
		Source:      AwsCostsMonthSumUnit,
		Namespace:   "CostsPerUnit",
		Body:        &inout.AwsCostsSumPerUnitBody{},
		Transformer: inout.Transform,
	},
	&navigation.Data{
		Source:      AwsCostsMonthSumUnitEnv,
		Namespace:   "CostsPerUnitEnv",
		Body:        &inout.AwsCostsSumPerUnitEnvBody{},
		Transformer: inout.Transform,
	},
)

// costsDetailed config
var costsDetailed = navigation.New(
	"Detailed costs",
	"/costs/detailed",
	&navigation.Display{PageTemplate: "costs-detailed"},
	&navigation.Data{
		Source:      AwsCostsMonthSumDetailed,
		Namespace:   "CostsDetailed",
		Body:        &inout.AwsCostsSumFullDetailsBody{},
		Transformer: inout.Transform,
	},
)

// costs is the overall cost navigation block
var costs = navigation.New(
	"Costs",
	"/costs",
	&navigation.Display{PageTemplate: "costs-overview", IsHeader: true},
	costsTaxOverview,
	costsPerTeam,
	costsDetailed,
)

// -- Standards navigation items

// Github repo standards
var ghStandards = navigation.New(
	"Repositories",
	"/standards/repositories",
	&navigation.Display{PageTemplate: "standards-github-repositories"},
	&navigation.Data{
		Source:    GitHubRespoitoryStandardList,
		Namespace: "RepositoryStandards",
		Body:      &inout.GitHubRepositoryStandardsListBody{},
	},
)

// wrapping standards
var standard = navigation.New(
	"Standards",
	"/standards",
	&navigation.Display{PageTemplate: "standards-overview", IsHeader: true},
	ghStandards,
)

// -- Uptime

var uptimeAws = navigation.New(
	"Service uptime",
	"/uptime/aws",
	&navigation.Display{PageTemplate: "uptime-aws"},
	&navigation.Data{
		Source:      AwsUptimeMonthAverage,
		Namespace:   "UptimeOverall",
		Body:        &inout.AwsUptimeAveragesBody{},
		Transformer: inout.Transform,
	},
	&navigation.Data{
		Source:      AwsUptimeMonthAverageUnit,
		Namespace:   "UptimeUnit",
		Body:        &inout.AwsUptimeAveragesPerUnitBody{},
		Transformer: inout.Transform,
	},
)

var up = navigation.New(
	"Uptime",
	"/uptime",
	&navigation.Display{PageTemplate: "uptime-overview", IsHeader: true},
	uptimeAws,
)

// -- Releases

var releasePerMonth = navigation.New(
	"Per month",
	"/releases/monthly",
	&navigation.Display{PageTemplate: "releases-github"},
	&navigation.Data{
		Source:      GitHubReleaseMonthCount,
		Namespace:   "ReleasesOverallMonthly",
		Body:        &inout.GitHubReleasesCountBody{},
		Transformer: inout.Transform,
	},
	&navigation.Data{
		Source:      GitHubReleaseMonthCountUnit,
		Namespace:   "ReleasesUnitMonthly",
		Body:        &inout.GitHubReleasesCountPerUnitBody{},
		Transformer: inout.Transform,
	},
)

var release = navigation.New(
	"Releases",
	"/releases",
	&navigation.Display{PageTemplate: "releases-overview", IsHeader: true},
	releasePerMonth,
)

// TODO: use a real team name for sirius
// -- team navigation - sirius
var siriusHistorical = navigation.New(
	"Historical Data",
	"/sirius/month",
	&navigation.Display{PageTemplate: "team-historical"},
	&navigation.Data{
		Source:      AwsUptimeMonthAverage + "?unit=sirius",
		Namespace:   "TeamUptimeUnit",
		Body:        &inout.AwsUptimeAveragesBody{},
		Transformer: inout.Transform,
	},
	&navigation.Data{
		Source:      AwsCostsMonthSumDetailed + "?unit=sirius",
		Namespace:   "TeamCostsPerUnit",
		Body:        &inout.AwsCostsSumFullDetailsBody{},
		Transformer: inout.Transform,
	},
	// &navigation.Data{
	// 	Source:      GitHubReleaseDayCount + "?unit=Sirius",
	// 	Namespace:   "TeamReleases",
	// 	Body:        &inout.GitHubReleasesCountBody{},
	// 	Transformer: inout.Transform,
	// },
)

var sirius = navigation.New(
	"Sirius",
	"/sirius",
	&navigation.Display{PageTemplate: "team-overview", IsHeader: true},
	&navigation.Data{
		Source:      AwsUptimeDayAverage + "?unit=sirius",
		Namespace:   "TeamUptimeUnit",
		Body:        &inout.AwsUptimeAveragesBody{},
		Transformer: inout.Transform,
	},
	// &navigation.Data{
	// 	Source:      GitHubReleaseDayCount + "?unit=Sirius",
	// 	Namespace:   "TeamReleases",
	// 	Body:        &inout.GitHubReleasesCountBody{},
	// 	Transformer: inout.Transform,
	// },
	siriusHistorical,
)

// -- Full navigation structure
var overview = navigation.New(
	"Overview",
	"/",
	&navigation.Display{PageTemplate: "homepage"},
	up,
	release,
	costs,
	standard,
)

// -- simple navigation structure

// replica of ghStandards so it doesnt get parent structure attached
// as that will then render the sidebar navigation
var single = navigation.New(
	"Repositories",
	"/standards/repositories",
	&navigation.Display{PageTemplate: "standards-github-repositories"},
	&navigation.Data{
		Source:    GitHubRespoitoryStandardList,
		Namespace: "RepositoryStandards",
		Body:      &inout.GitHubRepositoryStandardsListBody{},
	},
)

// navigation setup - picked by bi.Mode
var (
	full   = []*navigation.Navigation{overview, sirius}
	simple = []*navigation.Navigation{single}
)

// NavigationChoices is the set of all navigation structures
// to share
// This is the then selected in the sfront by using
// bi.Navigation as the key for this map
// This allows the navigation to be changed at run time
var NavigationChoices = map[string][]*navigation.Navigation{
	// "simple": simple,
	"simple": simple,
	"full":   full,
}
