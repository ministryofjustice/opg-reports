package awsaccount

// list of sql insert statements of various types
const stmtDropTable string = `DROP TABLE IF EXISTS aws_accounts`

const stmtDeleteAll string = `DELETE FROM aws_accounts`

const stmtSelectAll string = `
SELECT
	aws_accounts.id,
	aws_accounts.name,
	aws_accounts.label,
	aws_accounts.environment,
	json_object(
		'name', aws_accounts.team_name
	) as team
FROM aws_accounts
GROUP BY aws_accounts.id
ORDER BY aws_accounts.team_name ASC, aws_accounts.name ASC, aws_accounts.environment ASC;`

const stmtUpdateEmptyEnvironments string = `
UPDATE aws_accounts
SET environment = "production"
WHERE environment = ""
`

const stmtImport string = `
INSERT INTO aws_accounts (
	id,
	name,
	label,
	environment,
	team_name
) SELECT
	:id,
	:name,
	:label,
	:environment,
	teams.name
FROM teams WHERE name=:team_name LIMIT 1
ON CONFLICT (id)
 	DO UPDATE SET
		name=excluded.name,
		label=excluded.label,
		environment=excluded.environment
RETURNING id;`
