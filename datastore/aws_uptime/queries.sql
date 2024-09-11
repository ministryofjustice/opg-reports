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

-- name: UptimePerMonth :many
SELECT
    (coalesce(SUM(average), 0) / count(*) ) as average,
    strftime("%Y-%m", date) as interval
FROM aws_uptime
WHERE
    date >= @start AND
    date < @end
GROUP BY strftime("%Y-%m", date)
ORDER by strftime("%Y-%m", date) ASC;

-- name: UptimePerMonthUnit :many
SELECT
    (coalesce(SUM(average), 0) / count(*) ) as average,
    strftime("%Y-%m", date) as interval,
    unit
FROM aws_uptime
WHERE
    date >= @start AND
    date < @end
GROUP BY strftime("%Y-%m", date), unit
ORDER by strftime("%Y-%m", date) ASC;

-- name: UptimePerMonthFilterByUnit :many
SELECT
    (coalesce(SUM(average), 0) / count(*) ) as average,
    strftime("%Y-%m", date) as interval,
    unit
FROM aws_uptime
WHERE
    date >= @start AND
    date < @end AND
    unit = @unit
GROUP BY strftime("%Y-%m", date), unit
ORDER by strftime("%Y-%m", date) ASC;

-- name: UptimePerDay :many
SELECT
    (coalesce(SUM(average), 0) / count(*) ) as average,
    strftime("%Y-%m-%d", date) as interval
FROM aws_uptime
WHERE
    date >= @start AND
    date < @end
GROUP BY strftime("%Y-%m-%d", date)
ORDER by strftime("%Y-%m-%d", date) ASC;

-- name: UptimePerDayUnit :many
SELECT
    (coalesce(SUM(average), 0) / count(*) ) as average,
    strftime("%Y-%m-%d", date) as interval,
    unit
FROM aws_uptime
WHERE
    date >= @start AND
    date < @end
GROUP BY strftime("%Y-%m-%d", date), unit
ORDER by strftime("%Y-%m-%d", date) ASC;

-- name: UptimePerDayFilterByUnit :many
SELECT
    (coalesce(SUM(average), 0) / count(*) ) as average,
    strftime("%Y-%m-%d", date) as interval,
    unit
FROM aws_uptime
WHERE
    date >= @start AND
    date < @end AND
    unit = @unit
GROUP BY strftime("%Y-%m-%d", date), unit
ORDER by strftime("%Y-%m-%d", date) ASC;


-- name: Track :exec
INSERT INTO aws_uptime_tracker (run_date) VALUES(?) ;
-- name: Oldest :one
SELECT run_date FROM aws_uptime_tracker ORDER BY run_date ASC LIMIT 1;
-- name: Youngest :one
SELECT run_date FROM aws_uptime_tracker ORDER BY run_date DESC LIMIT 1;
