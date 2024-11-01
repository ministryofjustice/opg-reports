package navigation

import (
	"github.com/ministryofjustice/opg-reports/sources/costs/costsapi"
	"github.com/ministryofjustice/opg-reports/sources/costs/costsfront"
)

// -- Costs navigation items

// costsTaxOverview config
var costsTaxOverview = New(
	"Tax overview",
	"/costs/tax-overview",
	&Display{PageTemplate: "costs-tax"},
	&Data{
		Source:    costsapi.UriMonthlyTax,
		Namespace: "CostsTax",
		Body:      &costsapi.TaxOverviewBody{},
	},
)

// costsPerTeam config
var costsPerTeam = New(
	"Costs per team",
	"/costs/unit",
	&Display{PageTemplate: "costs-unit"},
	&Data{
		Source:      costsapi.UriMonthlyUnit,
		Namespace:   "CostsPerUnit",
		Body:        &costsapi.StandardBody{},
		Transformer: costsfront.TransformResult,
	},
	&Data{
		Source:      costsapi.UriMonthlyUnitEnvironment,
		Namespace:   "CostsPerUnitEnv",
		Body:        &costsapi.StandardBody{},
		Transformer: costsfront.TransformResult,
	},
)

// costsDetailed config
var costsDetailed = New(
	"Detailed costs",
	"/costs/detailed",
	&Display{PageTemplate: "costs-detailed"},
	&Data{
		Source:    costsapi.UriMonthlyDetailed,
		Namespace: "CostsDetailed",
		Body:      &costsapi.StandardBody{},
	},
)

// costs is the overall cost navigation block
var costs = New(
	"Costs",
	"/costs",
	&Display{PageTemplate: "costs-overview", IsHeader: true},
	costsTaxOverview,
	costsPerTeam,
	costsDetailed,
)

// Dummy holder for now
var simple = New(
	"Home",
	"/",
	&Display{PageTemplate: "homepage"},
)

// Configured is the set of all navigation structures
// to share
// This is the then selected in the sfront by using
// bi.Navigation as the key for this map
// This allows the navigation to be changed at run time
var Configured = map[string][]*Navigation{
	// "simple": {simple},
	"simple": {costsPerTeam, costs},
}
