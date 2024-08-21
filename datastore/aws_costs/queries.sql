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
    strftime("%Y-%m", date) as month
FROM aws_costs as incTax
WHERE
    incTax.date >= @start
    AND incTax.date < @end
GROUP BY strftime("%Y-%m", incTax.date)
UNION ALL
SELECT
    'Excluding Tax' as service,
    coalesce(SUM(cost), 0) as total,
    strftime("%Y-%m", date) as month
FROM aws_costs as excTax
WHERE
    excTax.service != 'Tax'
    AND excTax.date >= @start
    AND excTax.date < @end
GROUP BY strftime("%Y-%m", date)
ORDER by month ASC;

-- name: Total :one
SELECT
    coalesce(SUM(cost), 0) as total
FROM aws_costs
WHERE
    date >= @start AND date < @end
;
-- -- name: ByMonth :many
-- SELECT
--     SUM(cost) as total,
--     strftime("%Y-%m", date) as month
-- FROM aws_costs
-- GROUP BY strftime("%Y-%m", date);

-- name: Track :exec
INSERT INTO aws_costs_tracker (run_date) VALUES(?) ;

-- name: Oldest :one
SELECT run_date FROM aws_costs_tracker ORDER BY run_date ASC LIMIT 1;

-- name: Youngest :one
SELECT run_date FROM aws_costs_tracker ORDER BY run_date DESC LIMIT 1;
