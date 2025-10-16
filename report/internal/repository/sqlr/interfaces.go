package sqlr

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Model interface{}

type RepositoryIdentifer interface {
	ID() (string, error)
}

type RepositoryPinger interface {
	Ping() (err error)
}

type RepositoryWriter interface {
	RepositoryIdentifer
	Exec(statement string) (result sql.Result, err error)
	Insert(boundStatements ...*BoundStatement) (err error)
}
type RepositoryReader interface {
	RepositoryIdentifer
	Select(boundStatement *BoundStatement) (err error)
	ValidateSelect(boundStatement *BoundStatement) (valid bool, statement *sqlx.NamedStmt, err error)
}
