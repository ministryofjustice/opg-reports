// Package releasesdb contains all statments for releases api etc
package releasesdb

import "github.com/ministryofjustice/opg-reports/pkg/datastore"

// Create tables
const (
	CreateTeamTable            datastore.CreateStatement = `CREATE TABLE IF NOT EXISTS teams (team_id INTEGER PRIMARY KEY, team_name TEXT NOT NULL) STRICT;`
	CreateTeamReleaseJoinTable datastore.CreateStatement = `CREATE TABLE IF NOT EXISTS releases_teams (join_id INTEGER PRIMARY KEY,release_id INTEGER NOT NULL,team_id INTEGER NOT NULL) STRICT`
	CreateReleaseTable         datastore.CreateStatement = `CREATE TABLE IF NOT EXISTS releases (id INTEGER PRIMARY KEY,ts TEXT NOT NULL,repository TEXT NOT NULL,name TEXT NOT NULL,source TEXT NOT NULL,type TEXT NOT NULL,date TEXT NOT NULL,count INTEGER) STRICT;`
)

// Create Indexes
const (
	CreateReleaseDateIndex datastore.CreateStatement = `CREATE INDEX IF NOT EXISTS release_date_idx ON releases(date);`
	CreateTeamNameIndex    datastore.CreateStatement = `CREATE INDEX IF NOT EXISTS team_name_idx ON teams(team_name);`
)

// Inserts for tables
const (
	InsertTeam    datastore.InsertStatement = `INSERT INTO teams (team_name) VALUES (:team_name) RETURNING team_id;`
	InsertJoin    datastore.InsertStatement = `INSERT INTO releases_teams (release_id, team_id) VALUES (:release_id, :team_id) RETURNING join_id;`
	InsertRelease datastore.InsertStatement = `INSERT INTO releases (ts, repository, name, source, type, date,count) VALUES (:ts,:repository,:name,:source,:type,:date,:count) RETURNING id;`
)

// Counters
const (
	TeamsCount   datastore.SelectStatement = `SELECT count(*) as row_count FROM teams LIMIT 1;`
	JoinCount    datastore.SelectStatement = `SELECT count(*) as row_count FROM releases_teams LIMIT 1;`
	ReleaseCount datastore.SelectStatement = `SELECT count(*) as row_count FROM releases LIMIT 1;`
)

// Team selects
const (
	GetTeamByName datastore.NamedSelectStatement = `SELECT team_id, team_name FROM teams WHERE team_name = :team_name LIMIT 1`
	GetTeamByID   datastore.NamedSelectStatement = `SELECT team_id, team_name FROM teams WHERE team_id = :team_id LIMIT 1`
)

// Join Selects
const (
	GetJoin datastore.NamedSelectStatement = `SELECT join_id FROM releases_teams WHERE release_id = :release_id AND team_id = :team_id LIMIT 1`
)

// Release selects
const (
	GetRandomRelease   datastore.SelectStatement      = `SELECT * FROM releases ORDER BY RANDOM() LIMIT 1`
	ListReleases       datastore.SelectStatement      = `SELECT * FROM releases ORDER BY id ASC`
	GetTeamsForRelease datastore.NamedSelectStatement = `SELECT teams.team_id as team_id, teams.team_name as team_name FROM releases_teams LEFT JOIN teams ON releases_teams.team_id = teams.team_id WHERE releases_teams.release_id = :id`
)

// const PerInterval datastore.NamedSelectStatement = `
// SELECT
//     coalesce(SUM(count), 0) as count,
//     strftime(:date_format, date) as date
// FROM costs
// WHERE
//     date >= :start_date
//     AND date < :end_date
// GROUP BY strftime(:date_format, date)
// ORDER by strftime(:date_format, date) ASC
// `
