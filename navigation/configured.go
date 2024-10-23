package navigation

var costs *Navigation = &Navigation{
	Name:    "Costs",
	Uri:     "/costs",
	Display: &NavigationDisplay{IsHeader: true},
	Children: []*Navigation{
		{
			Name:    "Tax Overview",
			Uri:     "/costs/tax-overview",
			Display: &NavigationDisplay{},
		},
		{
			Name:    "Team costs",
			Uri:     "/costs/unit",
			Display: &NavigationDisplay{},
		},
		{
			Name:    "Detailed costs",
			Uri:     "/costs/detailed",
			Display: &NavigationDisplay{},
		},
	},
}

var team *Navigation = &Navigation{
	Name:    "{team}",
	Uri:     "/team/{team}",
	Display: &NavigationDisplay{IsHeader: true},
	Children: []*Navigation{
		{
			Name:    "Team costs",
			Uri:     "/team/costs",
			Display: &NavigationDisplay{},
		},
		{
			Name:    "Detailed costs",
			Uri:     "/costs/detailed",
			Display: &NavigationDisplay{},
		},
	},
}
