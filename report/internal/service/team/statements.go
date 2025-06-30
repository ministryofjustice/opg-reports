package team

// list of sql insert statements of various types
const stmtDropTable string = `DROP TABLE IF EXISTS teams`
const stmtDeleteAll string = `DELETE FROM teams`

// stmtImport is used by "existing" commands to insert data
const stmtImport string = `INSERT INTO teams (name) VALUES (:name) ON CONFLICT (name) DO UPDATE SET name=excluded.name RETURNING name;`

// stmtSelectAll and join the account data onto the team list
const stmtSelectAll string = `
SELECT
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
LEFT JOIN aws_accounts ON aws_accounts.team_name = teams.name
GROUP BY teams.name
ORDER BY teams.name ASC`
