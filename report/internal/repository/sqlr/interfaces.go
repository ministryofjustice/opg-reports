package sqlr

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Model interface{}

type Writer interface {
	Exec(statement string) (result sql.Result, err error)
	Insert(boundStatements ...*BoundStatement) (err error)
}
type Reader interface {
	Select(boundStatement *BoundStatement) (err error)
	ValidateSelect(boundStatement *BoundStatement) (valid bool, statement *sqlx.NamedStmt, err error)
}
