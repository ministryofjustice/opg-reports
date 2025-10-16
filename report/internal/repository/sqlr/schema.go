package sqlr

const migration_teams_table string = `
CREATE TABLE IF NOT EXISTS teams (
	name TEXT PRIMARY KEY,
	created_at TEXT NOT NULL DEFAULT (strftime('%FT%TZ', 'now') )
) STRICT;`

const migration_aws_accounts_table string = `
CREATE TABLE IF NOT EXISTS aws_accounts (
	id TEXT PRIMARY KEY,
	created_at TEXT NOT NULL DEFAULT (strftime('%FT%TZ', 'now') ),
	name TEXT NOT NULL,
	label TEXT NOT NULL,
	environment TEXT NOT NULL DEFAULT "production",
	uptime_tracking TEXT NOT NULL DEFAULT "false",
	team_name TEXT NOT NULL DEFAULT "ORG"
) WITHOUT ROWID;`
const migration_aws_account_idx string = `CREATE INDEX IF NOT EXISTS aws_accounts_id_idx ON aws_accounts(id);`

const migration_aws_costs_table string = `
CREATE TABLE IF NOT EXISTS aws_costs (
	id INTEGER PRIMARY KEY,
	created_at TEXT NOT NULL DEFAULT (strftime('%FT%TZ', 'now') ),
	region TEXT DEFAULT "NoRegion" NOT NULL,
	service TEXT NOT NULL,
	date TEXT NOT NULL,
	cost TEXT NOT NULL,
	aws_account_id TEXT,
	UNIQUE (aws_account_id,date,region,service)
) STRICT;`
const migration_aws_costs_idx_date string = `CREATE INDEX IF NOT EXISTS aws_costs_date_idx ON aws_costs(date);`
const migration_aws_costs_idx_date_account string = `CREATE INDEX IF NOT EXISTS aws_costs_date_account_idx ON aws_costs(date, aws_account_id);`
const migration_aws_costs_idx_merged string = `CREATE INDEX IF NOT EXISTS aws_costs_unique_idx ON aws_costs(aws_account_id,date,region,service);`

const migration_aws_uptime_table string = `
CREATE TABLE IF NOT EXISTS aws_uptime (
	id INTEGER PRIMARY KEY,
	created_at TEXT NOT NULL DEFAULT (strftime('%FT%TZ', 'now') ),
	date TEXT NOT NULL,
	aws_account_id TEXT,
	average TEXT NOT NULL,
	granularity TEXT NOT NULL,
	UNIQUE (aws_account_id,date)
) STRICT;
`
const migration_aws_uptime_idx_date string = `CREATE INDEX IF NOT EXISTS aws_uptime_date_idx ON aws_uptime(date);`
const migration_aws_uptime_idx_date_account string = `CREATE INDEX IF NOT EXISTS aws_uptime_account_date_idx ON aws_uptime(aws_account_id,date);`

var DB_MIGRATIONS_UP []string = []string{
	migration_teams_table,
	migration_aws_accounts_table,
	migration_aws_account_idx,
	migration_aws_costs_table,
	migration_aws_costs_idx_date,
	migration_aws_costs_idx_date_account,
	migration_aws_costs_idx_merged,
	migration_aws_uptime_table,
	migration_aws_uptime_idx_date,
	migration_aws_uptime_idx_date_account,
}
