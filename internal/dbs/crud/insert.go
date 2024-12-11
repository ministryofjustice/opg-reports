package crud

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/structs"
)

var (
	errNotInWriteMode      error = fmt.Errorf("Adaptor is not write enabled.")
	errNoInsertableColumns error = fmt.Errorf("No columns returned for inserting.")
)

// Insert writes the record to the table and returns slice of the inserted records.
// Uses a prepared statement to run the write and creates the insert sql from the
// InsertableRow InsertColumns
// Returns on the first error
func Insert[A dbs.Adaptor, T dbs.Insertable, R dbs.InsertableRow](ctx context.Context, adaptor A, table T, records ...R) (inserted []R, err error) {

	var (
		tx           *sqlx.Tx
		statement    *sqlx.NamedStmt
		transactions dbs.Transactioner = adaptor.TX()
		mode         dbs.Moder         = adaptor.Mode()
		connector    dbs.Connector     = adaptor.Connector()
		dber         dbs.DBer          = adaptor.DB()
		tableName    string            = table.TableName()
		columns      []string          = table.InsertColumns()
		sqlStmt      string            = ""
	)
	inserted = []R{}
	// If its not in write mode, then return error
	if !mode.Write() {
		err = errNotInWriteMode
		return
	}
	// check columns exist
	if len(columns) <= 0 {
		err = errNoInsertableColumns
		return
	}

	fields := strings.TrimSuffix(strings.Join(columns, ","), ",")
	values := strings.TrimSuffix(strings.Join(columns, ",:"), ",:")
	str := `INSERT INTO %s (%s) VALUES (:%s)`

	if table.UniqueField() != "" {
		str += ` ON CONFLICT (` + table.UniqueField() + `) DO UPDATE SET ` + table.UpsertUpdate()
	}

	str += ` RETURNING id;`

	sqlStmt = fmt.Sprintf(str, tableName, fields, values)
	slog.Debug("[crud] insert", slog.String("sqlStmt", sqlStmt))

	// connect to the database
	_, err = dber.Get(ctx, connector)
	if err != nil {
		err = errors.Join(fmt.Errorf("sql: %s", sqlStmt), err)
		return
	}
	defer dber.Close()

	// get a transaction
	tx, err = transactions.Get(ctx, dber, connector, mode)
	if err != nil {
		return
	}

	// make the statement
	statement, err = tx.PrepareNamedContext(ctx, sqlStmt)
	if err != nil {
		return
	}

	// for all records generate the call
	for _, record := range records {
		var id int
		if err = statement.GetContext(ctx, &id, record); err == nil {
			record.SetID(id)
			inserted = append(inserted, record)
		} else if err != nil && err != sql.ErrNoRows {
			thisErr := fmt.Errorf("record: [%v]", structs.Jsonify(record))
			err = errors.Join(thisErr, err)
			return
		}
	}
	// commit
	err = transactions.Commit(true)
	if err != nil {
		transactions.Rollback()
	}

	return
}

// Truncate drops an existing table and then recreates it
func Truncate[A dbs.Adaptor, T dbs.Insertable](ctx context.Context, adaptor A, table T) (err error) {
	var (
		tx           *sqlx.Tx
		transactions dbs.Transactioner = adaptor.TX()
		mode         dbs.Moder         = adaptor.Mode()
		connector    dbs.Connector     = adaptor.Connector()
		dber         dbs.DBer          = adaptor.DB()
		sqlStmt      string            = ""
	)
	// If its not in write mode, then return error
	if !mode.Write() {
		err = errNotInWriteMode
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
	// delete all of the table
	sqlStmt = fmt.Sprintf(`DROP TABLE IF EXISTS %s`, table.TableName())
	// execute
	_, err = tx.ExecContext(ctx, sqlStmt)
	if err != nil {
		return
	}
	// commit
	err = transactions.Commit(true)
	if err != nil {
		transactions.Rollback()
	}

	_, err = CreateTable(ctx, adaptor, table)
	return
}
