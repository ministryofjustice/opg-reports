package api

import (
	"opg-reports/report/internal/repository/sqlr"
)

// GithubCodeOwnerGetter interface is used for GetAllTeams calls
type GithubCodeOwnersGetter[T Model] interface {
	Closer
	GetAllGithubCodeOwners(store sqlr.RepositoryReader) (teams []T, err error)
}

// GithubCodeOwner
type GithubCodeOwner struct {
	CodeOwner  string `json:"codeowner,omitempty" db:"codeowner"`
	Repository string `json:"repository" db:"repository"`
	Team       string `json:"team,omitempty" db:"team"`
}

// GetAllGithubCodeOwners returns all teams and joins aws accounts as well
func (self *Service[T]) GetAllGithubCodeOwners(store sqlr.RepositoryReader) (teams []T, err error) {
	var statement = &sqlr.BoundStatement{Statement: stmtGithubCodeOwnerSelectAll}
	var log = self.log.With("operation", "GetAllGithubCodeOwners")

	teams = []T{}
	log.Debug("getting all github codeowners from database ...")

	if err = store.Select(statement); err == nil {
		// cast the data back to struct
		teams = statement.Returned.([]T)
	}

	return
}

// TruncateAndPutGithubCodeOwners inserts new records into the table.
//
// WARNING: This will truncate the DB table
// Note: Dont expose to the api endpoints
func (self *Service[T]) TruncateAndPutGithubCodeOwners(store sqlr.RepositoryWriter, data []T) (results []*sqlr.BoundStatement, err error) {

	_, err = store.Exec(stmtGithubCodeOwnerTruncate)
	if err != nil {
		return
	}
	return self.Put(store, stmtGithubCodeOwnerInsert, data)
}
