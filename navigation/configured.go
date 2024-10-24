package navigation

// org wide costs navigation items
var (
	taxOverview = &Navigation{
		Name:    "Tax Overview",
		Uri:     "/costs/tax-overview",
		Display: &NavigationDisplay{},
	}
	costsPerTeam = &Navigation{
		Name:    "Costs per team",
		Uri:     "/costs/unit",
		Display: &NavigationDisplay{},
	}
	detailedCosts = &Navigation{
		Name:    "Detailed costs",
		Uri:     "/costs/detailed",
		Display: &NavigationDisplay{},
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
