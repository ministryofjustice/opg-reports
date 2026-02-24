package migrations

const create_teams string = `
CREATE TABLE IF NOT EXISTS teams (
	name TEXT PRIMARY KEY,
	created_at TEXT NOT NULL DEFAULT (strftime('%FT%TZ', 'now') )
) STRICT
;
`

const migrate_costs string = `
INSERT INTO costs (created_at, region, service, month, cost, account_id)
	SELECT created_at, region, service, strftime("%Y-%m",date) as date, cost, aws_account_id FROM aws_costs;

DROP INDEX aws_costs_date_idx;
DROP INDEX aws_costs_date_account_idx;
DROP INDEX aws_costs_unique_idx;
DROP TABLE aws_costs;
`

const create_agnostic_costs string = `
CREATE TABLE IF NOT EXISTS costs (
	id INTEGER PRIMARY KEY,
	created_at TEXT NOT NULL DEFAULT (strftime('%FT%TZ', 'now') ),
	vendor TEXT NOT NULL DEFAULT 'aws',
	region TEXT DEFAULT "NoRegion" NOT NULL,
	service TEXT NOT NULL,
	month TEXT NOT NULL,
	cost TEXT NOT NULL,
	account_id TEXT,
	UNIQUE (account_id,month,region,service)
) STRICT;

CREATE INDEX IF NOT EXISTS idx_costs_date ON costs(month);
CREATE INDEX IF NOT EXISTS idx_costs_date_account ON costs(month, account_id);
CREATE INDEX IF NOT EXISTS idx_costs_unique ON costs(account_id, month, region, service);
`

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

const agnostic_accounts string = `
DROP INDEX aws_accounts_id_idx;
ALTER TABLE aws_accounts ADD vendor TEXT NOT NULL DEFAULT 'aws';
ALTER TABLE aws_accounts RENAME TO accounts;
CREATE INDEX IF NOT EXISTS idx_accounts_id ON accounts(id);
`

const lowercase_team_name string = `
UPDATE accounts SET(team_name) = LOWER(team_name);
UPDATE teams SET(name) = LOWER(name);
`

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

// agnostic_uptime removes the aws prefix
const agnostic_uptime string = `
CREATE TABLE IF NOT EXISTS uptime (
	id INTEGER PRIMARY KEY,
	created_at TEXT NOT NULL DEFAULT (strftime('%FT%TZ', 'now') ),
	vendor TEXT NOT NULL DEFAULT 'aws',
	month TEXT NOT NULL,
	account_id TEXT,
	average TEXT NOT NULL,
	granularity TEXT NOT NULL,
	UNIQUE (account_id,month)
) STRICT;

INSERT INTO uptime (month, account_id, average, granularity)
	SELECT strftime("%Y-%m", date) as date, aws_account_id, AVG(average), granularity FROM aws_uptime GROUP BY strftime("%Y-%m", date), aws_account_id;

DROP INDEX IF EXISTS aws_uptime_date_idx;
DROP INDEX IF EXISTS aws_uptime_account_date_idx;
CREATE INDEX IF NOT EXISTS idx_uptime_month ON uptime(month);
CREATE INDEX IF NOT EXISTS idx_uptime_account_month ON uptime(account_id,month);
DROP TABLE IF EXISTS aws_uptime;
`

const create_codebases string = `
CREATE TABLE IF NOT EXISTS codebases (
	id INTEGER PRIMARY KEY,
	created_at TEXT NOT NULL DEFAULT (strftime('%FT%TZ', 'now') ),
	name TEXT NOT NULL,
	full_name TEXT NOT NULL,
	url TEXT NOT NULL,
	compliance_level TEXT,
	compliance_report_url TEXT,
	compliance_badge TEXT,
	UNIQUE (full_name)
) STRICT;
CREATE INDEX IF NOT EXISTS idx_codeowners ON codebases(full_name);
`
