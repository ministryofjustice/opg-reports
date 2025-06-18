package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

type Repository struct {
	ctx  context.Context
	conf *config.Config
	log  *slog.Logger
}

// init is called via New and it creates database connection and then bootstraps
// the tables using schema within a transaction
func (self *Repository) init(initStatement string) (err error) {
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
	_, err = transaction.ExecContext(self.ctx, initStatement)
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

func (self *Repository) connection() (db *sqlx.DB, err error) {
	var dbSource string = self.conf.Database.Source()

	self.log.With("dbSource", dbSource).Info("connecting to database...")
	db, err = sqlx.ConnectContext(self.ctx, self.conf.Database.Driver, dbSource)
	if err != nil {
		self.log.Error("connection failed", "error", err.Error(), "dbSource", dbSource)
	}

	return
}

func (self *Repository) Insert(boundStatements ...*BoundStatement) (err error) {
	var (
		db          *sqlx.DB
		transaction *sqlx.Tx
		options     *sql.TxOptions = &sql.TxOptions{ReadOnly: false, Isolation: sql.LevelDefault}
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
		statement, err = transaction.PrepareNamedContext(self.ctx, boundStmt.Statement)
		if err != nil {
			self.log.Error("prepared stmt failed", "error", err.Error())
			return
		}
		err = statement.GetContext(self.ctx, &boundStmt.Returned, boundStmt.Data)
		if err != nil && err != sql.ErrNoRows {
			self.log.Error("stmt context failed", "error", err.Error())
			return
		}

	}
	self.log.Debug("executing transaction...")
	err = transaction.Commit()
	// if theres an error on commit, rollback and return
	if err != nil {
		self.log.Error("transaction commit failed", "error", err.Error())
		transaction.Rollback()
	}

	return
}

func New(ctx context.Context, log *slog.Logger, conf *config.Config) (rp *Repository, err error) {
	if log == nil {
		return nil, fmt.Errorf("no logger passed for owner service")
	}
	if conf == nil {
		return nil, fmt.Errorf("no config passed for owner service")
	}
	log = log.WithGroup("repository")
	rp = &Repository{
		ctx:  ctx,
		log:  log,
		conf: conf,
	}

	if !utils.FileExists(rp.conf.Database.Path) {
		log.With("dbPath", rp.conf.Database.Path).Warn("Database not found, so creating")
		err = rp.init(SCHEMA)
	}

	return
}
