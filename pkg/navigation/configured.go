package navigation

import (
	"github.com/ministryofjustice/opg-reports/sources/costs/costsapi"
)

// org wide costs navigation items
var (
	taxOverview = &Navigation{
		Name: "Tax Overview",
		Uri:  "/costs/tax-overview",
		Display: &Display{
			IsHeader:     false,
			PageTemplate: "costs-tax",
		},
		Data: []*Data{
			{Source: costsapi.UriMonthlyTax},
		},
	}
	costsPerTeam = &Navigation{
		Name: "Costs per team",
		Uri:  "/costs/unit",
		Display: &Display{
			IsHeader:     false,
			PageTemplate: "costs-unit",
		},
		Data: []*Data{
			{Source: costsapi.UriMonthlyUnit},
			{Source: costsapi.UriMonthlyUnitEnvironment},
		},
	}
	detailedCosts = &Navigation{
		Name: "Detailed costs",
		Uri:  "/costs/detailed",
		Display: &Display{
			IsHeader:     false,
			PageTemplate: "costs-detailed",
		},
		Data: []*Data{
			{Source: costsapi.UriMonthlyDetailed},
		},
	}
	costs *Navigation = &Navigation{
		Name: "Costs",
		Uri:  "/costs",
		Display: &Display{
			IsHeader:     true,
			PageTemplate: "costs-overview",
		},
		Children: []*Navigation{
			taxOverview,
			costsPerTeam,
			detailedCosts,
		},
	}
)
