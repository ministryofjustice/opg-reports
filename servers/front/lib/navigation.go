package lib

import (
	"fmt"
	"slices"
	"strings"

	"github.com/ministryofjustice/opg-reports/internal/dateformats"
	"github.com/ministryofjustice/opg-reports/internal/dateintervals"
	"github.com/ministryofjustice/opg-reports/internal/dateutils"
	"github.com/ministryofjustice/opg-reports/internal/endpoints"
	"github.com/ministryofjustice/opg-reports/internal/navigation"
	"github.com/ministryofjustice/opg-reports/internal/tmplfuncs"
	"github.com/ministryofjustice/opg-reports/servers/inout"
)

const (
	AwsAccountsList              endpoints.ApiEndpoint = "/{version}/aws/accounts/list"                                                             // Lists all aws accounts
	AwsCostsTotal                endpoints.ApiEndpoint = "/{version}/aws/costs/total/{start_billing_date:-5}/{end_billing_date:1}"                  // Returns the overal total between the dates (ex. Tax)
	AwsCostsList                 endpoints.ApiEndpoint = "/{version}/aws/costs/list/{start_billing_date:-5}/{end_billing_date:1}"                   // Returns all database entries between the dates (ex. Tax)
	AwsCostsMonthTaxes           endpoints.ApiEndpoint = "/{version}/aws/costs/tax/month/{start_billing_date:-5}/{end_billing_date:1}"              // Returns monthly totals with and without tax - can be filtered by unit
	AwsCostsMonthSum             endpoints.ApiEndpoint = "/{version}/aws/costs/sum/month/{start_billing_date:-5}/{end_billing_date:1}"              // Returns monthly totals without tax - can be filtered by unit
	AwsCostsMonthSumUnit         endpoints.ApiEndpoint = "/{version}/aws/costs/sum-per-unit/month/{start_billing_date:-5}/{end_billing_date:1}"     // Returns monthly totals group by unit
	AwsCostsMonthSumUnitEnv      endpoints.ApiEndpoint = "/{version}/aws/costs/sum-per-unit-env/month/{start_billing_date:-5}/{end_billing_date:1}" // Returns monthly totals grouped by the unit and account environment
	AwsCostsMonthSumDetailed     endpoints.ApiEndpoint = "/{version}/aws/costs/sum-detailed/month/{start_billing_date:-5}/{end_billing_date:1}"     // Returns costs grouped by month, unit, account number, account environment, aws service and aws region - can be filtered by unit
	AwsCostsDayTaxes             endpoints.ApiEndpoint = "/{version}/aws/costs/tax/day/{start_billing_date:0}/{end_billing_date:1}"                 // Returns daily totals with and without tax - can be filtered by unit
	AwsCostsDaySum               endpoints.ApiEndpoint = "/{version}/aws/costs/sum/day/{start_billing_date:0}/{end_billing_date:1}"                 // Returns daily totals without tax - can be filtered by unit
	AwsCostsDaySumUnit           endpoints.ApiEndpoint = "/{version}/aws/costs/sum-per-unit/day/{start_billing_date:0}/{end_billing_date:1}"        // Returns daily totals group by unit
	AwsCostsDaySumUnitEnv        endpoints.ApiEndpoint = "/{version}/aws/costs/sum-per-unit-env/day/{start_billing_date:0}/{end_billing_date:1}"    // Returns daily totals grouped by the unit and account environment
	AwsCostsDaySumDetailed       endpoints.ApiEndpoint = "/{version}/aws/costs/sum-detailed/day/{start_billing_date:0}/{end_billing_date:1}"        // Returns costs grouped by day, unit, account number, account environment, aws service and aws region - can be filtered by unit
	AwsUptimeList                endpoints.ApiEndpoint = "/{version}/aws/uptime/list/{start_month:-6}/{end_month:0}"                                // Returns all uptime data between dates
	AwsUptimeMonthAverage        endpoints.ApiEndpoint = "/{version}/aws/uptime/average/month/{start_month:-6}/{end_month:0}"                       // Returns average uptime % between dates
	AwsUptimeMonthAverageUnit    endpoints.ApiEndpoint = "/{version}/aws/uptime/average-per-unit/month/{start_month:-6}/{end_month:0}"              // Returns average uptime % between dates grouped by unit
	AwsUptimeDayAverage          endpoints.ApiEndpoint = "/{version}/aws/uptime/average/day/{start_day:-7}/{end_day:0}"                             // Returns average uptime % between dates
	AwsUptimeDayAverageUnit      endpoints.ApiEndpoint = "/{version}/aws/uptime/average-per-unit/day/{start_day:-7}/{end_day:0}"                    // Returns average uptime % between dates grouped by unit
	GitHubReleaseList            endpoints.ApiEndpoint = "/{version}/github/release/list/{start_month:-6}/{end_month:0}"                            // Return all releases between the dates
	GitHubReleaseMonthCount      endpoints.ApiEndpoint = "/{version}/github/release/count/month/{start_month:-6}/{end_month:0}"                     // Return count of all releases grouped by month
	GitHubReleaseMonthCountUnit  endpoints.ApiEndpoint = "/{version}/github/release/count-per-unit/month/{start_month:-6}/{end_month:0}"            // Return count of releases grouped by month and unit
	GitHubReleaseDayCount        endpoints.ApiEndpoint = "/{version}/github/release/count/day/{start_day:-7}/{end_day:0}"                           // Return count of all releases grouped by day
	GitHubReleaseDayCountUnit    endpoints.ApiEndpoint = "/{version}/github/release/count-per-unit/day/{start_day:-7}/{end_day:0}"                  // Return count of releases grouped by day and unit
	GitHubRepositoryList         endpoints.ApiEndpoint = "/{version}/github/repository/list"                                                        // Return list of all repositories
	GitHubRespoitoryStandardList endpoints.ApiEndpoint = "/{version}/github/standard/list"                                                          // Return list of all repository standards
	GitHubTeamList               endpoints.ApiEndpoint = "/{version}/github/team/list"                                                              // Return all github teams
	UnitList                     endpoints.ApiEndpoint = "/{version}/unit/list"                                                                     // Return all units
)

// NavigationChoices is the set of all navigation structures
// to share and should then use info.Fixtures to choose
var NavigationChoices = map[string][]*navigation.Navigation{}

func overviewItems() (overview *navigation.Navigation) {
	// costsTaxOverview config
	costsTaxOverview := navigation.New(
		"Tax overview",
		"/costs/tax-overview",
		&navigation.Display{PageTemplate: "costs-tax"},
		&navigation.Data{
			Source:      AwsCostsMonthTaxes,
			Namespace:   "CostsTax",
			Body:        &inout.AwsCostsTaxesBody{},
			Transformer: inout.TransformToDateWideTable,
		})

	// costsPerTeam config
	costsPerTeam := navigation.New(
		"Costs per team",
		"/costs/unit",
		&navigation.Display{PageTemplate: "costs-unit"},
		&navigation.Data{
			Source:      AwsCostsMonthSumUnit,
			Namespace:   "CostsPerUnit",
			Body:        &inout.AwsCostsSumPerUnitBody{},
			Transformer: inout.TransformToDateWideTable,
		},
		&navigation.Data{
			Source:      AwsCostsMonthSumUnitEnv,
			Namespace:   "CostsPerUnitEnv",
			Body:        &inout.AwsCostsSumPerUnitEnvBody{},
			Transformer: inout.TransformToDateWideTable,
		},
	)

	// costsDetailed config
	costsDetailed := navigation.New(
		"Detailed costs",
		"/costs/detailed",
		&navigation.Display{PageTemplate: "costs-detailed"},
		&navigation.Data{
			Source:      AwsCostsMonthSumDetailed,
			Namespace:   "CostsDetailed",
			Body:        &inout.AwsCostsSumFullDetailsBody{},
			Transformer: inout.TransformToDateWideTable,
		},
	)

	// costs is the overall cost navigation block
	costs := navigation.New(
		"Costs",
		"/costs",
		&navigation.Display{PageTemplate: "costs-overview", IsHeader: true},
		costsTaxOverview,
		costsPerTeam,
		costsDetailed,
	)

	// -- Standards navigation items

	// Github repo standards
	ghStandards := navigation.New(
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
	standard := navigation.New(
		"Standards",
		"/standards",
		&navigation.Display{PageTemplate: "standards-overview", IsHeader: true},
		ghStandards,
	)

	// -- Uptime

	uptimeAws := navigation.New(
		"Service uptime",
		"/uptime/aws",
		&navigation.Display{PageTemplate: "uptime-aws"},
		&navigation.Data{
			Source:      AwsUptimeMonthAverage,
			Namespace:   "UptimeOverall",
			Body:        &inout.AwsUptimeAveragesBody{},
			Transformer: inout.TransformToDateWideTable,
		},
		&navigation.Data{
			Source:      AwsUptimeMonthAverageUnit,
			Namespace:   "UptimeUnit",
			Body:        &inout.AwsUptimeAveragesPerUnitBody{},
			Transformer: inout.TransformToDateWideTable,
		},
	)

	up := navigation.New(
		"Uptime",
		"/uptime",
		&navigation.Display{PageTemplate: "uptime-overview", IsHeader: true},
		uptimeAws,
	)

	// -- Releases

	releasePerMonth := navigation.New(
		"Per month",
		"/releases/monthly",
		&navigation.Display{PageTemplate: "releases-github"},
		&navigation.Data{
			Source:      GitHubReleaseMonthCount,
			Namespace:   "ReleasesOverallMonthly",
			Body:        &inout.GitHubReleasesCountBody{},
			Transformer: inout.TransformToDateWideTable,
		},
		&navigation.Data{
			Source:      GitHubReleaseMonthCountUnit,
			Namespace:   "ReleasesUnitMonthly",
			Body:        &inout.GitHubReleasesCountPerUnitBody{},
			Transformer: inout.TransformToDateWideTable,
		},
	)

	release := navigation.New(
		"Releases",
		"/releases",
		&navigation.Display{PageTemplate: "releases-overview", IsHeader: true},
		releasePerMonth,
	)
	overview = navigation.New(
		"Overview",
		"/",
		&navigation.Display{PageTemplate: "homepage"},
		up,
		release,
		costs,
		standard,
	)

	return
}

func teamItems() (navs []*navigation.Navigation) {
	var teams []string = []string{"digideps", "make", "modernise", "serve", "sirius", "use"}
	navs = []*navigation.Navigation{}

	for _, team := range teams {
		months := dateutils.Range(-4, 1, dateintervals.Month)
		slices.Reverse(months)
		unitFilter := endpoints.ApiEndpoint("?unit=" + team)

		monthlyOverview := navigation.New(
			"Overview",
			fmt.Sprintf("/%s/month/overview", team),
			&navigation.Display{PageTemplate: "team-per-month-overview"},
			&navigation.Data{
				Source:      AwsUptimeMonthAverage + unitFilter,
				Namespace:   "TeamUptimeUnit",
				Body:        &inout.AwsUptimeAveragesBody{},
				Transformer: inout.TransformToDateWideTable,
			},
			&navigation.Data{
				Source:      GitHubReleaseMonthCount + unitFilter,
				Namespace:   "TeamReleases",
				Body:        &inout.GitHubReleasesCountBody{},
				Transformer: inout.TransformToDateWideTable,
			},
			&navigation.Data{
				Source:      GitHubRepositoryList + unitFilter,
				Namespace:   "TeamRepositories",
				Body:        &inout.GitHubRepositoriesListBody{},
				Transformer: inout.TransformToDateWideTable,
			},
		)
		monthlyCosts := navigation.New(
			"Costs",
			fmt.Sprintf("/%s/month/costs", team),
			&navigation.Display{PageTemplate: "team-per-month-costs"},
			&navigation.Data{
				Source:      AwsCostsMonthSumDetailed + unitFilter,
				Namespace:   "TeamCostsPerUnit",
				Body:        &inout.AwsCostsSumFullDetailsBody{},
				Transformer: inout.TransformToDateWideTable,
			},
		)
		// -- uptime data per month
		monthly := navigation.New(
			"Monthly data",
			fmt.Sprintf("/%s/month", team),
			&navigation.Display{IsHeader: true},
			monthlyOverview,
			monthlyCosts,
		)

		uptimes := []*navigation.Navigation{}
		for _, m := range months {
			ymd := m.Format(dateformats.YMD)
			ym := m.Format(dateformats.YM)
			end := dateutils.Reset(m.AddDate(0, 1, 0), dateintervals.Month).Format(dateformats.YMD)

			uri := string(AwsUptimeDayAverage)
			uri = strings.ReplaceAll(uri, "{start_day:-7}", ymd)
			uri = strings.ReplaceAll(uri, "{end_day:0}", end)

			n := navigation.New(
				m.Format(dateformats.YM),
				fmt.Sprintf("/%s/uptime/day/%s", team, ym),
				&navigation.Display{PageTemplate: "team-uptime-day"},
				&navigation.Data{
					Source:      endpoints.ApiEndpoint(uri) + unitFilter,
					Namespace:   "TeamUptimeUnit",
					Body:        &inout.AwsUptimeAveragesBody{},
					Transformer: inout.TransformToDateDeepTable,
				},
			)
			uptimes = append(uptimes, n)
		}

		uptime := navigation.New(
			"Uptime",
			fmt.Sprintf("/%s/uptime", team),
			&navigation.Display{IsHeader: true},
			uptimes,
		)

		teamNav := navigation.New(
			tmplfuncs.Title(team),
			fmt.Sprintf("/%s", team),
			&navigation.Display{PageTemplate: "team-overview"},
			&navigation.Data{
				Source:      AwsUptimeDayAverage + unitFilter,
				Namespace:   "TeamUptimeUnit",
				Body:        &inout.AwsUptimeAveragesBody{},
				Transformer: inout.TransformToDateWideTable,
			},
			&navigation.Data{
				Source:      GitHubReleaseDayCount + unitFilter,
				Namespace:   "TeamReleases",
				Body:        &inout.GitHubReleasesCountBody{},
				Transformer: inout.TransformToDateWideTable,
			},
			&navigation.Data{
				Source:      GitHubRepositoryList + unitFilter,
				Namespace:   "TeamRepositories",
				Body:        &inout.GitHubRepositoriesListBody{},
				Transformer: inout.TransformToDateWideTable,
			},
			monthly,
			uptime,
		)
		navs = append(navs, teamNav)
	}

	return
}

func init() {

	overview := overviewItems()
	teamNavs := teamItems()

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

	full := []*navigation.Navigation{overview}
	simple := []*navigation.Navigation{single}
	full = append(full, teamNavs...)

	// NavigationChoices is the set of all navigation structures
	// to share and should then use info.Fixtures to choose
	NavigationChoices = map[string][]*navigation.Navigation{
		"simple": simple,
		"full":   full,
	}
}
