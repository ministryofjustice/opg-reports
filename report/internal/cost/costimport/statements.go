package costimport

const importStatement string = `
INSERT INTO costs (
	region,
	service,
	month,
	cost,
	account_id
) VALUES (
	:region,
	:service,
	:month,
	:cost,
	:account_id
) ON CONFLICT (account_id, month, region, service)
 	DO UPDATE SET cost=excluded.cost
RETURNING id
;
`
