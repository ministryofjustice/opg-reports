package api

import (
	"fmt"

	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/utils"
)

// AwsAccountsGetter interface is used for GetAllAccounts calls
type AwsAccountsGetter[T Model] interface {
	Closer
	GetAllAwsAccounts(store sqlr.Reader) (accounts []T, err error)
}

const stmtAwsAccountsSelectAll string = `
SELECT
	aws_accounts.id,
	aws_accounts.name,
	aws_accounts.label,
	aws_accounts.environment,
	json_object(
		'name', aws_accounts.team_name
	) as team
FROM aws_accounts
GROUP BY aws_accounts.id
ORDER BY aws_accounts.team_name ASC, aws_accounts.name ASC, aws_accounts.environment ASC;`

// AwsAccount is api response model
type AwsAccount struct {
	ID          string `json:"id,omitempty" db:"id" example:"012345678910"` // This is the AWS Account ID as a string
	Name        string `json:"name,omitempty" db:"name" example:"Public API"`
	Label       string `json:"label,omitempty" db:"label" example:"aurora-cluster"`
	Environment string `json:"environment,omitempty" db:"environment" example:"development|preproduction|production"`
	// Joins to team
	Team *hasOneTeam `json:"team,omitempty" db:"team"`
}

// team is internal and used for handling the account->team join
type accountTeam struct {
	Name string `json:"name,omitempty" db:"name" example:"SRE"`
}

type hasOneTeam accountTeam

// Scan handles the processing of the join data
func (self *hasOneTeam) Scan(src interface{}) (err error) {
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

// GetAllAccounts returns all accounts as a slice from the database
func (self *Service[T]) GetAllAwsAccounts(store sqlr.Reader) (accounts []T, err error) {
	var selectStmt = &sqlr.BoundStatement{Statement: stmtAwsAccountsSelectAll}
	var log = self.log.With("operation", "GetAllAccounts")

	accounts = []T{}
	log.Debug("getting all awsaccounts from database ...")

	// cast the data back to struct
	if err = store.Select(selectStmt); err == nil {
		accounts = selectStmt.Returned.([]T)
	}

	return
}
