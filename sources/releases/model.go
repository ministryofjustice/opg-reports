package releases

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/pkg/datastore"
	"github.com/ministryofjustice/opg-reports/pkg/record"
	"github.com/ministryofjustice/opg-reports/sources/releases/releasesdb"
)

// Release captures either the merged pull request
type Release struct {
	ID         int     `json:"id,omitempty" db:"id" faker:"unique, boundary_start=1, boundary_end=2000000" doc:"Database primary key."` // ID is a generated primary key
	Ts         string  `json:"ts,omitempty" db:"ts" faker:"time_string" doc:"Time the record was created."`                             // TS is timestamp when the record was created
	Repository string  `json:"repository,omitempty" db:"repository" faker:"word"`
	Name       string  `json:"name,omitempty" db:"name" faker:"word"`
	Source     string  `json:"source,omitempty" db:"source" doc:"url of source event." faker:"uri"`
	Type       string  `json:"type,omitempty" db:"type" faker:"oneof: workflow_run, pull_request" enum:"workflow,merge" doc:"tracks if this was a workflow run or a merge"`
	Date       string  `json:"date,omitempty" db:"date" faker:"date_string" doc:"Date this release happened."`
	Count      int     `json:"count,omitempty" db:"count" faker:"oneof: 1" enum:"1" doc:"Number of releases"`
	TeamList   []*Team `json:"teams,omitempty" db:"teams" faker:"slice_len=1" doc:"pulled from a many to many join table"`
}

// InsertJoins.
// On insert of a Release row we then want to check and convert any associated values
// within TeamList field (imported from json normally) into database joins
// To do that we:
//   - find each team with matching name (or insert a new one)
//   - find each (or create) the join between the team and this release
//
// RecordInsertJoiner interface
func (self *Release) InsertJoins(ctx context.Context, db *sqlx.DB) (err error) {
	var tx *sqlx.Tx = db.MustBeginTx(ctx, datastore.TxOptions)
	var teams []*Team = []*Team{}

	// Loop over each team and sort out the joins
	for _, team := range self.TeamList {
		var join *Join
		team.ID = 0
		// Find / create team DB entry
		err = team.UpdateSelf(ctx, db, tx)
		if err != nil {
			return
		}
		// Find / create the join between both
		join = &Join{ReleaseID: self.ID, TeamID: team.ID}
		err = join.UpdateSelf(ctx, db, tx)
		if err != nil {
			return
		}
		teams = append(teams, team)
	}
	self.TeamList = teams
	// commit
	err = tx.Commit()
	// roll back and warn of error
	if err != nil {
		err = tx.Rollback()
		slog.Error("[releases.InsertJoins]", slog.String("err", err.Error()))
	}

	return
}

// Teams fetches the teams from the database
func (self *Release) Teams(ctx context.Context, db *sqlx.DB) (teams []*Team, err error) {
	teams, err = datastore.SelectMany[*Team](ctx, db, releasesdb.GetTeamsForRelease, self)
	return
}

// New
// Record interface
func (self *Release) New() record.Record {
	return &Release{}
}

// UID
// Record interface
func (self *Release) UID() string {
	return fmt.Sprintf("%s-%d", "releases", self.ID)
}

// SetID
// Record interface
func (self *Release) SetID(id int) {
	self.ID = id
}

// TDate
// Transformable interface
func (self *Release) TDate() string {
	return self.Date
}

// TValue
// Transformable interface
// Always 1 so it increments per period
func (self *Release) TValue() string {
	return strconv.Itoa(self.Count)
}

// Team captures the github team from the repository
type Team struct {
	ID   int    `json:"team_id,omitempty" db:"team_id" faker:"unique, boundary_start=1, boundary_end=2000000" doc:"Database primary key."`
	Name string `json:"team_name,omitempty" db:"team_name" faker:"oneof: A, B, C, D"`
}

// UpdateSelf finds its record within the database (by checking name) and updates its ID to match
// the result.
// If a record is not found, it creates a new entry and then updates its own ID with the
// results from InsertOne
func (self *Team) UpdateSelf(ctx context.Context, db *sqlx.DB, tx *sqlx.Tx) (err error) {
	var t *Team = &Team{}
	t, err = datastore.SelectOne[*Team](ctx, db, releasesdb.GetTeamByName, self)

	if t != nil {
		self.ID = t.ID
	} else {
		self.ID, err = datastore.InsertOne(ctx, db, releasesdb.InsertTeam, self, tx)
	}

	return
}

// New
// Record interface
func (self *Team) New() record.Record {
	return &Team{}
}

// UID
// Record interface
func (self *Team) UID() string {
	return fmt.Sprintf("%s-%d", "team", self.ID)
}

// SetID
// Record interface
func (self *Team) SetID(id int) {
	self.ID = id
}

// Join is many to many table between both team and releases
type Join struct {
	ID        int `json:"join_id,omitempty" db:"join_id" faker:"unique, boundary_start=1, boundary_end=2000000" doc:"Database primary key."`
	TeamID    int `json:"team_id,omitempty" db:"team_id"`
	ReleaseID int `json:"release_id,omitempty" db:"release_id"`
}

// UpdateSelf finds its record within the database (by checking ids) and updates its ID to match
// the result.
// If a record is not found, it creates a new entry and then updates its own ID with the
// results from InsertOne
func (self *Join) UpdateSelf(ctx context.Context, db *sqlx.DB, tx *sqlx.Tx) (err error) {
	var j *Join = &Join{}
	j, err = datastore.SelectOne[*Join](ctx, db, releasesdb.GetJoin, self)

	if j != nil {
		self.ID = j.ID
	} else {
		self.ID, err = datastore.InsertOne(ctx, db, releasesdb.InsertJoin, self, tx)
	}

	return
}

// New
// Record interface
func (self *Join) New() record.Record {
	return &Join{}
}

// UID
// Record interface
func (self *Join) UID() string {
	return fmt.Sprintf("%s-%d", "join", self.ID)
}

// SetID
// Record interface
func (self *Join) SetID(id int) {
	self.ID = id
}
