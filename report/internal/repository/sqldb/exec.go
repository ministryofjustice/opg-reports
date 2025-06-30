package sqldb

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// Exec runs a complete statement against the database and returns any error
// Used for mostly calls without parameters (like create / delete) that either
// return no result or simple value
func (self *Repository[T]) Exec(statement string) (result sql.Result, err error) {
	var (
		db          *sqlx.DB
		transaction *sqlx.Tx
		options     *sql.TxOptions = &sql.TxOptions{ReadOnly: false, Isolation: sql.LevelDefault}
	)
	db, err = self.connection()
	if err != nil {
		return
	}
	defer db.Close()
	// start the transaction
	transaction, err = db.BeginTxx(self.ctx, options)
	// try to execute all schema
	result, err = transaction.ExecContext(self.ctx, statement)
	if err != nil {
		self.log.Error("exec failed", "error", err.Error())
		return
	}
	// if no error, commit the transaction
	self.log.Debug("executing transaction...")
	err = transaction.Commit()
	// if theres an error on commit, rollback and return
	if err != nil {
		self.log.Error("transaction commit failed", "error", err.Error())
		transaction.Rollback()
	}
	return
}
