package migrations

const run_vacuum string = `VACUUM;`

const lowercase_team_name string = `
UPDATE accounts SET(team_name) = LOWER(team_name);
UPDATE teams SET(name) = LOWER(name);
`

const create_teams string = `
CREATE TABLE IF NOT EXISTS teams (
	name TEXT PRIMARY KEY,
	created_at TEXT NOT NULL DEFAULT (strftime('%FT%TZ', 'now') )
) STRICT
;
`
const create_accounts string = `
CREATE TABLE IF NOT EXISTS accounts (
	id TEXT PRIMARY KEY,
	created_at TEXT NOT NULL DEFAULT (strftime('%FT%TZ', 'now') ),
	vendor TEXT NOT NULL DEFAULT 'aws',
	name TEXT NOT NULL,
	label TEXT NOT NULL,
	environment TEXT NOT NULL DEFAULT "production",
	uptime_tracking TEXT NOT NULL DEFAULT "false",
	team_name TEXT NOT NULL DEFAULT "ORG"
) WITHOUT ROWID;
CREATE INDEX IF NOT EXISTS accounts_idx ON accounts(id);
`

const create_costs string = `
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

// agnostic_uptime removes the aws prefix
const create_uptime string = `
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
CREATE INDEX IF NOT EXISTS idx_uptime_month ON uptime(month);
CREATE INDEX IF NOT EXISTS idx_uptime_account_month ON uptime(account_id,month);
`

const create_codebases string = `
CREATE TABLE IF NOT EXISTS codebases (
	id INTEGER PRIMARY KEY,
	created_at TEXT NOT NULL DEFAULT (strftime('%FT%TZ', 'now') ),
	name TEXT NOT NULL,
	full_name TEXT NOT NULL,
	url TEXT NOT NULL,
	archived INTEGER NOT NULL DEFAULT 0,
	UNIQUE (full_name)
) STRICT;
CREATE INDEX IF NOT EXISTS idx_codebases ON codebases(full_name);
`

const create_codebase_stats string = `
CREATE TABLE IF NOT EXISTS codebase_stats (
	id INTEGER PRIMARY KEY,
	codebase TEXT NOT NULL,
	visibility TEXT NOT NULL,

	compliance_level TEXT,
	compliance_report_url TEXT,
	compliance_badge TEXT,
	compliance_grade INTEGER,

	trivy_usage INTEGER,
	trivy_sbom_usage INTEGER,
	trivy_locations TEXT,
	UNIQUE (codebase)
) STRICT;
CREATE INDEX IF NOT EXISTS idx_codebase_stats ON codebase_stats(codebase);
`

const create_codeowner string = `
CREATE TABLE IF NOT EXISTS codebase_owners (
	id INTEGER PRIMARY KEY,
	created_at TEXT NOT NULL DEFAULT (strftime('%FT%TZ', 'now') ),
	owner TEXT NOT NULL,
	codebase TEXT NOT NULL,
	team_name TEXT NOT NULL,
	UNIQUE (owner,codebase,team_name)
) STRICT;
CREATE INDEX IF NOT EXISTS idx_codeowners ON codebase_owners(codebase,team_name);
`

const create_codebase_metrics string = `
CREATE TABLE IF NOT EXISTS codebase_metrics (
	id INTEGER PRIMARY KEY,
	codebase TEXT NOT NULL,
	month TEXT NOT NULL,
	releases INTEGER NOT NULL,
	releases_securityish INTEGER DEFAULT 0,
	releases_average_time TEXT DEFAULT "0.0",
	pr_count INTEGER DEFAULT 0,
	pr_count_securityish INTEGER DEFAULT 0,
	pr_count_stale  INTEGER DEFAULT 0,
	pr_average_time TEXT DEFAULT "0.0",
	UNIQUE (codebase,month)
) STRICT;
CREATE INDEX IF NOT EXISTS idx_codebase_metrics_month ON codebase_metrics(codebase,month);
CREATE INDEX IF NOT EXISTS idx_codebase_metrics ON codebase_metrics(codebase);
`
