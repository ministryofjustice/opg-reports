package main

// apiResponseTeams captures name only data from the api response
// which is then used for generating the navigation structure
//
// endpoint: `/v1/teams/all`
type apiResponseTeams struct {
	Count int `json:"count,omityempty"`
	Data  []*struct {
		Name string `json:"name"`
	} `json:"data"`
}

// parseAllTeamsForNavigation excludes Legacy & ORG from team listing in the navigation
// for ease
func parseAllTeamsForNavigation(response *apiResponseTeams) (teams []string, err error) {
	teams = []string{}
	for _, team := range response.Data {
		if team.Name != "Legacy" && team.Name != "ORG" {
			teams = append(teams, team.Name)
		}
	}
	return
}

// apiAwsCost is flat version of cost data that comes back from the api
type apiAwsCost struct {
	Cost        string `json:"cost"`
	Date        string `json:"date,omitempty"`
	Region      string `json:"region,omitempty"`
	Service     string `json:"service,omitempty"`
	TeamName    string `json:"team_name,omitempty"`
	AccountID   string `json:"aws_account_id,omitempty"`
	Account     string `json:"aws_account_name,omitempty"`
	Label       string `json:"aws_account_label,omitempty"`
	Environment string `json:"aws_account_environment,omitempty"`
}

// apiResponseAwsCostsGrouped
type apiResponseAwsCostsGrouped struct {
	Count   int `json:"count,omityempty"`
	Request *struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	} `json:"request"`
	Data []*apiAwsCost `json:"data"`
}

func parseAwsCostsGrouped(response *apiResponseAwsCostsGrouped) (rows []map[string]string, err error) {

	return
}
