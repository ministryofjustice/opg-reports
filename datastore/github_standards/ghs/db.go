// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package ghs

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
	if q.ageStmt, err = db.PrepareContext(ctx, age); err != nil {
		return nil, fmt.Errorf("error preparing query Age: %w", err)
	}
	if q.allStmt, err = db.PrepareContext(ctx, all); err != nil {
		return nil, fmt.Errorf("error preparing query All: %w", err)
	}
	if q.countStmt, err = db.PrepareContext(ctx, count); err != nil {
		return nil, fmt.Errorf("error preparing query Count: %w", err)
	}
	if q.filterByIsArchivedStmt, err = db.PrepareContext(ctx, filterByIsArchived); err != nil {
		return nil, fmt.Errorf("error preparing query FilterByIsArchived: %w", err)
	}
	if q.filterByIsArchivedAndTeamStmt, err = db.PrepareContext(ctx, filterByIsArchivedAndTeam); err != nil {
		return nil, fmt.Errorf("error preparing query FilterByIsArchivedAndTeam: %w", err)
	}
	if q.filterByTeamStmt, err = db.PrepareContext(ctx, filterByTeam); err != nil {
		return nil, fmt.Errorf("error preparing query FilterByTeam: %w", err)
	}
	if q.insertStmt, err = db.PrepareContext(ctx, insert); err != nil {
		return nil, fmt.Errorf("error preparing query Insert: %w", err)
	}
	if q.totalCountCompliantBaselineStmt, err = db.PrepareContext(ctx, totalCountCompliantBaseline); err != nil {
		return nil, fmt.Errorf("error preparing query TotalCountCompliantBaseline: %w", err)
	}
	if q.totalCountCompliantExtendedStmt, err = db.PrepareContext(ctx, totalCountCompliantExtended); err != nil {
		return nil, fmt.Errorf("error preparing query TotalCountCompliantExtended: %w", err)
	}
	if q.trackStmt, err = db.PrepareContext(ctx, track); err != nil {
		return nil, fmt.Errorf("error preparing query Track: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.ageStmt != nil {
		if cerr := q.ageStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing ageStmt: %w", cerr)
		}
	}
	if q.allStmt != nil {
		if cerr := q.allStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing allStmt: %w", cerr)
		}
	}
	if q.countStmt != nil {
		if cerr := q.countStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing countStmt: %w", cerr)
		}
	}
	if q.filterByIsArchivedStmt != nil {
		if cerr := q.filterByIsArchivedStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing filterByIsArchivedStmt: %w", cerr)
		}
	}
	if q.filterByIsArchivedAndTeamStmt != nil {
		if cerr := q.filterByIsArchivedAndTeamStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing filterByIsArchivedAndTeamStmt: %w", cerr)
		}
	}
	if q.filterByTeamStmt != nil {
		if cerr := q.filterByTeamStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing filterByTeamStmt: %w", cerr)
		}
	}
	if q.insertStmt != nil {
		if cerr := q.insertStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing insertStmt: %w", cerr)
		}
	}
	if q.totalCountCompliantBaselineStmt != nil {
		if cerr := q.totalCountCompliantBaselineStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing totalCountCompliantBaselineStmt: %w", cerr)
		}
	}
	if q.totalCountCompliantExtendedStmt != nil {
		if cerr := q.totalCountCompliantExtendedStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing totalCountCompliantExtendedStmt: %w", cerr)
		}
	}
	if q.trackStmt != nil {
		if cerr := q.trackStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing trackStmt: %w", cerr)
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
	db                              DBTX
	tx                              *sql.Tx
	ageStmt                         *sql.Stmt
	allStmt                         *sql.Stmt
	countStmt                       *sql.Stmt
	filterByIsArchivedStmt          *sql.Stmt
	filterByIsArchivedAndTeamStmt   *sql.Stmt
	filterByTeamStmt                *sql.Stmt
	insertStmt                      *sql.Stmt
	totalCountCompliantBaselineStmt *sql.Stmt
	totalCountCompliantExtendedStmt *sql.Stmt
	trackStmt                       *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                              tx,
		tx:                              tx,
		ageStmt:                         q.ageStmt,
		allStmt:                         q.allStmt,
		countStmt:                       q.countStmt,
		filterByIsArchivedStmt:          q.filterByIsArchivedStmt,
		filterByIsArchivedAndTeamStmt:   q.filterByIsArchivedAndTeamStmt,
		filterByTeamStmt:                q.filterByTeamStmt,
		insertStmt:                      q.insertStmt,
		totalCountCompliantBaselineStmt: q.totalCountCompliantBaselineStmt,
		totalCountCompliantExtendedStmt: q.totalCountCompliantExtendedStmt,
		trackStmt:                       q.trackStmt,
	}
}