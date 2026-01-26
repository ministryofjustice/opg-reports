package accountimports

// insertStmt used to insert records
const insertStmt string = `
INSERT INTO aws_accounts (
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
