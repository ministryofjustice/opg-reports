package codebaseimports

// insertStmt used to insert records
const insertStmt string = `
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
