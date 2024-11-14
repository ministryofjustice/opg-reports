package crud

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/dbs/statements"
)

type emptyNamed struct{}

// Select runs the sql statement using the adaptor and named parameter struct as value subsitutions (`struct.ID` => `:id` etc).
// If `named` is passed as nil, an empty struct is used for value replacements in the sql statement.
// Uses transactions to reduce locking
// Returns []R even when result is singular (like a COUNT(*) )
func Select[R any, A dbs.Adaptor](ctx context.Context, adaptor A, stmt string, named statements.Named) (results []R, err error) {

	var (
		tx           *sqlx.Tx
		statement    *sqlx.NamedStmt
		mode         dbs.Moder         = adaptor.Mode()
		transactions dbs.Transactioner = adaptor.TX()
		connector    dbs.Connector     = adaptor.Connector()
		dber         dbs.DBer          = adaptor.DB()
	)
	results = []R{}
	// validate the select statement passed
	err = statements.Validate(stmt, named)
	if err != nil {
		return
	}

	// connect to the database
	_, err = dber.Get(ctx, connector)
	if err != nil {
		return
	}
	defer dber.Close()

	// get a transaction
	tx, err = transactions.Get(ctx, dber, connector, mode)
	if err != nil {
		return
	}
	// prepare statement
	statement, err = tx.PrepareNamedContext(ctx, stmt)
	if err != nil {
		return
	}
	// if no arguments are passed use any empty struct
	if named == nil {
		named = &emptyNamed{}
	}
	// get the data
	err = statement.SelectContext(ctx, &results, named)
	if err != nil {
		return
	}
	//
	err = tx.Commit()
	return
}
