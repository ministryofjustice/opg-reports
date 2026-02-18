package teammigrations

import (
	"opg-reports/report/internal/global/migrations"
)

// Migrations is a slice to ensure ordering is kept
var Migrations = []*migrations.Migration{
	{Key: "create_teams", Stmt: create_teams},
	{Key: "lower_case_teams", Stmt: lower_case_teams},
}

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
