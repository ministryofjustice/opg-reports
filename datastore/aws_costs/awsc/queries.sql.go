// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: queries.sql

package awsc

import (
	"context"
	"database/sql"
)

const byMonth = `-- name: ByMonth :many
SELECT
    SUM(cost) as total,
    strftime("%Y-%m", date) as month
FROM aws_costs
GROUP BY strftime("%Y-%m", date)
`

type ByMonthRow struct {
	Total sql.NullFloat64 `json:"total"`
	Month interface{}     `json:"month"`
}

func (q *Queries) ByMonth(ctx context.Context) ([]ByMonthRow, error) {
	rows, err := q.query(ctx, q.byMonthStmt, byMonth)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ByMonthRow
	for rows.Next() {
		var i ByMonthRow
		if err := rows.Scan(&i.Total, &i.Month); err != nil {
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

const count = `-- name: Count :one
SELECT count(*) FROM aws_costs
`

func (q *Queries) Count(ctx context.Context) (int64, error) {
	row := q.queryRow(ctx, q.countStmt, count)
	var count int64
	err := row.Scan(&count)
	return count, err
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

const oldest = `-- name: Oldest :one
;

SELECT MIN(date) FROM aws_costs_tracker LIMIT 1
`

func (q *Queries) Oldest(ctx context.Context) (interface{}, error) {
	row := q.queryRow(ctx, q.oldestStmt, oldest)
	var min interface{}
	err := row.Scan(&min)
	return min, err
}

const track = `-- name: Track :exec
INSERT INTO aws_costs_tracker (run_date) VALUES(?)
`

func (q *Queries) Track(ctx context.Context, runDate string) error {
	_, err := q.exec(ctx, q.trackStmt, track, runDate)
	return err
}

const youngest = `-- name: Youngest :one
SELECT MAX(date) FROM aws_costs_tracker LIMIT 1
`

func (q *Queries) Youngest(ctx context.Context) (interface{}, error) {
	row := q.queryRow(ctx, q.youngestStmt, youngest)
	var max interface{}
	err := row.Scan(&max)
	return max, err
}
