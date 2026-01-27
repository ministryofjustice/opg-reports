package uptimeimports

// insertStmt used to insert records
const insertStmt string = `
INSERT INTO aws_uptime (
	date,
	average,
	granularity,
	aws_account_id
) VALUES (
	:date,
	:average,
	:granularity,
	:aws_account_id
) ON CONFLICT (aws_account_id,date)
 	DO UPDATE SET average=excluded.average, granularity=excluded.granularity
RETURNING id;
`
