package awsaccount

import (
	"fmt"

	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

type AwsAccount struct {
	ID          string `json:"id,omitempty" db:"id" example:"012345678910"` // This is the AWS Account ID as a string
	Name        string `json:"name,omitempty" db:"name" example:"Public API"`
	Label       string `json:"label,omitempty" db:"label" example:"aurora-cluster"`
	Environment string `json:"environment,omitempty" db:"environment" example:"development|preproduction|production"`
	CreatedAt   string `json:"created_at,omitempty" db:"created_at" example:"2019-08-24T14:15:22Z"`

	// Joins
	// TeamID is the raw db column that stores the association, should not
	TeamID int         `json:"team_id,omitempty" db:"team_id"`
	Team   *hasOneTeam `json:"team,omitempty" db:"team"`
}

// team is internal and used for handling the account->team join
type team struct {
	Name string `json:"name,omitempty" db:"name" example:"SRE"`
}

type hasOneTeam team

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

// AwsAccountImport captures an extra field from the metadata which
// is used in the stmtInsert to create the initial join to team based
// on the billing_unit name
type AwsAccountImport struct {
	AwsAccount
	BillingUnit string `json:"billing_unit,omitempty" db:"billing_unit"`
}

func newImportAccount(account *AwsAccount, billingUnit string) (acc *AwsAccountImport) {
	acc = &AwsAccountImport{}
	utils.Convert(account, &acc)
	acc.BillingUnit = billingUnit
	return
}
