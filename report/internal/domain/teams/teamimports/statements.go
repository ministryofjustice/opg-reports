package teamimports

// insertStmt used to insert records
const insertStmt string = `
INSERT INTO teams (
	name
) VALUES (
	:name
) ON CONFLICT (name)
 	DO UPDATE SET name=excluded.name
RETURNING name
;
`
