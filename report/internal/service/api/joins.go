package api

import (
	"fmt"
	"opg-reports/report/internal/utils"
)

// joinedAwsAccount is used to capture sql join data for accounts
type joinedAwsAccount struct {
	ID          string `json:"id,omitempty" db:"id" example:"012345678910"` // This is the AWS Account ID as a string
	Name        string `json:"name,omitempty" db:"name" example:"Public API"`
	Label       string `json:"label,omitempty" db:"label" example:"aurora-cluster"`
	Environment string `json:"environment,omitempty" db:"environment" example:"development|preproduction|production"`
}

// hasOneAwsAccount used for the join to awsAccount to handle the scaning into a struct
type hasOneAwsAccount joinedAwsAccount

// Scan handles the processing of the join data
func (self *hasOneAwsAccount) Scan(src interface{}) (err error) {
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

// joinedTeam maps to the team model in sequel
type joinedTeam struct {
	Name string `json:"name,omitempty" db:"name" example:"SRE"`
}

type hasOneTeam joinedTeam

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
