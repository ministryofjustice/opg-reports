// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: queries.sql

package awsc

import (
	"context"
)

const count = `-- name: Count :one
SELECT count(*) FROM aws_costs
`

func (q *Queries) Count(ctx context.Context) (int64, error) {
	row := q.queryRow(ctx, q.countStmt, count)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const dailyCostsDetailed = `-- name: DailyCostsDetailed :many
SELECT
    account_id,
    unit,
    IIF(environment != "null", environment, "production") as environment,
    service,
    coalesce(SUM(cost), 0) as total,
    strftime("%Y-%m-%d", date) as interval
FROM aws_costs
WHERE
    date >= ?1 AND
    date < ?2
GROUP BY strftime("%Y-%m-%d", date), account_id, unit, environment, service
ORDER by strftime("%Y-%m-%d", date) ASC
`

type DailyCostsDetailedParams struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type DailyCostsDetailedRow struct {
	AccountID   string      `json:"account_id"`
	Unit        string      `json:"unit"`
	Environment interface{} `json:"environment"`
	Service     string      `json:"service"`
	Total       interface{} `json:"total"`
	Interval    interface{} `json:"interval"`
}

func (q *Queries) DailyCostsDetailed(ctx context.Context, arg DailyCostsDetailedParams) ([]DailyCostsDetailedRow, error) {
	rows, err := q.query(ctx, q.dailyCostsDetailedStmt, dailyCostsDetailed, arg.Start, arg.End)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DailyCostsDetailedRow
	for rows.Next() {
		var i DailyCostsDetailedRow
		if err := rows.Scan(
			&i.AccountID,
			&i.Unit,
			&i.Environment,
			&i.Service,
			&i.Total,
			&i.Interval,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const dailyCostsDetailedForUnit = `-- name: DailyCostsDetailedForUnit :many
SELECT
    account_id,
    unit,
    IIF(environment != "null", environment, "production") as environment,
    service,
    coalesce(SUM(cost), 0) as total,
    strftime("%Y-%m-%d", date) as interval
FROM aws_costs
WHERE
    date >= ?1 AND
    date < ?2 AND
    unit = ?3
GROUP BY strftime("%Y-%m-%d", date), account_id, environment, service
ORDER by strftime("%Y-%m-%d", date) ASC
`

type DailyCostsDetailedForUnitParams struct {
	Start string `json:"start"`
	End   string `json:"end"`
	Unit  string `json:"unit"`
}

type DailyCostsDetailedForUnitRow struct {
	AccountID   string      `json:"account_id"`
	Unit        string      `json:"unit"`
	Environment interface{} `json:"environment"`
	Service     string      `json:"service"`
	Total       interface{} `json:"total"`
	Interval    interface{} `json:"interval"`
}

func (q *Queries) DailyCostsDetailedForUnit(ctx context.Context, arg DailyCostsDetailedForUnitParams) ([]DailyCostsDetailedForUnitRow, error) {
	rows, err := q.query(ctx, q.dailyCostsDetailedForUnitStmt, dailyCostsDetailedForUnit, arg.Start, arg.End, arg.Unit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DailyCostsDetailedForUnitRow
	for rows.Next() {
		var i DailyCostsDetailedForUnitRow
		if err := rows.Scan(
			&i.AccountID,
			&i.Unit,
			&i.Environment,
			&i.Service,
			&i.Total,
			&i.Interval,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const dailyCostsPerUnit = `-- name: DailyCostsPerUnit :many
SELECT
    unit,
    coalesce(SUM(cost), 0) as total,
    strftime("%Y-%m-%d", date) as interval
FROM aws_costs
WHERE
    date >= ?1 AND
    date < ?2
GROUP BY strftime("%Y-%m-%d", date), unit
ORDER by strftime("%Y-%m-%d", date) ASC
`

type DailyCostsPerUnitParams struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type DailyCostsPerUnitRow struct {
	Unit     string      `json:"unit"`
	Total    interface{} `json:"total"`
	Interval interface{} `json:"interval"`
}

func (q *Queries) DailyCostsPerUnit(ctx context.Context, arg DailyCostsPerUnitParams) ([]DailyCostsPerUnitRow, error) {
	rows, err := q.query(ctx, q.dailyCostsPerUnitStmt, dailyCostsPerUnit, arg.Start, arg.End)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DailyCostsPerUnitRow
	for rows.Next() {
		var i DailyCostsPerUnitRow
		if err := rows.Scan(&i.Unit, &i.Total, &i.Interval); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const dailyCostsPerUnitEnvironment = `-- name: DailyCostsPerUnitEnvironment :many
SELECT
    unit,
    IIF(environment != "null", environment, "production") as environment,
    coalesce(SUM(cost), 0) as total,
    strftime("%Y-%m-%d", date) as interval
FROM aws_costs
WHERE
    date >= ?1 AND
    date < ?2
GROUP BY strftime("%Y-%m-%d", date), unit, environment
ORDER by strftime("%Y-%m-%d", date) ASC
`

type DailyCostsPerUnitEnvironmentParams struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type DailyCostsPerUnitEnvironmentRow struct {
	Unit        string      `json:"unit"`
	Environment interface{} `json:"environment"`
	Total       interface{} `json:"total"`
	Interval    interface{} `json:"interval"`
}

func (q *Queries) DailyCostsPerUnitEnvironment(ctx context.Context, arg DailyCostsPerUnitEnvironmentParams) ([]DailyCostsPerUnitEnvironmentRow, error) {
	rows, err := q.query(ctx, q.dailyCostsPerUnitEnvironmentStmt, dailyCostsPerUnitEnvironment, arg.Start, arg.End)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DailyCostsPerUnitEnvironmentRow
	for rows.Next() {
		var i DailyCostsPerUnitEnvironmentRow
		if err := rows.Scan(
			&i.Unit,
			&i.Environment,
			&i.Total,
			&i.Interval,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insert = `-- name: Insert :one
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
) RETURNING id
`

type InsertParams struct {
	Ts           string `json:"ts"`
	Organisation string `json:"organisation"`
	AccountID    string `json:"account_id"`
	AccountName  string `json:"account_name"`
	Unit         string `json:"unit"`
	Label        string `json:"label"`
	Environment  string `json:"environment"`
	Service      string `json:"service"`
	Region       string `json:"region"`
	Date         string `json:"date"`
	Cost         string `json:"cost"`
}

func (q *Queries) Insert(ctx context.Context, arg InsertParams) (int, error) {
	row := q.queryRow(ctx, q.insertStmt, insert,
		arg.Ts,
		arg.Organisation,
		arg.AccountID,
		arg.AccountName,
		arg.Unit,
		arg.Label,
		arg.Environment,
		arg.Service,
		arg.Region,
		arg.Date,
		arg.Cost,
	)
	var id int
	err := row.Scan(&id)
	return id, err
}

const monthlyCostsDetailed = `-- name: MonthlyCostsDetailed :many
SELECT
    account_id,
    unit,
    IIF(environment != "null", environment, "production") as environment,
    service,
    coalesce(SUM(cost), 0) as total,
    strftime("%Y-%m", date) as interval
FROM aws_costs
WHERE
    date >= ?1 AND
    date < ?2
GROUP BY strftime("%Y-%m", date), account_id, unit, environment, service
ORDER by strftime("%Y-%m", date) ASC
`

type MonthlyCostsDetailedParams struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type MonthlyCostsDetailedRow struct {
	AccountID   string      `json:"account_id"`
	Unit        string      `json:"unit"`
	Environment interface{} `json:"environment"`
	Service     string      `json:"service"`
	Total       interface{} `json:"total"`
	Interval    interface{} `json:"interval"`
}

func (q *Queries) MonthlyCostsDetailed(ctx context.Context, arg MonthlyCostsDetailedParams) ([]MonthlyCostsDetailedRow, error) {
	rows, err := q.query(ctx, q.monthlyCostsDetailedStmt, monthlyCostsDetailed, arg.Start, arg.End)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MonthlyCostsDetailedRow
	for rows.Next() {
		var i MonthlyCostsDetailedRow
		if err := rows.Scan(
			&i.AccountID,
			&i.Unit,
			&i.Environment,
			&i.Service,
			&i.Total,
			&i.Interval,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const monthlyCostsDetailedForUnit = `-- name: MonthlyCostsDetailedForUnit :many
SELECT
    account_id,
    unit,
    IIF(environment != "null", environment, "production") as environment,
    service,
    coalesce(SUM(cost), 0) as total,
    strftime("%Y-%m", date) as interval
FROM aws_costs
WHERE
    date >= ?1 AND
    date < ?2 AND
    unit = ?3
GROUP BY strftime("%Y-%m", date), account_id, unit, environment, service
ORDER by strftime("%Y-%m", date) ASC
`

type MonthlyCostsDetailedForUnitParams struct {
	Start string `json:"start"`
	End   string `json:"end"`
	Unit  string `json:"unit"`
}

type MonthlyCostsDetailedForUnitRow struct {
	AccountID   string      `json:"account_id"`
	Unit        string      `json:"unit"`
	Environment interface{} `json:"environment"`
	Service     string      `json:"service"`
	Total       interface{} `json:"total"`
	Interval    interface{} `json:"interval"`
}

func (q *Queries) MonthlyCostsDetailedForUnit(ctx context.Context, arg MonthlyCostsDetailedForUnitParams) ([]MonthlyCostsDetailedForUnitRow, error) {
	rows, err := q.query(ctx, q.monthlyCostsDetailedForUnitStmt, monthlyCostsDetailedForUnit, arg.Start, arg.End, arg.Unit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MonthlyCostsDetailedForUnitRow
	for rows.Next() {
		var i MonthlyCostsDetailedForUnitRow
		if err := rows.Scan(
			&i.AccountID,
			&i.Unit,
			&i.Environment,
			&i.Service,
			&i.Total,
			&i.Interval,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const monthlyCostsPerUnit = `-- name: MonthlyCostsPerUnit :many
SELECT
    unit,
    coalesce(SUM(cost), 0) as total,
    strftime("%Y-%m", date) as interval
FROM aws_costs
WHERE
    date >= ?1 AND
    date < ?2
GROUP BY strftime("%Y-%m", date), unit
ORDER by strftime("%Y-%m", date) ASC
`

type MonthlyCostsPerUnitParams struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type MonthlyCostsPerUnitRow struct {
	Unit     string      `json:"unit"`
	Total    interface{} `json:"total"`
	Interval interface{} `json:"interval"`
}

func (q *Queries) MonthlyCostsPerUnit(ctx context.Context, arg MonthlyCostsPerUnitParams) ([]MonthlyCostsPerUnitRow, error) {
	rows, err := q.query(ctx, q.monthlyCostsPerUnitStmt, monthlyCostsPerUnit, arg.Start, arg.End)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MonthlyCostsPerUnitRow
	for rows.Next() {
		var i MonthlyCostsPerUnitRow
		if err := rows.Scan(&i.Unit, &i.Total, &i.Interval); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const monthlyCostsPerUnitEnvironment = `-- name: MonthlyCostsPerUnitEnvironment :many
SELECT
    unit,
    IIF(environment != "null", environment, "production") as environment,
    coalesce(SUM(cost), 0) as total,
    strftime("%Y-%m", date) as interval
FROM aws_costs
WHERE
    date >= ?1 AND
    date < ?2
GROUP BY strftime("%Y-%m", date), unit, environment
ORDER by strftime("%Y-%m", date) ASC
`

type MonthlyCostsPerUnitEnvironmentParams struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type MonthlyCostsPerUnitEnvironmentRow struct {
	Unit        string      `json:"unit"`
	Environment interface{} `json:"environment"`
	Total       interface{} `json:"total"`
	Interval    interface{} `json:"interval"`
}

func (q *Queries) MonthlyCostsPerUnitEnvironment(ctx context.Context, arg MonthlyCostsPerUnitEnvironmentParams) ([]MonthlyCostsPerUnitEnvironmentRow, error) {
	rows, err := q.query(ctx, q.monthlyCostsPerUnitEnvironmentStmt, monthlyCostsPerUnitEnvironment, arg.Start, arg.End)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MonthlyCostsPerUnitEnvironmentRow
	for rows.Next() {
		var i MonthlyCostsPerUnitEnvironmentRow
		if err := rows.Scan(
			&i.Unit,
			&i.Environment,
			&i.Total,
			&i.Interval,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const monthlyTotalsTaxSplit = `-- name: MonthlyTotalsTaxSplit :many
SELECT
    'Including Tax' as service,
    coalesce(SUM(cost), 0) as total,
    strftime("%Y-%m", date) as interval
FROM aws_costs as incTax
WHERE
    incTax.date >= ?1
    AND incTax.date < ?2
GROUP BY strftime("%Y-%m", incTax.date)
UNION ALL
SELECT
    'Excluding Tax' as service,
    coalesce(SUM(cost), 0) as total,
    strftime("%Y-%m", date) as interval
FROM aws_costs as excTax
WHERE
    excTax.service != 'Tax'
    AND excTax.date >= ?1
    AND excTax.date < ?2
GROUP BY strftime("%Y-%m", date)
ORDER by interval ASC
`

type MonthlyTotalsTaxSplitParams struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type MonthlyTotalsTaxSplitRow struct {
	Service  string      `json:"service"`
	Total    interface{} `json:"total"`
	Interval interface{} `json:"interval"`
}

func (q *Queries) MonthlyTotalsTaxSplit(ctx context.Context, arg MonthlyTotalsTaxSplitParams) ([]MonthlyTotalsTaxSplitRow, error) {
	rows, err := q.query(ctx, q.monthlyTotalsTaxSplitStmt, monthlyTotalsTaxSplit, arg.Start, arg.End)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MonthlyTotalsTaxSplitRow
	for rows.Next() {
		var i MonthlyTotalsTaxSplitRow
		if err := rows.Scan(&i.Service, &i.Total, &i.Interval); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const oldest = `-- name: Oldest :one
;
SELECT run_date FROM aws_costs_tracker ORDER BY run_date ASC LIMIT 1
`

func (q *Queries) Oldest(ctx context.Context) (string, error) {
	row := q.queryRow(ctx, q.oldestStmt, oldest)
	var run_date string
	err := row.Scan(&run_date)
	return run_date, err
}

const total = `-- name: Total :one
SELECT
    coalesce(SUM(cost), 0) as total
FROM aws_costs
WHERE
    date >= ?1 AND date < ?2
`

type TotalParams struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

func (q *Queries) Total(ctx context.Context, arg TotalParams) (interface{}, error) {
	row := q.queryRow(ctx, q.totalStmt, total, arg.Start, arg.End)
	var total interface{}
	err := row.Scan(&total)
	return total, err
}

const track = `-- name: Track :exec
INSERT INTO aws_costs_tracker (run_date) VALUES(?)
`

func (q *Queries) Track(ctx context.Context, runDate string) error {
	_, err := q.exec(ctx, q.trackStmt, track, runDate)
	return err
}

const youngest = `-- name: Youngest :one
SELECT run_date FROM aws_costs_tracker ORDER BY run_date DESC LIMIT 1
`

func (q *Queries) Youngest(ctx context.Context) (string, error) {
	row := q.queryRow(ctx, q.youngestStmt, youngest)
	var run_date string
	err := row.Scan(&run_date)
	return run_date, err
}
