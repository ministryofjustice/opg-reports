package sqldb

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type SQLer interface {
	Exec(statement string) (result sql.Result, err error)
	Insert(boundStatements ...*BoundStatement) (err error)
	Select(boundStatement *BoundStatement) (err error)
	ValidateSelect(boundStatement *BoundStatement) (valid bool, statement *sqlx.NamedStmt, err error)
}
