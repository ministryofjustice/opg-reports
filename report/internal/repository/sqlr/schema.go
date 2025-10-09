package sqlr

// SCHEMA contains all of the database tables and indexes
const SCHEMA string = `
CREATE TABLE IF NOT EXISTS teams (
	name TEXT PRIMARY KEY,
	created_at TEXT NOT NULL DEFAULT (strftime('%FT%TZ', 'now') )
) STRICT;

CREATE TABLE IF NOT EXISTS aws_accounts (
	id TEXT PRIMARY KEY,
	created_at TEXT NOT NULL DEFAULT (strftime('%FT%TZ', 'now') ),
	name TEXT NOT NULL,
	label TEXT NOT NULL,
	environment TEXT NOT NULL DEFAULT "production",
	uptime_tracking TEXT NOT NULL DEFAULT "false",
	team_name TEXT NOT NULL DEFAULT "ORG"
) WITHOUT ROWID;

CREATE INDEX IF NOT EXISTS aws_accounts_id_idx ON aws_accounts(id);

CREATE TABLE IF NOT EXISTS aws_costs (
	id INTEGER PRIMARY KEY,
	created_at TEXT NOT NULL DEFAULT (strftime('%FT%TZ', 'now') ),
	region TEXT DEFAULT "NoRegion" NOT NULL,
	service TEXT NOT NULL,
	date TEXT NOT NULL,
	cost TEXT NOT NULL,
	aws_account_id TEXT,
	UNIQUE (aws_account_id,date,region,service)
) STRICT;

CREATE INDEX IF NOT EXISTS aws_costs_date_idx ON aws_costs(date);
CREATE INDEX IF NOT EXISTS aws_costs_date_account_idx ON aws_costs(date, aws_account_id);
CREATE INDEX IF NOT EXISTS aws_costs_unique_idx ON aws_costs(aws_account_id,date,region,service);

CREATE TABLE IF NOT EXISTS aws_uptime (
	id INTEGER PRIMARY KEY,
	created_at TEXT NOT NULL DEFAULT (strftime('%FT%TZ', 'now') ),
	date TEXT NOT NULL,
	aws_account_id TEXT,
	average TEXT NOT NULL,
	granularity TEXT NOT NULL,
	UNIQUE (aws_account_id,date)
) STRICT;

CREATE INDEX IF NOT EXISTS aws_uptime_date_idx ON aws_uptime(date);
CREATE INDEX IF NOT EXISTS aws_uptime_account_date_idx ON aws_uptime(aws_account_id,date);
`
