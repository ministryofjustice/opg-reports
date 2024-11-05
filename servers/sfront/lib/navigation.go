package lib

import (
	"github.com/ministryofjustice/opg-reports/pkg/endpoints"
	"github.com/ministryofjustice/opg-reports/pkg/navigation"
	"github.com/ministryofjustice/opg-reports/sources/costs/costsfront"
	"github.com/ministryofjustice/opg-reports/sources/costs/costsio"
	"github.com/ministryofjustice/opg-reports/sources/standards/standardsio"
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

// -- Costs navigation items

// costsTaxOverview config
var costsTaxOverview = navigation.New(
	"Tax overview",
	"/costs/tax-overview",
	&navigation.Display{PageTemplate: "costs-tax"},
	&navigation.Data{
		Source:      CostsUriMonthlyTax,
		Namespace:   "CostsTax",
		Body:        &costsio.TaxOverviewBody{},
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
		Body:        &costsio.StandardBody{},
		Transformer: costsfront.TransformResult,
	},
	&navigation.Data{
		Source:      CostsUriMonthlyUnitEnvironment,
		Namespace:   "CostsPerUnitEnv",
		Body:        &costsio.StandardBody{},
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
		Body:        &costsio.StandardBody{},
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
		Body:      &standardsio.Body{},
	},
)

// wrapping standards
var standard = navigation.New(
	"Standards",
	"/standards",
	&navigation.Display{PageTemplate: "standards-overview", IsHeader: true},
	ghStandards,
)

// -- simple navigation structure
// replica of ghStandards so it doesnt get parent structure attached
// as that will then render the sidebar navigation
var simple = navigation.New(
	"Repositories",
	"/standards/repositories",
	&navigation.Display{PageTemplate: "standards-github-repositories"},
	&navigation.Data{
		Source:    StandardsUri,
		Namespace: "RepositoryStandards",
		Body:      &standardsio.Body{},
	},
)

// -- Full navigation structure
var overview = navigation.New(
	"Overview",
	"/",
	&navigation.Display{PageTemplate: "homepage"},
	standard,
	costs,
)

// NavigationChoices is the set of all navigation structures
// to share
// This is the then selected in the sfront by using
// bi.Navigation as the key for this map
// This allows the navigation to be changed at run time
var NavigationChoices = map[string][]*navigation.Navigation{
	// "simple": {simple},
	"simple": {overview},
}