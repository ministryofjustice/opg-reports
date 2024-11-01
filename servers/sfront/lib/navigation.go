package lib

import (
	"github.com/ministryofjustice/opg-reports/pkg/navigation"
	"github.com/ministryofjustice/opg-reports/sources/costs/costsapi"
	"github.com/ministryofjustice/opg-reports/sources/costs/costsfront"
)

// -- Costs navigation items

// costsTaxOverview config
var costsTaxOverview = navigation.New(
	"Tax overview",
	"/costs/tax-overview",
	&navigation.Display{PageTemplate: "costs-tax"},
	&navigation.Data{
		Source:      costsapi.UriMonthlyTax,
		Namespace:   "CostsTax",
		Body:        &costsapi.TaxOverviewBody{},
		Transformer: costsfront.TransformResult,
	},
)

// costsPerTeam config
var costsPerTeam = navigation.New(
	"Costs per team",
	"/costs/unit",
	&navigation.Display{PageTemplate: "costs-unit"},
	&navigation.Data{
		Source:      costsapi.UriMonthlyUnit,
		Namespace:   "CostsPerUnit",
		Body:        &costsapi.StandardBody{},
		Transformer: costsfront.TransformResult,
	},
	&navigation.Data{
		Source:      costsapi.UriMonthlyUnitEnvironment,
		Namespace:   "CostsPerUnitEnv",
		Body:        &costsapi.StandardBody{},
		Transformer: costsfront.TransformResult,
	},
)

// costsDetailed config
var costsDetailed = navigation.New(
	"Detailed costs",
	"/costs/detailed",
	&navigation.Display{PageTemplate: "costs-detailed"},
	&navigation.Data{
		Source:      costsapi.UriMonthlyDetailed,
		Namespace:   "CostsDetailed",
		Body:        &costsapi.StandardBody{},
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

// Dummy holder for now
var simple = navigation.New(
	"Home",
	"/",
	&navigation.Display{PageTemplate: "homepage"},
)

// NavigationChoices is the set of all navigation structures
// to share
// This is the then selected in the sfront by using
// bi.Navigation as the key for this map
// This allows the navigation to be changed at run time
var NavigationChoices = map[string][]*navigation.Navigation{
	// "simple": {simple},
	"simple": {costs, simple},
}
