package navigation

import (
	"github.com/ministryofjustice/opg-reports/sources/costs/costsapi"
)

// org wide costs navigation items
var (
	taxOverview = &Navigation{
		Name: "Tax Overview",
		Uri:  "/costs/tax-overview",
		Display: &NavigationDisplay{
			IsHeader: false,
		},
		Data: []*NavigationData{
			{Source: costsapi.UriMonthlyTax},
		},
	}
	costsPerTeam = &Navigation{
		Name: "Costs per team",
		Uri:  "/costs/unit",
		Display: &NavigationDisplay{
			IsHeader: false,
		},
		Data: []*NavigationData{
			{Source: costsapi.UriMonthlyUnit},
			{Source: costsapi.UriMonthlyUnitEnvironment},
		},
	}
	detailedCosts = &Navigation{
		Name: "Detailed costs",
		Uri:  "/costs/detailed",
		Display: &NavigationDisplay{
			IsHeader: false,
		},
		Data: []*NavigationData{
			{Source: costsapi.UriMonthlyDetailed},
		},
	}
	costs *Navigation = &Navigation{
		Name:    "Costs",
		Uri:     "/costs",
		Display: &NavigationDisplay{IsHeader: true},
		Children: []*Navigation{
			taxOverview,
			costsPerTeam,
			detailedCosts,
		},
	}
)
