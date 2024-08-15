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
    'WithTax' as service,
    coalesce(SUM(cost), 0) as total,
    strftime("%Y-%m", date) as month
FROM aws_costs
GROUP BY strftime("%Y-%m", date)
UNION
SELECT
    'WithoutTax' as service,
    coalesce(SUM(cost), 0) as total,
    strftime("%Y-%m", date) as month
FROM aws_costs
WHERE service != 'Tax'
GROUP BY strftime("%Y-%m", date);

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
