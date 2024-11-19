package lib

import (
	"github.com/ministryofjustice/opg-reports/pkg/endpoints"
	"github.com/ministryofjustice/opg-reports/pkg/navigation"
	"github.com/ministryofjustice/opg-reports/sources/costs/costsfront"
	"github.com/ministryofjustice/opg-reports/sources/costs/costsio"
	"github.com/ministryofjustice/opg-reports/sources/releases/releasesfront"
	"github.com/ministryofjustice/opg-reports/sources/releases/releasesio"
	"github.com/ministryofjustice/opg-reports/sources/standards/standardsio"
	"github.com/ministryofjustice/opg-reports/sources/uptime/uptimefront"
	"github.com/ministryofjustice/opg-reports/sources/uptime/uptimeio"
)

// Cost endpoints to call formatted with required placeholders
const (
	CostsUriTotal                  endpoints.ApiEndpoint = "/{version}/costs/aws/total/{billing_date:-11}/{billing_date:0}"
	CostsUriMonthlyTax             endpoints.ApiEndpoint = "/{version}/costs/aws/tax-overview/{billing_date:-11}/{billing_date:0}/month"
	CostsUriMonthlyUnit            endpoints.ApiEndpoint = "/{version}/costs/aws/unit/{billing_date:-9}/{billing_date:0}/month"
	CostsUriMonthlyUnitEnvironment endpoints.ApiEndpoint = "/{version}/costs/aws/unit-environment/{billing_date:-9}/{billing_date:0}/month"
	CostsUriMonthlyDetailed        endpoints.ApiEndpoint = "/{version}/costs/aws/detailed/{billing_date:-6}/{billing_date:0}/month"
	CostsUriDailyUnit              endpoints.ApiEndpoint = "/{version}/costs/aws/unit/{billing_date:-1}/{billing_date:0}/day"
	CostsUriDailyUnitEnvironment   endpoints.ApiEndpoint = "/{version}/costs/aws/unit-environment/{billing_date:-1}/{billing_date:0}/day"
	CostsUriDailyDetailed          endpoints.ApiEndpoint = "/{version}/costs/aws/detailed/{billing_date:-1}/{billing_date:0}/day"
)

// Standards endpoints
const (
	StandardsUri endpoints.ApiEndpoint = "/{version}/standards/github/false"
)

// Uptime endpoints
const (
	UptimeOverallMonthlyUri endpoints.ApiEndpoint = "/{version}/uptime/aws/overall/{month:-9}/{month:0}/month"
	UptimePerUnitMonthlyUri endpoints.ApiEndpoint = "/{version}/uptime/aws/unit/{month:-9}/{month:0}/month"
	UptimePerUnitBillingUri endpoints.ApiEndpoint = "/{version}/uptime/aws/unit/{billing_date:-6}/{billing_date:0}/month"
	UptimeOverallDailylUri  endpoints.ApiEndpoint = "/{version}/uptime/aws/overall/{day:-14}/{day:0}/day"
	UptimePerUnitDailylUri  endpoints.ApiEndpoint = "/{version}/uptime/aws/unit/{day:-14}/{day:0}/day"
)

// Releases endpoints
const (
	// /v1/releases/github/counts/2022-01-01/2024-04-01/month
	ReleasesOverallMonthly endpoints.ApiEndpoint = "/{version}/releases/github/counts/all/{month:-9}/{month:0}/month"
	ReleasesUnitMonthly    endpoints.ApiEndpoint = "/{version}/releases/github/counts/unit/{month:-9}/{month:0}/month"
)

// -- Costs navigation items

// costsTaxOverview config
var costsTaxOverview = navigation.New(
	"Tax overview",
	"/costs/tax-overview",
	&navigation.Display{PageTemplate: "costs-tax"},
	&navigation.Data{
		Source:      CostsUriMonthlyTax,
		Namespace:   "CostsTax",
		Body:        &costsio.CostsTaxOverviewBody{},
		Transformer: costsfront.TransformResult,
	})

// costsPerTeam config
var costsPerTeam = navigation.New(
	"Costs per team",
	"/costs/unit",
	&navigation.Display{PageTemplate: "costs-unit"},
	&navigation.Data{
		Source:      CostsUriMonthlyUnit,
		Namespace:   "CostsPerUnit",
		Body:        &costsio.CostsStandardBody{},
		Transformer: costsfront.TransformResult,
	},
	&navigation.Data{
		Source:      CostsUriMonthlyUnitEnvironment,
		Namespace:   "CostsPerUnitEnv",
		Body:        &costsio.CostsStandardBody{},
		Transformer: costsfront.TransformResult,
	},
)

// costsDetailed config
var costsDetailed = navigation.New(
	"Detailed costs",
	"/costs/detailed",
	&navigation.Display{PageTemplate: "costs-detailed"},
	&navigation.Data{
		Source:      CostsUriMonthlyDetailed,
		Namespace:   "CostsDetailed",
		Body:        &costsio.CostsStandardBody{},
		Transformer: costsfront.TransformResult,
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
		Source:    StandardsUri,
		Namespace: "RepositoryStandards",
		Body:      &standardsio.StandardsBody{},
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
		Source:      UptimeOverallMonthlyUri,
		Namespace:   "UptimeOverall",
		Body:        &uptimeio.UptimeBody{},
		Transformer: uptimefront.TransformResult,
	},
	&navigation.Data{
		Source:      UptimePerUnitMonthlyUri,
		Namespace:   "UptimeUnit",
		Body:        &uptimeio.UptimeBody{},
		Transformer: uptimefront.TransformResult,
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
		Source:      ReleasesOverallMonthly,
		Namespace:   "ReleasesOverallMonthly",
		Body:        &releasesio.ReleasesBody{},
		Transformer: releasesfront.TransformResult,
	},
	&navigation.Data{
		Source:      ReleasesUnitMonthly,
		Namespace:   "ReleasesUnitMonthly",
		Body:        &releasesio.ReleasesBody{},
		Transformer: releasesfront.TransformResult,
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
		Source:      UptimePerUnitBillingUri + "?unit=Sirius",
		Namespace:   "TeamUptimeUnit",
		Body:        &uptimeio.UptimeBody{},
		Transformer: uptimefront.TransformResult,
	},
	&navigation.Data{
		Source:      CostsUriMonthlyDetailed + "?unit=Sirius",
		Namespace:   "TeamCostsPerUnit",
		Body:        &costsio.CostsStandardBody{},
		Transformer: costsfront.TransformResult,
	},
)

var sirius = navigation.New(
	"Sirius",
	"/sirius",
	&navigation.Display{PageTemplate: "team-overview", IsHeader: true},
	&navigation.Data{
		Source:      UptimePerUnitDailylUri + "?unit=Sirius",
		Namespace:   "TeamUptimeUnit",
		Body:        &uptimeio.UptimeBody{},
		Transformer: uptimefront.TransformResult,
	},
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
		Source:    StandardsUri,
		Namespace: "RepositoryStandards",
		Body:      &standardsio.StandardsBody{},
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
