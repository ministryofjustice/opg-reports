package seed

import (
	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/utils"
)

// stmtGithubCodeOwnerSeed is used to insert records into the team table
// for seed / fixture data
const stmtGithubCodeOwnerSeed string = `
INSERT INTO github_codeownership (
	codeowner,
	repository,
	team
) VALUES (
	:codeowner,
	:repository,
	:team
) ON CONFLICT (codeowner,repository,team)
 	DO UPDATE SET team=excluded.team
RETURNING id;`

// teamSeed is used for seeding team data only
type githubCodeOwnerSeed struct {
	CodeOwner  string `json:"codeowner" db:"codeowner"`
	Repository string `json:"repository" db:"repository"`
	Team       string `json:"team,omitempty" db:"team"`
}

var githubCodeOwnerSeeds = []*sqlr.BoundStatement{
	// repo-a has 1 codeowner, that is shared by 3 teams
	{Statement: stmtGithubCodeOwnerSeed, Data: &githubCodeOwnerSeed{Team: "TEAM-A", Repository: "moj/repo-a", CodeOwner: "opg/gh-team-joint"}},
	{Statement: stmtGithubCodeOwnerSeed, Data: &githubCodeOwnerSeed{Team: "TEAM-B", Repository: "moj/repo-a", CodeOwner: "opg/gh-team-joint"}},
	{Statement: stmtGithubCodeOwnerSeed, Data: &githubCodeOwnerSeed{Team: "TEAM-C", Repository: "moj/repo-a", CodeOwner: "opg/gh-team-joint"}},
	// repo-b has 3 codeowners shared between 2 teams
	{Statement: stmtGithubCodeOwnerSeed, Data: &githubCodeOwnerSeed{Team: "TEAM-A", Repository: "moj/repo-b", CodeOwner: "opg/gh-team-b1"}},
	{Statement: stmtGithubCodeOwnerSeed, Data: &githubCodeOwnerSeed{Team: "TEAM-B", Repository: "moj/repo-b", CodeOwner: "opg/gh-team-b2"}},
	{Statement: stmtGithubCodeOwnerSeed, Data: &githubCodeOwnerSeed{Team: "TEAM-B", Repository: "moj/repo-b", CodeOwner: "opg/gh-team-b3"}},
	// repo-c has 1 codeowner and 1 team
	{Statement: stmtGithubCodeOwnerSeed, Data: &githubCodeOwnerSeed{Team: "TEAM-D", Repository: "moj/repo-c", CodeOwner: "opg/gh-team-c"}},
	// repo-d has a codeowner, but no team
	{Statement: stmtGithubCodeOwnerSeed, Data: &githubCodeOwnerSeed{Team: "NONE", Repository: "moj/repo-d", CodeOwner: "opg/gh-team-d"}},
}

// GithubCodeOwners populates the database (via the sqc var) with standard known enteries
// that can be used for testing and development databases
func (self *Service) GithubCodeOwners(sqc sqlr.RepositoryWriter) (results []*sqlr.BoundStatement, err error) {
	var sw = utils.Stopwatch()

	defer func() {
		self.log.With("seconds", sw.Stop().Seconds(), "inserted", len(results)).
			Info("[seed:CodeOwners] seeding finished.")
	}()
	self.log.Info("[seed:CodeOwners] starting seeding ...")
	err = sqc.Insert(githubCodeOwnerSeeds...)
	if err != nil {
		return
	}
	self.log.Info("[seed:CodeOwners] seeding successful")
	results = githubCodeOwnerSeeds
	return
}
