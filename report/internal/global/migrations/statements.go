package migrations

const lower_case_teams string = `
UPDATE teams SET(name) = LOWER(name)
;
`

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
