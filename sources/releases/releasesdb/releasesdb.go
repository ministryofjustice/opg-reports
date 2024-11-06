// Package releasesdb contains all statments for releases api etc
package releasesdb

import "github.com/ministryofjustice/opg-reports/pkg/datastore"

// Create tables
const (
	CreateTeamTable            datastore.CreateStatement = `CREATE TABLE IF NOT EXISTS teams (id INTEGER PRIMARY KEY,name TEXT NOT NULL) STRICT;`
	CreateTeamReleaseJoinTable datastore.CreateStatement = `CREATE TABLE IF NOT EXISTS releases_teams (id INTEGER PRIMARY KEY,release_id INTEGER NOT NULL,team_id INTEGER NOT NULL) STRICT`
	CreateReleaseTable         datastore.CreateStatement = `CREATE TABLE IF NOT EXISTS releases (id INTEGER PRIMARY KEY,ts TEXT NOT NULL,repository TEXT NOT NULL,name TEXT NOT NULL,source TEXT NOT NULL,type TEXT NOT NULL,date TEXT NOT NULL,count INTEGER) STRICT;`
)

// Create Indexes
const (
	CreateReleaseDateIndex datastore.CreateStatement = `CREATE INDEX IF NOT EXISTS release_date_idx ON releases(date);`
	CreateTeamNameIndex    datastore.CreateStatement = `CREATE INDEX IF NOT EXISTS team_name_idx ON teams(name);`
)

// Inserts for tables
const (
	InsertTeam    datastore.InsertStatement = `INSERT INTO teams (name) VALUES (:name) RETURNING id;`
	InsertJoin    datastore.InsertStatement = `INSERT INTO releases_teams (release_id, team_id) VALUES (:release_id, :team_id) RETURNING id;`
	InsertRelease datastore.InsertStatement = `INSERT INTO releases (ts,repository,name,source,type,date,count) VALUES (:ts,:repository,:name,:source,:type,:date,:count) RETURNING id;`
)

// Counters
const (
	TeamsCount   datastore.SelectStatement = `SELECT count(*) as row_count FROM teams LIMIT 1;`
	JoinCount    datastore.SelectStatement = `SELECT count(*) as row_count FROM releases_teams LIMIT 1;`
	ReleaseCount datastore.SelectStatement = `SELECT count(*) as row_count FROM releases LIMIT 1;`
)

const GetTeamID datastore.SelectStatement = `
SELECT
	id
FROM teams
WHERE
	name = ?
LIMIT 1
`

const AllReleases datastore.SelectStatement = `
SELECT
	*
FROM releases
ORDER BY id ASC
`

const PerInterval datastore.NamedSelectStatement = `
SELECT
    coalesce(SUM(count), 0) as count,
    strftime(:date_format, date) as date
FROM costs
WHERE
    date >= :start_date
    AND date < :end_date
GROUP BY strftime(:date_format, date)
ORDER by strftime(:date_format, date) ASC
`
