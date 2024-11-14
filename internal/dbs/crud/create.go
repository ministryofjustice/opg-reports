package crud

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/internal/dbs"
)

var (
	errCreateTableNoColumns error = fmt.Errorf("Columns() returned empty.")
	errCreateIndexNoIndexes error = fmt.Errorf("Indexes() returned empty.")
	errCreateIndexNoColumns error = fmt.Errorf("Indexes() contained empty index columns.")
)

// CreateTable uses the adaptor and table passed in to connect to a sql database and create
// a new table using the table name and column details from `table`
// Uses transactions within this function
func CreateTable[A dbs.Adaptor, T dbs.CreateableTable](ctx context.Context, adaptor A, table T) (result sql.Result, err error) {
	var (
		tx           *sqlx.Tx
		transactions dbs.Transactioner = adaptor.TX()
		mode         dbs.Moder         = adaptor.Mode()
		connector    dbs.Connector     = adaptor.Connector()
		dber         dbs.DBer          = adaptor.DB()
		tableName    string            = table.TableName()
		columns      map[string]string = table.Columns()
		rollBack     bool              = true
		sqlStmt      string            = ""
	)

	// if there are no columns, return error
	if len(columns) <= 0 {
		err = errCreateTableNoColumns
		return
	}

	// Create the table create sql statement
	sqlStmt = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (`, tableName)
	for field, definition := range columns {
		var column string = fmt.Sprintf("%s %s,", field, definition)
		sqlStmt += column
	}
	sqlStmt = strings.TrimSuffix(sqlStmt, ",")
	sqlStmt += `) STRICT;`
	slog.Debug("[crud] createtable", slog.String("sqlStmt", sqlStmt))
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
	// execute the statement
	if result, err = tx.ExecContext(ctx, sqlStmt); err == nil {
		err = transactions.Commit(rollBack)
	}

	return
}

// CreateIndexes uses the sql adaptor and table structure to generate a series of indexes to create
// against itself
// Uses transactions within this function
func CreateIndexes[A dbs.Adaptor, T dbs.CreateableTable](ctx context.Context, adaptor A, table T) (results []sql.Result, err error) {

	var (
		tx           *sqlx.Tx
		transactions dbs.Transactioner   = adaptor.TX()
		mode         dbs.Moder           = adaptor.Mode()
		connector    dbs.Connector       = adaptor.Connector()
		dber         dbs.DBer            = adaptor.DB()
		rollBack     bool                = true
		indexes      map[string][]string = table.Indexes()
		tableName    string              = table.TableName()
	)

	// if there are no columns, return error
	if len(indexes) <= 0 {
		err = errCreateIndexNoIndexes
		return
	}

	// connect to the database
	_, err = dber.Get(ctx, connector)
	if err != nil {
		return
	}
	defer dber.Close()
	// get a transaction started
	tx, err = transactions.Get(ctx, dber, connector, mode)
	if err != nil {
		return
	}

	for idxName, columns := range indexes {
		// check there are columns to use for the index
		if len(columns) <= 0 {
			err = errCreateTableNoColumns
			return
		}
		var result sql.Result
		var cols = strings.TrimSuffix(strings.Join(columns, ","), ",")
		var stmt = fmt.Sprintf(`CREATE INDEX IF NOT EXISTS %s on %s(%s);`, idxName, tableName, cols)

		slog.Debug("[crud] createindex", slog.String("sqlStmt", stmt))

		result, err = tx.ExecContext(ctx, stmt)
		results = append(results, result)
		if err != nil {
			return
		}
	}

	err = transactions.Commit(rollBack)

	return
}
