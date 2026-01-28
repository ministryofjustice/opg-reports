package uptimeimports

// insertStmt used to insert records
const insertStmt string = `
INSERT INTO uptime (
	date,
	average,
	granularity,
	account_id
) VALUES (
	:date,
	:average,
	:granularity,
	:account_id
) ON CONFLICT (account_id,date)
 	DO UPDATE SET average=excluded.average, granularity=excluded.granularity
RETURNING id;
`
