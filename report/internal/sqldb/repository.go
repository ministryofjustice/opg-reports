package sqldb

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/interfaces"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

type Repository[T interfaces.Model] struct {
	ctx  context.Context
	conf *config.Config
	log  *slog.Logger
}

// empty is used for selects that require a non-nil interface
// passed
type empty struct{}

// init is called via New and it creates database using the
// full SCHEMA
func (self *Repository[T]) init() (err error) {
	_, err = self.Exec(SCHEMA)
	return
}

// connection internal helper to handle connecting to the db
func (self *Repository[T]) connection() (db *sqlx.DB, err error) {
	var dbSource string = self.conf.Database.Source()

	// create the file path
	dir := filepath.Dir(self.conf.Database.Path)
	os.MkdirAll(dir, os.ModePerm)

	self.log.With("dbSource", dbSource).Debug("connecting to database...")
	db, err = sqlx.ConnectContext(self.ctx, self.conf.Database.Driver, dbSource)
	if err != nil {
		self.log.Error("connection failed", "error", err.Error(), "dbSource", dbSource)
	}

	return
}

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

// Insert creates a transaction for each bound statement and will fail if any
// insert fails.
// On fail a rollback is triggered
func (self *Repository[T]) Insert(boundStatements ...*BoundStatement) (err error) {
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
	}

	return
}

// Select uses the boundStatement to run command against the database
// and attach the result to a data item
func (self *Repository[T]) Select(boundStatement *BoundStatement) (err error) {
	var (
		db          *sqlx.DB
		transaction *sqlx.Tx
		statement   *sqlx.NamedStmt
		data                       = boundStatement.Data
		options     *sql.TxOptions = &sql.TxOptions{ReadOnly: false, Isolation: sql.LevelDefault}
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

// New creates a new repo
func New[T interfaces.Model](ctx context.Context, log *slog.Logger, conf *config.Config) (rp *Repository[T], err error) {
	if log == nil {
		return nil, fmt.Errorf("no logger passed for sqldb repository")
	}
	if conf == nil {
		return nil, fmt.Errorf("no config passed for sqldb repository")
	}
	log = log.WithGroup("sqldb")
	rp = &Repository[T]{
		ctx:  ctx,
		log:  log,
		conf: conf,
	}

	if !utils.FileExists(rp.conf.Database.Path) {
		log.With("database.path", rp.conf.Database.Path).Warn("Database not found, so creating")
		err = rp.init()
	}

	return
}
