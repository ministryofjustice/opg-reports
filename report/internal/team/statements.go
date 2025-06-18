package team

// list of sql insert statements of various types
const stmtSelectAll string = `
SELECT
	id,
	name,
	created_at
FROM team
ORDER BY name ASC;
`
const stmtInsert string = `
INSERT INTO team (
	name,
	created_at
) VALUES (
	:name,
	:created_at
) ON CONFLICT (name)
 	DO UPDATE SET name=excluded.name
RETURNING id;`
