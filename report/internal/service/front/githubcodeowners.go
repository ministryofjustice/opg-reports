package front

import (
	"net/http"
	"opg-reports/report/internal/endpoints"
	"opg-reports/report/internal/repository/restr"
)

// GithubCodeOwner
type GithubCodeOwner struct {
	CodeOwner  string `json:"codeowner"`
	Repository string `json:"repository"`
	Team       string `json:"team"`
}

// apiResponseGithubCodeOwnersForTeam represents the api data structure returned
// by the github codeowners for team end point
//
// endpoint: `/v1/githubcodeowners/team/{team}`
type apiResponseGithubCodeOwnersForTeam struct {
	Count int                `json:"count,omityempty"`
	Data  []*GithubCodeOwner `json:"data"`
}

// apiResponseGithubCodeOwnersForCodeOwners represents the api data structure returned
//
// endpoint: `/v1/githubcodeowners/codeowners/{codeowner}`
type apiResponseGithubCodeOwnersForCodeOwners struct {
	Count int                `json:"count,omityempty"`
	Data  []*GithubCodeOwner `json:"data"`
}

type githubCodeOwnersPreCallF func(params map[string]string)

// GetGithubCodeOwnersForCodeOwners
func (self *Service) GetGithubCodeOwnersForCodeOwners(
	client restr.RepositoryRestGetter,
	request *http.Request,
	apiParameters map[string]string,
	adjusters ...githubCodeOwnersPreCallF) (owners []*GithubCodeOwner, err error) {
	var (
		log      = self.log.With("operation", "GetGithubCodeOwnersForCodeOwners")
		defaults = githubCodeOwnerParams()
		params   = mergeRequestWithMaps(request, defaults, apiParameters)
		endpoint string
	)
	// allow function to overwrite parameters
	for _, adjustF := range adjusters {
		adjustF(params)
	}
	endpoint = endpoints.Parse(endpoints.GITHUBCODEOWNERS_FOR_CODEOWNER, params)

	log.With("defaults", defaults, "api", apiParameters, "merged", params).Debug("calling api for github codeowner codeowners")
	owners, err = getFromAPI(self.ctx, self.log, self.conf,
		client,
		endpoint,
		parseGithubCodeOwnersForCodeOwnersF,
	)

	log.Debug("returning api data ... ")
	return
}

// GetGithubCodeOwnersForTeam
func (self *Service) GetGithubCodeOwnersForTeam(
	client restr.RepositoryRestGetter,
	request *http.Request,
	apiParameters map[string]string,
	adjusters ...githubCodeOwnersPreCallF) (owners []*GithubCodeOwner, err error) {
	var (
		log      = self.log.With("operation", "GetGithubCodeOwnersForTeam")
		defaults = githubCodeOwnerParams()
		params   = mergeRequestWithMaps(request, defaults, apiParameters)
		endpoint string
	)
	// allow function to overwrite parameters
	for _, adjustF := range adjusters {
		adjustF(params)
	}
	endpoint = endpoints.Parse(endpoints.GITHUBCODEOWNERS_FOR_TEAM, params)

	log.With("defaults", defaults, "api", apiParameters, "merged", params).Debug("calling api for github codeowner team data")
	owners, err = getFromAPI(self.ctx, self.log, self.conf,
		client,
		endpoint,
		parseGithubCodeOwnersForTeamF,
	)

	log.Debug("returning api data ... ")
	return
}

func parseGithubCodeOwnersForTeamF(response *apiResponseGithubCodeOwnersForTeam) (data []*GithubCodeOwner, err error) {
	data = response.Data
	return
}

func parseGithubCodeOwnersForCodeOwnersF(response *apiResponseGithubCodeOwnersForCodeOwners) (data []*GithubCodeOwner, err error) {
	data = response.Data
	return
}

// githubCodeOwnerParams returns a map of the values for the input struct.
func githubCodeOwnerParams() map[string]string {

	return map[string]string{
		"team": "",
	}
}
