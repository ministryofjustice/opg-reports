package awscost

// stmtDropTable deletes the table
const stmtDropTable string = `DROP TABLE IF EXISTS aws_costs;`

// stmtDeleteAll removes all records - used by fixture seeding to avoid duplicates
const stmtDeleteAll string = `DELETE FROM aws_costs;`

// stmtImport is used by "existing" commands to insert data
// while handling the joins to account data
const stmtImport string = `
INSERT INTO aws_costs (
	region,
	service,
	date,
	cost,
	aws_account_id
) SELECT
	:region,
	:service,
	:date,
	:cost,
	id
FROM aws_accounts WHERE aws_accounts.id = :account_id
ON CONFLICT (aws_account_id,date,region,service)
 	DO UPDATE SET cost=excluded.cost
RETURNING id;
`

// stmtSelectAll fetches all records and orders them by cost.
//
// This should never be exposed to the api layer as table
// has millions of rows
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
	) as aws_account,
	json_object(
		'name', aws_accounts.team_name
	) as team
FROM aws_costs
LEFT JOIN aws_accounts on aws_accounts.id = aws_costs.aws_account_id
GROUP BY aws_costs.id
ORDER BY
	CAST(aws_costs.cost as REAL) DESC,
	aws_accounts.team_name ASC,
	aws_accounts.name ASC,
	aws_accounts.environment ASC,
	aws_costs.region ASC,
	aws_costs.service ASC;
`

// stmtSelectTop20 returns the top20 most expensive costs stored.
//
// Excludes tax. as that is worked out on a single day for the
// month so would always fill this list.
const stmtSelectTop20 string = `
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
	) as aws_account,
	json_object(
		'name', aws_accounts.team_name
	) as team
FROM aws_costs
LEFT JOIN aws_accounts on aws_accounts.id = aws_costs.aws_account_id
WHERE
	aws_costs.service != 'Tax'
GROUP BY aws_costs.id
ORDER BY
	CAST(aws_costs.cost as REAL) DESC,
	aws_accounts.team_name ASC,
	aws_accounts.name ASC,
	aws_accounts.environment ASC,
	aws_costs.region ASC,
	aws_costs.service ASC
LIMIT 20;
`

// stmGroupedCosts is the base sql statements used for most cost database calls
// that filters out Tax and groups values by at least the date column.
//
// It contains extra :params to allow extension of the query and typically
// generated from an api input dataset
//
// :date_format = used for date time grouping via strftime on the date column
// :start_date	= lower bound on the date where
// :end_date 	= upper bound on the date where
// :select 		= additional columns to fetch
// :where 		= additional where queries to include
// :groupby 	= extra group by clauses
// :orderby 	= extra order by options
const stmGroupedCosts string = `
SELECT
	{SELECT}
    coalesce(SUM(cost), 0) as cost,
    strftime(:date_format, date) as date
FROM aws_costs
LEFT JOIN aws_accounts ON aws_accounts.id = aws_costs.aws_account_id
WHERE
	{WHERE}
    date >= :start_date
    AND date < :end_date
GROUP BY
	{GROUP_BY}
	strftime(:date_format, date)
ORDER BY
	CAST(aws_costs.cost as REAL) DESC,
	{ORDER_BY}
	strftime(:date_format, date) ASC
;
`
