package seed

import (
	"github.com/ministryofjustice/opg-reports/report/internal/repository/sqlr"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

// stmtTeamSeed is used to insert records into the team table
// for seed / fixture data
const stmtTeamSeed string = `
INSERT INTO teams (
	name
) VALUES (
	:name
) ON CONFLICT (name)
 	DO UPDATE SET name=excluded.name
RETURNING name;`

// teamSeed is used for seeding team data only
type teamSeed struct {
	Name string `json:"name,omitempty" db:"name""`
}

var teamSeeds = []*sqlr.BoundStatement{
	{Statement: stmtTeamSeed, Data: &teamSeed{Name: "TEAM-A"}},
	{Statement: stmtTeamSeed, Data: &teamSeed{Name: "TEAM-B"}},
	{Statement: stmtTeamSeed, Data: &teamSeed{Name: "TEAM-C"}},
	{Statement: stmtTeamSeed, Data: &teamSeed{Name: "TEAM-D"}},
	{Statement: stmtTeamSeed, Data: &teamSeed{Name: "TEAM-E"}},
}

// Teams populates the database (via the sqc var) with standard known enteries
// that can be used for testing and development databases
func (self *Service) Teams(sqc sqlr.Writer) (results []*sqlr.BoundStatement, err error) {
	var sw = utils.Stopwatch()

	defer func() {
		self.log.With("seconds", sw.Stop().Seconds(), "inserted", len(results)).
			Info("[seed:Teams] seeding finished.")
	}()
	self.log.Info("[seed:Teams] starting seeding ...")
	err = sqc.Insert(teamSeeds...)
	if err != nil {
		return
	}
	self.log.Info("[seed:Teams] seeding successful")
	results = teamSeeds
	return
}
