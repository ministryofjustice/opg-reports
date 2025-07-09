package api

import (
	"fmt"

	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/utils"
)

// TeamGetter interface is used for GetAllTeams calls
type TeamGetter[T Model] interface {
	Closer
	GetAllTeams(store sqlr.Reader) (teams []T, err error)
}

// stmtTeamsSelectAll is sql used to fetch all teams and the join to aws accounts
const stmtTeamsSelectAll string = `
SELECT
	teams.name,
	json_group_array(
		DISTINCT json_object(
			'id', aws_accounts.id,
			'name', aws_accounts.name,
			'label', aws_accounts.label,
			'environment', aws_accounts.environment
		)
	) filter ( where aws_accounts.id is not null) as aws_accounts
FROM teams
LEFT JOIN aws_accounts ON aws_accounts.team_name = teams.name
GROUP BY teams.name
ORDER BY teams.name ASC;`

// Team acts as a top level group matching accounts, github repos etc to be attached
// as owners
type Team struct {
	Name string `json:"name,omitempty" db:"name" example:"SREs"`
	// Joins
	// AwsAccounts is the one->many join from teams to awsaccounts
	AwsAccounts hasManyAwsAccounts `json:"aws_accounts,omitempty" db:"aws_accounts"`
}

// teamAwsAccount is internal and used for the join on the team model only.
//
// Matches with team.Team fields but copied over to only references values in
// the sql statements and to avoide circular references.
//
// Note: struct name is unique due to how huma generates schema.
type teamAwsAccount struct {
	ID          string `json:"id,omitempty" db:"id" example:"012345678910"`
	Name        string `json:"name,omitempty" db:"name" example:"Public API"`
	Label       string `json:"label,omitempty" db:"label" example:"production database"`
	Environment string `json:"environment,omitempty" db:"environment" example:"development|preproduction|production"`
}

// hasManyAwsAccounts is used for one->many join from Team to AwsAccount
//
// Interfaces:
//   - sql.Scanner
type hasManyAwsAccounts []*teamAwsAccount

// Scan handles the processing of the join data
func (self *hasManyAwsAccounts) Scan(src interface{}) (err error) {
	switch src.(type) {
	case []byte:
		err = utils.Unmarshal(src.([]byte), self)
	case string:
		err = utils.Unmarshal([]byte(src.(string)), self)
	default:
		err = fmt.Errorf("unsupported scan src type")
	}
	return
}

// GetAllTeams returns all teams and joins aws accounts as well
func (self *Service[T]) GetAllTeams(store sqlr.Reader) (teams []T, err error) {
	var statement = &sqlr.BoundStatement{Statement: stmtTeamsSelectAll}
	var log = self.log.With("operation", "GetAllTeams")

	teams = []T{}
	log.Debug("getting all teams from database ...")

	if err = store.Select(statement); err == nil {
		// cast the data back to struct
		teams = statement.Returned.([]T)
	}

	return
}
