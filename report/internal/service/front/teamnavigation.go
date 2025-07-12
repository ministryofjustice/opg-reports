package front

import (
	"net/http"
	"opg-reports/report/internal/endpoints"
	"opg-reports/report/internal/repository/restr"
	"slices"
)

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

// GetTeamNavigation calls the api (determined via values in conf) on the TEAMS_GET_ALL
// endpoint (`/v1/teams/all`), converts the response to a struct (*apiResponseTeams) and
// then parses the `data` slice, ignoring some team names, to provide a suitable list
// for the navigation structure. The list of teams is sorted alphabetically.
//
// Note: ignored team names are: "Legacy" & "ORG"
func (self *Service) GetTeamNavigation(client restr.RepositoryRestGetter, request *http.Request) (teams []string, err error) {
	var log = self.log.With("operation", "GetTeamNavigation")

	log.Debug("getting team navigation data ... ")
	teams, err = getFromAPI(self.ctx, log, self.conf,
		client,
		endpoints.TEAMS_GET_ALL,
		parseTeamNavigationF,
	)
	slices.Sort(teams)
	log.With("count", len(teams)).Debug("returning team navigation data ...")
	return
}

// parseTeamNavigationF handles the api response structure, parses out the teams
// list and returns them as slice of strings.
func parseTeamNavigationF(response *apiResponseTeams) (list []string, err error) {
	var exclude = []string{"Legacy", "ORG"}
	list = []string{}
	for _, team := range response.Data {
		if !slices.Contains(exclude, team.Name) {
			list = append(list, team.Name)
		}
	}
	return
}
