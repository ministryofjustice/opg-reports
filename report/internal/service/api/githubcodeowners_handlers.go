package api

import (
	"opg-reports/report/internal/repository/sqlr"
)

// GithubCodeOwnerGetter interface
type GithubCodeOwnersGetter[T Model] interface {
	Closer
	GetAllGithubCodeOwners(store sqlr.RepositoryReader) (data []T, err error)
}

// GithubCodeOwnersForTeamGetter interface
type GithubCodeOwnersForTeamGetter[T Model] interface {
	Closer
	GetAllGithubCodeOwnersForTeam(store sqlr.RepositoryReader, options *GetAllGithubCodeOwnersForTeamOptions) (data []T, err error)
}

// GithubCodeOwnersForCodeOwnerGetter interface
type GithubCodeOwnersForCodeOwnerGetter[T Model] interface {
	Closer
	GetAllGithubCodeOwnersForCodeOwner(store sqlr.RepositoryReader, options *GetAllGithubCodeOwnersForCodeOwnerOptions) (data []T, err error)
}

// GithubCodeOwner
type GithubCodeOwner struct {
	CodeOwner  string `json:"codeowner,omitempty" db:"codeowner"`
	Repository string `json:"repository" db:"repository"`
	Team       string `json:"team,omitempty" db:"team"`
}

type GetAllGithubCodeOwnersForTeamOptions struct {
	Team string `json:"team" db:"team"`
}
type GetAllGithubCodeOwnersForCodeOwnerOptions struct {
	CodeOwner string `json:"codeowner" db:"codeowner"`
}

// GetAllGithubCodeOwnersForTeam
func (self *Service[T]) GetAllGithubCodeOwnersForTeam(store sqlr.RepositoryReader, options *GetAllGithubCodeOwnersForTeamOptions) (data []T, err error) {
	var statement = &sqlr.BoundStatement{Statement: stmtGithubCodeOwnerSelectForTeam, Data: options}
	var log = self.log.With("operation", "GetAllGithubCodeOwnersForTeam")

	data = []T{}
	log.Debug("getting all github codeowners from database for team ...")

	if err = store.Select(statement); err == nil {
		// cast the data back to struct
		data = statement.Returned.([]T)
	}

	return
}

// GetAllGithubCodeOwnersForCodeOwner
func (self *Service[T]) GetAllGithubCodeOwnersForCodeOwner(store sqlr.RepositoryReader, options *GetAllGithubCodeOwnersForCodeOwnerOptions) (data []T, err error) {
	var statement = &sqlr.BoundStatement{Statement: stmtGithubCodeOwnerSelectForCodeOwner, Data: options}
	var log = self.log.With("operation", "GetAllGithubCodeOwnersForCodeOwner")

	data = []T{}
	log.Debug("getting all github codeowners from database for codeowner ...")

	if err = store.Select(statement); err == nil {
		// cast the data back to struct
		data = statement.Returned.([]T)
	}

	return
}

// GetAllGithubCodeOwners returns all teams and joins aws accounts as well
func (self *Service[T]) GetAllGithubCodeOwners(store sqlr.RepositoryReader) (data []T, err error) {
	var statement = &sqlr.BoundStatement{Statement: stmtGithubCodeOwnerSelectAll}
	var log = self.log.With("operation", "GetAllGithubCodeOwners")

	data = []T{}
	log.Debug("getting all github codeowners from database ...")

	if err = store.Select(statement); err == nil {
		// cast the data back to struct
		data = statement.Returned.([]T)
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
