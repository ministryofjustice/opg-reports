package team

// list of sql insert statements of various types
const stmtDropTable string = `DROP TABLE IF EXISTS teams`
const stmtSelectAll string = `
SELECT
	id,
	name,
	created_at
FROM teams
ORDER BY name ASC;`
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
