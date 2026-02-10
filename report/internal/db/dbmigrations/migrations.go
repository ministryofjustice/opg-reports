package dbmigrations

// list of ALL migrations to run in sequence
var MIGRATIONS []*Migration = []*Migration{
	{Name: "create_migration", SQL: create_migration},
	{Name: "create_teams", SQL: create_teams},
	{Name: "create_aws_accounts", SQL: create_aws_accounts},
	{Name: "create_aws_costs", SQL: create_aws_costs},
	{Name: "create_aws_uptime", SQL: create_aws_uptime},
	{Name: "create_codebases", SQL: create_codebases},
	{Name: "create_codeowners", SQL: create_codeowners},

	{Name: "agnostic_accounts", SQL: agnostic_accounts},
	{Name: "agnostic_costs", SQL: agnostic_costs},
	{Name: "agnostic_uptime", SQL: agnostic_uptime},

	{Name: "lowercase_team_name", SQL: lowercase_team_name},
	{Name: "copy_from_github_codewners", SQL: copy_from_github_codewners},
}

// copy from codeowners into codebases and codeowners
const copy_from_github_codewners string = `
INSERT INTO codebases (created_at, full_name, name, url)
	SELECT created_at, repository, replace(repository, "ministryofjustice/", ""), concat("https://github.com/", repository) FROM github_codeownership GROUP BY repository;
INSERT INTO codeowners (created_at, name, codebase_full_name, team_name)
	SELECT created_at, codeowner, repository, team FROM github_codeownership;
UPDATE codeowners SET(team_name) = LOWER(team_name);
`

// drop things to lower cases
const lowercase_team_name string = `
UPDATE codeowners SET(team_name) = LOWER(team_name);
UPDATE accounts SET(team_name) = LOWER(team_name);
UPDATE teams SET(name) = LOWER(name);
`

// agnostic_uptime removes the aws prefix
const agnostic_uptime string = `
CREATE TABLE IF NOT EXISTS uptime (
	id INTEGER PRIMARY KEY,
	created_at TEXT NOT NULL DEFAULT (strftime('%FT%TZ', 'now') ),
	vendor TEXT NOT NULL DEFAULT 'aws',
	date TEXT NOT NULL,
	account_id TEXT,
	average TEXT NOT NULL,
	granularity TEXT NOT NULL,
	UNIQUE (account_id,date)
) STRICT;

INSERT INTO uptime (created_at, date, account_id, average, granularity)
	SELECT created_at, date, aws_account_id, average, granularity FROM aws_uptime;

DROP INDEX aws_uptime_date_idx;
DROP INDEX aws_uptime_account_date_idx;

CREATE INDEX IF NOT EXISTS idx_uptime_date ON uptime(date);
CREATE INDEX IF NOT EXISTS idx_uptime_account_date ON uptime(account_id,date);

DROP TABLE aws_uptime;
`

// agnostic_costs removes the aws prefix
const agnostic_costs string = `
CREATE TABLE IF NOT EXISTS infracosts (
	id INTEGER PRIMARY KEY,
	created_at TEXT NOT NULL DEFAULT (strftime('%FT%TZ', 'now') ),
	vendor TEXT NOT NULL DEFAULT 'aws',
	region TEXT DEFAULT "NoRegion" NOT NULL,
	service TEXT NOT NULL,
	date TEXT NOT NULL,
	cost TEXT NOT NULL,
	account_id TEXT,
	UNIQUE (account_id,date,region,service)
) STRICT;

INSERT INTO infracosts (created_at, region, service, date, cost, account_id)
	SELECT created_at, region, service, date, cost, aws_account_id FROM aws_costs;

DROP INDEX aws_costs_date_idx;
DROP INDEX aws_costs_date_account_idx;
DROP INDEX aws_costs_unique_idx;

CREATE INDEX IF NOT EXISTS idx_infracosts_date ON infracosts(date);
CREATE INDEX IF NOT EXISTS idx_infracosts_date_account ON infracosts(date, account_id);
CREATE INDEX IF NOT EXISTS idx_infracosts_unique ON infracosts(account_id,date,region,service);

DROP TABLE aws_costs;
`

// rename aws_accounts to accounts
// add vendor column
const agnostic_accounts string = `
DROP INDEX aws_accounts_id_idx;
ALTER TABLE aws_accounts ADD vendor TEXT NOT NULL DEFAULT 'aws';
ALTER TABLE aws_accounts RENAME TO accounts;
CREATE INDEX IF NOT EXISTS idx_accounts_id ON accounts(id);
`
const create_codeowners string = `
CREATE TABLE IF NOT EXISTS codeowners (
	id INTEGER PRIMARY KEY,
	created_at TEXT NOT NULL DEFAULT (strftime('%FT%TZ', 'now') ),
	name TEXT NOT NULL,
	codebase_full_name TEXT NOT NULL,
	team_name TEXT NOT NULL,
	UNIQUE (name,codebase_full_name,team_name)
) STRICT;

CREATE INDEX IF NOT EXISTS idx_codeowners_join ON codeowners(codebase_full_name,team_name);
`

// create the codebase table
const create_codebases string = `
CREATE TABLE IF NOT EXISTS codebases (
	id INTEGER PRIMARY KEY,
	created_at TEXT NOT NULL DEFAULT (strftime('%FT%TZ', 'now') ),
	name TEXT NOT NULL,
	full_name TEXT NOT NULL,
	url TEXT NOT NULL,
	UNIQUE (full_name)
) STRICT;
CREATE INDEX IF NOT EXISTS idx_codebases_fullname ON codebases(full_name);
`

// create the aws uptime tracking table
const create_aws_uptime string = `
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

// create the aws_costs table
const create_aws_costs string = `
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
`

// create the aws_accounts table & indexes
const create_aws_accounts string = `
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
`

// creates the main team table
const create_teams string = `
CREATE TABLE IF NOT EXISTS teams (
	name TEXT PRIMARY KEY,
	created_at TEXT NOT NULL DEFAULT (strftime('%FT%TZ', 'now') )
) STRICT;
`

// create the main migration table for tracking
const create_migration string = `
CREATE TABLE IF NOT EXISTS migrations (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL
) STRICT;
`
