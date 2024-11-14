package adaptors

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/internal/dbs"
)

// SqlxTransaction provides methods for getting an committing database transactions
// using sqlx
//
// Implements dbs.Transactioner
type SqlxTransaction struct {
	tx *sqlx.Tx
}

// GetTransaction will used the db porinter and the connection details passed to either
// create a new transaction or return a current version
func (self *SqlxTransaction) Get(ctx context.Context, dbp dbs.DBer, connection dbs.Connector, mode dbs.Moder) (tx *sqlx.Tx, err error) {
	var (
		db       *sqlx.DB
		readOnly bool           = !mode.Write()
		options  *sql.TxOptions = &sql.TxOptions{ReadOnly: readOnly, Isolation: sql.LevelDefault}
	)
	// grab the db
	db, err = dbp.Get(ctx, connection)
	if err != nil {
		return
	}

	if self.tx == nil {
		self.tx, err = db.BeginTxx(ctx, options)
	}
	tx = self.tx
	return
}

// CommitTransaction tries to commit the passed along transaction.
// If it failes and rollback is true it will attempt ad .RollBack() call
func (self *SqlxTransaction) Commit(rollback bool) (err error) {

	if self.tx != nil {
		err = self.tx.Commit()
		if err != nil && rollback {
			rollError := self.Rollback()
			err = errors.Join(err, rollError)
		}
		// reset the transaction to nil
		self.tx = nil
	}
	return
}

func (self *SqlxTransaction) Rollback() (err error) {

	if self.tx != nil {
		err = self.tx.Rollback()
		self.tx = nil
	}
	return
}
