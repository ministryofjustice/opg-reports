-- name: Insert :one
INSERT INTO aws_costs(
    ts,
    organisation,
    account_id,
    account_name,
    unit,
    label,
    environment,
    service,
    region,
    date,
    cost
) VALUES (
    ?,?,?,?,?,?,?,?,?,?,?
) RETURNING id;

-- name: Count :one
SELECT count(*) FROM aws_costs;

-- name: MonthlyTotalsTaxSplit :many
SELECT
    'Including Tax' as service,
    coalesce(SUM(cost), 0) as total,
    strftime("%Y-%m", date) as interval
FROM aws_costs as incTax
WHERE
    incTax.date >= @start
    AND incTax.date < @end
GROUP BY strftime("%Y-%m", incTax.date)
UNION ALL
SELECT
    'Excluding Tax' as service,
    coalesce(SUM(cost), 0) as total,
    strftime("%Y-%m", date) as interval
FROM aws_costs as excTax
WHERE
    excTax.service != 'Tax'
    AND excTax.date >= @start
    AND excTax.date < @end
GROUP BY strftime("%Y-%m", date)
ORDER by interval ASC;

-- name: Total :one
SELECT
    coalesce(SUM(cost), 0) as total
FROM aws_costs
WHERE
    date >= @start AND date < @end;

-- name: MonthlyCostsPerUnit :many
SELECT
    unit,
    coalesce(SUM(cost), 0) as total,
    strftime("%Y-%m", date) as interval
FROM aws_costs
WHERE
    date >= @start AND
    date < @end
GROUP BY strftime("%Y-%m", date), unit
ORDER by strftime("%Y-%m", date) ASC;

-- name: MonthlyCostsPerUnitEnvironment :many
SELECT
    unit,
    IIF(environment != "null", environment, "production") as environment,
    coalesce(SUM(cost), 0) as total,
    strftime("%Y-%m", date) as interval
FROM aws_costs
WHERE
    date >= @start AND
    date < @end
GROUP BY strftime("%Y-%m", date), unit, environment
ORDER by strftime("%Y-%m", date) ASC;

-- name: MonthlyCostsDetailed :many
SELECT
    account_id,
    unit,
    IIF(environment != "null", environment, "production") as environment,
    service,
    coalesce(SUM(cost), 0) as total,
    strftime("%Y-%m", date) as interval
FROM aws_costs
WHERE
    date >= @start AND
    date < @end
GROUP BY strftime("%Y-%m", date), account_id, unit, environment, service
ORDER by strftime("%Y-%m", date) ASC;

-- name: MonthlyCostsDetailedForUnit :many
SELECT
    account_id,
    unit,
    IIF(environment != "null", environment, "production") as environment,
    service,
    coalesce(SUM(cost), 0) as total,
    strftime("%Y-%m", date) as interval
FROM aws_costs
WHERE
    date >= @start AND
    date < @end AND
    unit = @unit
GROUP BY strftime("%Y-%m", date), account_id, unit, environment, service
ORDER by strftime("%Y-%m", date) ASC;

-- name: DailyCostsPerUnit :many
SELECT
    unit,
    coalesce(SUM(cost), 0) as total,
    strftime("%Y-%m-%d", date) as interval
FROM aws_costs
WHERE
    date >= @start AND
    date < @end
GROUP BY strftime("%Y-%m-%d", date), unit
ORDER by strftime("%Y-%m-%d", date) ASC;

-- name: DailyCostsPerUnitEnvironment :many
SELECT
    unit,
    IIF(environment != "null", environment, "production") as environment,
    coalesce(SUM(cost), 0) as total,
    strftime("%Y-%m-%d", date) as interval
FROM aws_costs
WHERE
    date >= @start AND
    date < @end
GROUP BY strftime("%Y-%m-%d", date), unit, environment
ORDER by strftime("%Y-%m-%d", date) ASC;

-- name: DailyCostsDetailed :many
SELECT
    account_id,
    unit,
    IIF(environment != "null", environment, "production") as environment,
    service,
    coalesce(SUM(cost), 0) as total,
    strftime("%Y-%m-%d", date) as interval
FROM aws_costs
WHERE
    date >= @start AND
    date < @end
GROUP BY strftime("%Y-%m-%d", date), account_id, unit, environment, service
ORDER by strftime("%Y-%m-%d", date) ASC;


-- name: DailyCostsDetailedForUnit :many
SELECT
    account_id,
    unit,
    IIF(environment != "null", environment, "production") as environment,
    service,
    coalesce(SUM(cost), 0) as total,
    strftime("%Y-%m-%d", date) as interval
FROM aws_costs
WHERE
    date >= @start AND
    date < @end AND
    unit = @unit
GROUP BY strftime("%Y-%m-%d", date), account_id, environment, service
ORDER by strftime("%Y-%m-%d", date) ASC;


-- name: Track :exec
INSERT INTO aws_costs_tracker (run_date) VALUES(?) ;
-- name: Oldest :one
SELECT run_date FROM aws_costs_tracker ORDER BY run_date ASC LIMIT 1;
-- name: Youngest :one
SELECT run_date FROM aws_costs_tracker ORDER BY run_date DESC LIMIT 1;
