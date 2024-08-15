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

-- name: Track :exec
INSERT INTO aws_costs_tracker (run_date) VALUES(?) ;

-- name: Oldest :one
SELECT MIN(date) FROM aws_costs_tracker LIMIT 1;

-- name: Youngest :one
SELECT MAX(date) FROM aws_costs_tracker LIMIT 1;
