package teamimport

const importStatement string = `
INSERT INTO teams (
	name
) VALUES (
	:name
) ON CONFLICT (name)
 	DO UPDATE SET name=excluded.name
RETURNING name
;
`
