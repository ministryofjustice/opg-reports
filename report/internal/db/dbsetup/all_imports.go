package dbsetup

var _IMPORTS map[string]*ImportStatements = map[string]*ImportStatements{
	"accounts": {
		Insert: account_insert,
	},
	"codebases": {
		Pre:    codebase_truncate,
		Insert: codebase_insert,
	},
	"codeowners": {
		Pre:    codeowner_truncate,
		Insert: codeowner_insert,
	},
	"infracosts": {
		Insert: infracost_insert,
	},
	"teams": {
		Insert: team_insert,
	},
	"uptime": {
		Insert: uptime_insert,
	},
}

const uptime_insert string = `
INSERT INTO uptime (
	date,
	average,
	granularity,
	account_id
) VALUES (
	:date,
	:average,
	:granularity,
	:account_id
) ON CONFLICT (account_id,date)
 	DO UPDATE SET average=excluded.average, granularity=excluded.granularity
RETURNING id
;
`

const team_insert string = `
INSERT INTO teams (
	name
) VALUES (
	:name
) ON CONFLICT (name)
 	DO UPDATE SET name=excluded.name
RETURNING name
;
`

const infracost_insert string = `
INSERT INTO infracosts (
	region,
	service,
	date,
	cost,
	account_id
) VALUES (
	:region,
	:service,
	:date,
	:cost,
	:account_id
) ON CONFLICT (account_id,date,region,service)
 	DO UPDATE SET cost=excluded.cost
RETURNING id
;
`

const codeowner_truncate string = `DELETE FROM codeowners;`
const codeowner_insert string = `
INSERT INTO codeowners (
	name,
	codebase_full_name,
	team_name
) VALUES (
 	:name,
	:codebase_full_name,
	:team_name
)
ON CONFLICT (name,codebase_full_name,team_name)
 	DO UPDATE SET
		team_name=excluded.team_name,
		name=excluded.name,
		team_name=excluded.team_name
RETURNING id
;
`

const codebase_truncate string = `DELETE FROM codebases;`
const codebase_insert string = `
INSERT INTO codebases (
	name,
	full_name,
	url
) VALUES (
	:name,
	:full_name,
	:url
)
ON CONFLICT (full_name)
 	DO UPDATE SET
		name=excluded.name,
		url=excluded.url
RETURNING id
;
`

const account_insert string = `
INSERT INTO accounts (
	id,
	name,
	label,
	environment,
	team_name
) VALUES (
	:id,
	:name,
	:label,
	:environment,
	:team_name
)
ON CONFLICT (id)
 	DO UPDATE SET
		name=excluded.name,
		label=excluded.label,
		environment=excluded.environment
RETURNING id
;
`
