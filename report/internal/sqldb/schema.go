package sqldb

// SCHEMA contains all of the database tables and indexes
const SCHEMA string = `
CREATE TABLE IF NOT EXISTS teams (
	id INTEGER PRIMARY KEY,
	created_at TEXT,
	name TEXT NOT NULL UNIQUE
) STRICT;

CREATE INDEX IF NOT EXISTS team_name ON teams(name)
`
