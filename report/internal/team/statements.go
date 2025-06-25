package team

// list of sql insert statements of various types
const stmtDropTable string = `DROP TABLE IF EXISTS teams`
const stmtSelectAll string = `
SELECT
	teams.id,
	teams.name,
	json_group_array(
		DISTINCT json_object(
			'id', aws_accounts.id,
			'name', aws_accounts.name,
			'label', aws_accounts.label,
			'environment', aws_accounts.environment
		)
	) filter ( where aws_accounts.id is not null) as aws_accounts
FROM teams
LEFT JOIN aws_accounts ON aws_accounts.team_id = teams.id
GROUP BY teams.id
ORDER BY teams.name ASC`

const stmtImport string = `
INSERT INTO teams (
	name,
	created_at
) VALUES (
	:name,
	:created_at
) ON CONFLICT (name)
 	DO UPDATE SET name=excluded.name
RETURNING id;`

const stmtInsert string = `
INSERT INTO teams (
	name,
	created_at
) VALUES (
	:name,
	:created_at
) ON CONFLICT (name)
 	DO UPDATE SET name=excluded.name
RETURNING id;`
