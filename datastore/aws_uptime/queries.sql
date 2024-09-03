-- name: Insert :one
INSERT INTO aws_uptime(
    ts,
    unit,
    date,
    average
) VALUES (
    ?,?,?,?
) RETURNING id;

-- name: Count :one
SELECT count(*) FROM aws_uptime;



-- name: Track :exec
INSERT INTO aws_uptime_tracker (run_date) VALUES(?) ;
-- name: Oldest :one
SELECT run_date FROM aws_uptime_tracker ORDER BY run_date ASC LIMIT 1;
-- name: Youngest :one
SELECT run_date FROM aws_uptime_tracker ORDER BY run_date DESC LIMIT 1;
