package sqlr

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Model interface{}

type RepositoryWriter interface {
	Exec(statement string) (result sql.Result, err error)
	Insert(boundStatements ...*BoundStatement) (err error)
}
type RepositoryReader interface {
	Select(boundStatement *BoundStatement) (err error)
	ValidateSelect(boundStatement *BoundStatement) (valid bool, statement *sqlx.NamedStmt, err error)
}
