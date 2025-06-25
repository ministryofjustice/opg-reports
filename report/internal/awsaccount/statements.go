package awsaccount

// list of sql insert statements of various types
const stmtDropTable string = `DROP TABLE IF EXISTS awsaccounts`
const stmtSelectAll string = `
SELECT
	id,
	name,
	label,
	environment,
	created_at
FROM awsaccounts
ORDER BY id ASC;`

const stmtImport string = `
INSERT INTO awsaccounts (
	id,
	name,
	label,
	environment,
	created_at,
	team_id
) SELECT
	:id,
	:name,
	:label,
	:environment,
	:created_at,
	id
FROM teams WHERE name=:billing_unit LIMIT 1
ON CONFLICT (id)
 	DO UPDATE SET
		name=excluded.name,
		label=excluded.label,
		environment=excluded.environment
RETURNING id;`

const stmtInsert string = `
INSERT INTO awsaccounts (
	id,
	name,
	label,
	environment,
	created_at
) VALUES (
	:id,
	:name,
	:label,
	:environment,
	:created_at
) ON CONFLICT (id)
 	DO UPDATE SET name=excluded.name, label=excluded.label, environment=excluded.environment
RETURNING id;`
