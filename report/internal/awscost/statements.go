package awscost

// list of sql insert statements of various types
const stmtDropTable string = `DROP TABLE IF EXISTS aws_costs`

const stmtDeleteAll string = `DELETE FROM aws_costs`

const stmtSelectAll string = `
SELECT
	aws_costs.region,
	aws_costs.service,
	aws_costs.date,
	aws_costs.cost,
	json_object(
		'id', aws_accounts.id,
		'name', aws_accounts.name,
		'label', aws_accounts.label,
		'environment', aws_accounts.environment
	) as aws_account
FROM aws_costs
LEFT JOIN aws_accounts on aws_accounts.id = aws_costs.aws_account_id
GROUP BY aws_costs.id
ORDER BY
	CAST(aws_costs.cost as REAL) DESC,
	aws_accounts.name ASC,
	aws_accounts.environment ASC,
	aws_costs.region ASC,
	aws_costs.service ASC;`

const stmtImport string = `
INSERT INTO aws_costs (
	region,
	service,
	date,
	cost,
	created_at,
	aws_account_id
) SELECT
	:region,
	:service,
	:date,
	:cost,
	:created_at,
	id
FROM aws_accounts WHERE aws_accounts.id = :account_id
ON CONFLICT (aws_account_id,date,region,service)
 	DO UPDATE SET cost=excluded.cost
RETURNING id;`
