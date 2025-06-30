package sqlr

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// Select uses the boundStatement to run command against the database
// and attach the result to a data item
func (self *Repository[T]) Select(boundStatement *BoundStatement) (err error) {
	var (
		db          *sqlx.DB
		transaction *sqlx.Tx
		statement   *sqlx.NamedStmt
		data                       = boundStatement.Data
		options     *sql.TxOptions = &sql.TxOptions{ReadOnly: true, Isolation: sql.LevelDefault}
		log                        = self.log.With("operation", "select")
		returned                   = []T{}
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

	statement, err = transaction.PrepareNamedContext(self.ctx, boundStatement.Statement)
	if err != nil {
		log.Error("prepared stmt failed", "error", err.Error())
		return
	}
	// data needs to be non-nil
	if data == nil {
		data = &empty{}
	}

	err = statement.SelectContext(self.ctx, &returned, data)
	if err != nil && err != sql.ErrNoRows {
		log.Error("stmt context failed", "error", err.Error())
		return
	}
	boundStatement.Returned = returned

	log.Debug("executing transaction...")
	err = transaction.Commit()

	return
}

// ValidateSelect makes a short connection the database and tries to prepare the provided statement
// in order to validate it without running it - allows a way to test sql ahead of running it
func (self *Repository[T]) ValidateSelect(boundStatement *BoundStatement) (valid bool, statement *sqlx.NamedStmt, err error) {
	var db *sqlx.DB
	valid = false
	// db connection
	db, err = self.connection()
	if err != nil {
		return
	}
	defer db.Close()

	statement, err = db.PrepareNamedContext(self.ctx, boundStatement.Statement)
	if err == nil && statement.QueryString != "" {
		valid = true
	}
	return
}

// checkSuccess looks at each statments to see if there are any without a
// Returned value, if so it flags an error
func checkSuccess(boundStatements ...*BoundStatement) (err error) {
	// check if everything was inserted
	for _, stmt := range boundStatements {
		if stmt.Returned == nil {
			return fmt.Errorf("some statments failed to complete")
		}
	}
	return nil
}
