package awscosts

import (
	"context"

	"github.com/jmoiron/sqlx"
)

const createCostTable string = `CREATE TABLE IF NOT EXISTS aws_costs (
    id INTEGER PRIMARY KEY,
    ts TEXT NOT NULL,

    organisation TEXT NOT NULL,
    account_id TEXT NOT NULL,
    account_name TEXT NOT NULL,
    unit TEXT NOT NULL,
    label TEXT NOT NULL,
    environment TEXT NOT NULL,

	region TEXT NOT NULL,
    service TEXT NOT NULL,
    date TEXT NOT NULL,
    cost TEXT NOT NULL
    --
) STRICT;`
const createCostTableIndex string = `CREATE INDEX IF NOT EXISTS aws_costs_date_idx ON aws_costs(date);`

// Create will run the table and index creation for aws_costs table
// It uses the MustExecContext function, so any errors will throw a
// panic
func Create(ctx context.Context, db *sqlx.DB) {
	db.MustExecContext(ctx, createCostTable)
	db.MustExecContext(ctx, createCostTableIndex)
}
