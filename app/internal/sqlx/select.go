package sqlx

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"opg-reports/app/internal/logx"
)

// Result interface is used by Select to ensure the structs provide the
// BindResult method used to map the SQL result into the struct.
type Result interface {
	// BindResult returns a slice of pointers to the fields of this struct
	// in the correct order the SQL would return them.
	BindResult() []any
}

// RowScanner is a function used within Select to map the SQL rows into a
// slice of structs.
type RowScanner[T Result] func(rows *sql.Rows, results *[]T) (err error)

// Select intends to be a way to run a SELECT style sql statement against the database.
//
// Presumes the `sqlStmt` uses fieldname placeholders (`:fieldName`) to allow varied
// filter values to be set by the filter struct. Allows things like dynamic where
// values (such as `month IN (:months)`) to be used.
//
// The result slice should be passed by pointer to it can be updated within the
// row scanning function (`scanF`).
//
// `sqlStmt` and `filter` are based to `Bind` to generate a correctly structured
// string and argument slice which is then passed to `QueryContext`.
//
// The rows returned by `QueryContext` are then interated over with `rows.Next()`
// and then `scanF` is called within that loop to update the result slice.
//
// If scanF is nil, then `RowScan` is used by default.
//
// Note: Does not close the db connection.
func Select[T Result, F any](ctx context.Context, conn Connector, sqlStmt string, filter F, results *[]T, scanF RowScanner[T]) (err error) {
	var (
		db     *sql.DB
		sqlStr string
		rows   *sql.Rows
		args   []interface{} = []interface{}{}
		lg     *slog.Logger  = logx.Default()
	)
	// open the db connection
	db, err = conn.Open()
	if err != nil {
		lg.Error("failed to open database connection.", "err", err.Error())
		return
	}
	// create the sql
	lg.Debug("creating the sql & args from args.")
	sqlStr, args, err = Bind(ctx, sqlStmt, filter)
	if err != nil {
		lg.Error("failed to bind entry to sql.", "sqlStmt", sqlStmt, "err", err.Error())
		return
	}
	// run the query
	lg.Debug("run the query and get resulting rows.")
	rows, err = db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		lg.Error("error with the sql query.", "sqlStr", sqlStr, "err", err.Error())
		return
	}
	// provide a default row scanner if one wasnt set
	if scanF == nil {
		scanF = RowScan
	}

	defer rows.Close()

	lg.Debug("scanning rows from the result with scanner func.", "scanF", fmt.Sprintf("%T", scanF))
	for rows.Next() {
		err = scanF(rows, results)
		if err != nil {
			lg.Error("error while running row scan.", "err", err.Error())
			return
		}
	}

	lg.Debug("selection done.", "count", len(*results))
	return
}

// RowScan is the default method for attaching results from sql results to a
// slice of structs.
//
// Requires the use of BindResult on the struct, which returns a slice of
// pointers to the struct fields in the correct order the SQL would
// return them.
func RowScan[T Result](rows *sql.Rows, results *[]T) (err error) {
	var (
		r        T     = instanceOf[T]()
		sequence []any = r.BindResult()
	)

	err = rows.Scan(sequence...)
	if err == nil {
		*results = append(*results, r)
	}
	return
}
