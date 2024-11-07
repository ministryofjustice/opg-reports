// Package releasesdb contains all statments for releases api etc
//   - Create* is a creation
//   - Insert* is an insert
//   - Get* is a select returning a single item
//   - List* is a select returning many
package releasesdb

import "github.com/ministryofjustice/opg-reports/pkg/datastore"

// Create tables
const (
	CreateTeamTable            datastore.CreateStatement = `CREATE TABLE IF NOT EXISTS teams (id INTEGER PRIMARY KEY, name TEXT NOT NULL) STRICT;`
	CreateTeamReleaseJoinTable datastore.CreateStatement = `CREATE TABLE IF NOT EXISTS releases_teams (id INTEGER PRIMARY KEY, release_id INTEGER NOT NULL,team_id INTEGER NOT NULL) STRICT`
	CreateReleaseTable         datastore.CreateStatement = `CREATE TABLE IF NOT EXISTS releases (id INTEGER PRIMARY KEY, ts TEXT NOT NULL,repository TEXT NOT NULL,name TEXT NOT NULL,source TEXT NOT NULL,type TEXT NOT NULL,date TEXT NOT NULL,count INTEGER) STRICT;`
)

// Create Indexes for common queries
const (
	CreateReleaseDateIndex datastore.CreateStatement = `CREATE INDEX IF NOT EXISTS release_date_idx ON releases(date);`
	CreateTeamNameIndex    datastore.CreateStatement = `CREATE INDEX IF NOT EXISTS team_name_idx ON teams(name);`
	CreateJoinTeamIndex    datastore.CreateStatement = `CREATE INDEX IF NOT EXISTS join_team_idx ON releases_teams(team_id);`
	CreateJoinReleaseIndex datastore.CreateStatement = `CREATE INDEX IF NOT EXISTS join_release_idx ON releases_teams(release_id);`
)

// Inserts for tables
const (
	InsertTeam    datastore.InsertStatement = `INSERT INTO teams (name) VALUES (:name) RETURNING id;`
	InsertJoin    datastore.InsertStatement = `INSERT INTO releases_teams (release_id, team_id) VALUES (:release_id, :team_id) RETURNING id;`
	InsertRelease datastore.InsertStatement = `INSERT INTO releases (ts, repository, name, source, type, date, count) VALUES (:ts,:repository,:name,:source,:type,:date,:count) RETURNING id;`
)

// Counters
const (
	CountTeams    datastore.SelectStatement = `SELECT count(*) as row_count FROM teams LIMIT 1;`
	CountJoins    datastore.SelectStatement = `SELECT count(*) as row_count FROM releases_teams LIMIT 1;`
	CountReleases datastore.SelectStatement = `SELECT count(*) as row_count FROM releases LIMIT 1;`
)

// Team selects
const (
	GetTeamByName datastore.NamedSelectStatement = `SELECT id, name FROM teams WHERE name = :name LIMIT 1`
	GetTeamByID   datastore.NamedSelectStatement = `SELECT id, name FROM teams WHERE id = :id LIMIT 1`
)

// Join Selects
const (
	GetJoin datastore.NamedSelectStatement = `SELECT id FROM releases_teams WHERE release_id = :release_id AND team_id = :team_id LIMIT 1`
)

// Release selects
const (
	GetRandomRelease   datastore.SelectStatement      = `SELECT * FROM releases ORDER BY RANDOM() LIMIT 1`
	ListReleases       datastore.SelectStatement      = `SELECT * FROM releases ORDER BY id ASC`
	GetTeamsForRelease datastore.NamedSelectStatement = `SELECT teams.id as id, teams.name as name FROM releases_teams LEFT JOIN teams ON releases_teams.team_id = teams.id WHERE releases_teams.release_id = :id`
)

const ListReleasesGroupedByInterval datastore.NamedSelectStatement = `
SELECT
	strftime(:date_format, releases.date) as date,
	COUNT(releases.id) as count
FROM releases
GROUP BY strftime(:date_format, releases.date)
ORDER BY strftime(:date_format, releases.date) ASC
;`

const ListReleasesGroupedByIntervalAndTeam datastore.NamedSelectStatement = `
SELECT
	teams.name as team_name,
	strftime(:date_format, releases.date) as date,
	COUNT(releases.id) as count
FROM releases
LEFT JOIN releases_teams on releases_teams.release_id = releases.id
LEFT JOIN teams on teams.id = releases_teams.team_id
GROUP BY strftime(:date_format, releases.date), releases_teams.team_id
ORDER BY strftime(:date_format, releases.date), teams.name ASC
;`
