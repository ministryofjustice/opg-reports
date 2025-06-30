package sqlr

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// Insert creates a transaction for each bound statement and will fail if any
// insert fails.
// On fail a rollback is triggered
func (self *Repository) Insert(boundStatements ...*BoundStatement) (err error) {
	var (
		db          *sqlx.DB
		transaction *sqlx.Tx
		options     *sql.TxOptions = &sql.TxOptions{ReadOnly: false, Isolation: sql.LevelDefault}
		log                        = self.log.With("operation", "select")
	)
	// db connection
	db, err = self.connection()
	if err != nil {
		return
	}
	defer db.Close()
	// start the transaction
	transaction, err = db.BeginTxx(self.ctx, options)
	if err != nil {
		return
	}
	// iterate over all the boundStatement and generate transactions
	for _, boundStmt := range boundStatements {
		var statement *sqlx.NamedStmt
		var data = boundStmt.Data

		statement, err = transaction.PrepareNamedContext(self.ctx, boundStmt.Statement)

		if err != nil {
			log.Error("prepared stmt failed", "error", err.Error())
			return
		}
		// data needs to be non-nil
		if data == nil {
			data = &empty{}
		}

		err = statement.GetContext(self.ctx, &boundStmt.Returned, data)
		if err != nil && err != sql.ErrNoRows {
			log.Error("stmt context failed", "error", err.Error(), "sql", statement.QueryString)
			return
		}

	}
	log.Debug("executing transaction...")
	err = transaction.Commit()
	// if theres an error on commit, rollback and return
	if err != nil {
		log.Error("transaction commit failed", "error", err.Error())
		transaction.Rollback()
		return
	}
	// now check all the bound statements executed
	err = checkSuccess(boundStatements...)

	return
}
