package navigation

import (
	"github.com/ministryofjustice/opg-reports/sources/costs/costsapi"
)

// -- Costs navigation items

// costsTaxOverview config
var costsTaxOverview = New(
	"Tax overview",
	"/costs/tax-overview",
	&Display{PageTemplate: "costs-tax"},
	&Data{Source: costsapi.UriMonthlyTax, Namespace: "CostsTax"},
)

// costsPerTeam config
var costsPerTeam = New(
	"Costs per team",
	"/costs/unit",
	&Display{PageTemplate: "costs-unit"},
	&Data{Source: costsapi.UriMonthlyUnit, Namespace: "CostsPerUnit"},
	&Data{Source: costsapi.UriMonthlyUnitEnvironment, Namespace: "CostsPerUnitEnv"},
)

// costsDetailed config
var costsDetailed = New(
	"Detailed costs",
	"/costs/detailed",
	&Display{PageTemplate: "costs-detailed"},
	&Data{Source: costsapi.UriMonthlyDetailed, Namespace: "CostsDetailed"},
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
