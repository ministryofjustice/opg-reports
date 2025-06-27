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
		'name', teams.name
	) as team
FROM aws_accounts
LEFT JOIN teams on teams.id = aws_accounts.team_id
GROUP BY aws_accounts.id
ORDER BY teams.name ASC, aws_accounts.name ASC, aws_accounts.environment ASC;`

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
	team_id
) SELECT
	:id,
	:name,
	:label,
	:environment,
	id
FROM teams WHERE name=:team_name LIMIT 1
ON CONFLICT (id)
 	DO UPDATE SET
		name=excluded.name,
		label=excluded.label,
		environment=excluded.environment
RETURNING id;`
