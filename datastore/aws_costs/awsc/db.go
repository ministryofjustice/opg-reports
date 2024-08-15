// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package awsc

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.countStmt, err = db.PrepareContext(ctx, count); err != nil {
		return nil, fmt.Errorf("error preparing query Count: %w", err)
	}
	if q.insertStmt, err = db.PrepareContext(ctx, insert); err != nil {
		return nil, fmt.Errorf("error preparing query Insert: %w", err)
	}
	if q.oldestStmt, err = db.PrepareContext(ctx, oldest); err != nil {
		return nil, fmt.Errorf("error preparing query Oldest: %w", err)
	}
	if q.trackStmt, err = db.PrepareContext(ctx, track); err != nil {
		return nil, fmt.Errorf("error preparing query Track: %w", err)
	}
	if q.youngestStmt, err = db.PrepareContext(ctx, youngest); err != nil {
		return nil, fmt.Errorf("error preparing query Youngest: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.countStmt != nil {
		if cerr := q.countStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing countStmt: %w", cerr)
		}
	}
	if q.insertStmt != nil {
		if cerr := q.insertStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing insertStmt: %w", cerr)
		}
	}
	if q.oldestStmt != nil {
		if cerr := q.oldestStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing oldestStmt: %w", cerr)
		}
	}
	if q.trackStmt != nil {
		if cerr := q.trackStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing trackStmt: %w", cerr)
		}
	}
	if q.youngestStmt != nil {
		if cerr := q.youngestStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing youngestStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db           DBTX
	tx           *sql.Tx
	countStmt    *sql.Stmt
	insertStmt   *sql.Stmt
	oldestStmt   *sql.Stmt
	trackStmt    *sql.Stmt
	youngestStmt *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:           tx,
		tx:           tx,
		countStmt:    q.countStmt,
		insertStmt:   q.insertStmt,
		oldestStmt:   q.oldestStmt,
		trackStmt:    q.trackStmt,
		youngestStmt: q.youngestStmt,
	}
}
