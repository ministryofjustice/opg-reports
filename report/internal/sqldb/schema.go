package sqldb

// SCHEMA contains all of the database tables and indexes
const SCHEMA string = `
CREATE TABLE IF NOT EXISTS teams (
	id INTEGER PRIMARY KEY,
	created_at TEXT,
	name TEXT NOT NULL UNIQUE
) STRICT;

CREATE INDEX IF NOT EXISTS team_name ON teams(name);

CREATE TABLE IF NOT EXISTS aws_accounts (
	id TEXT PRIMARY KEY,
	created_at TEXT,
	name TEXT NOT NULL,
	label TEXT NOT NULL,
	environment TEXT DEFAULT "production" NOT NULL,
	team_id INTEGER
) WITHOUT ROWID;

CREATE INDEX IF NOT EXISTS aws_accounts_id ON aws_accounts(id);
`
