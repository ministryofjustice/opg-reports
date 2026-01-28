package infracostimports

// insertStmt used to insert records
const insertStmt string = `
INSERT INTO infracosts (
	region,
	service,
	date,
	cost,
	account_id
) VALUES (
	:region,
	:service,
	:date,
	:cost,
	:account_id
) ON CONFLICT (account_id,date,region,service)
 	DO UPDATE SET cost=excluded.cost
RETURNING id;
`
