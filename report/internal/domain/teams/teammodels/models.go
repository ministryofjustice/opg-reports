package teammodels

import (
	"fmt"
	"opg-reports/report/internal/utils/unmarshal"
)

type Team struct {
	Name string `json:"name" db:"name"`
}

type TeamData struct {
	Name     string          `json:"name" db:"team_name"`
	Accounts hasManyAccounts `json:"accounts,omitempty" db:"account_list"` // list from the join
}

type joinedAccount struct {
	ID          string `json:"id" db:"id"`                   // This is the Account ID as a string - they can have leading 0
	Name        string `json:"name" db:"name"`               // account name as used internally
	Label       string `json:"label" db:"label"`             // internal label
	Environment string `json:"environment" db:"environment"` // environment type
}

// hasManyAccounts represents the join
// Interfaces:
//   - sql.Scanner
type hasManyAccounts []*joinedAccount

// Scan handles the processing of the join data
func (self *hasManyAccounts) Scan(src interface{}) (err error) {
	switch src.(type) {
	case []byte:
		err = unmarshal.Unmarshal(src.([]byte), self)
	case string:
		err = unmarshal.Unmarshal([]byte(src.(string)), self)
	default:
		err = fmt.Errorf("unsupported scan src type")
	}
	return
}
