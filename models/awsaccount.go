package models

import (
	"fmt"

	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/structs"
)

// AwsAccount
//
// Interfaces:
//   - dbs.Table
//   - dbs.CreateableTable
//   - dbs.Insertable
//   - dbs.Row
//   - dbs.InsertableRow
//   - dbs.Record
//   - dbs.Cloneable
type AwsAccount struct {
	ID          int    `json:"id,omitempty" db:"id" faker:"-"`
	Ts          string `json:"ts,omitempty" db:"ts" faker:"time_string" doc:"Time the record was created."`
	Number      string `json:"number,omitempty" db:"number" doc:"Account number"`
	Name        string `json:"name,omitempty" db:"name" faker:"unique, oneof: sirius prod, use prod, make prod, digideps prod, serve prod, refunds dev, sirius dev"`
	Label       string `json:"label,omitempty" db:"label" faker:"word" doc:"A supplimental lavel to provide extra detail on the account type."`       // Label is passed string that sets a more exact name - so DB account production
	Environment string `json:"environment,omitempty" db:"environment" faker:"oneof: production, pre-production, development" doc:"Environment type."` // Environment is passed along to show if this is production, development etc account

	// Unit join - one to many (accoutn can be for one unit only)
	UnitID int             `json:"unit_id,omitempty" db:"unit_id" faker:"-"`
	Unit   *UnitForeignKey `json:"unit,omitempty" db:"unit" faker:"-"`
}

// UniqueValue returns the value representing the value of
// UniqueField
//
// Interfaces:
//   - dbs.Row
func (self *AwsAccount) UniqueValue() string {
	return self.Number
}

// UniqueField for this model returns the name of the number field
//
// Interfaces:
//   - dbs.Insertable
func (self *AwsAccount) UniqueField() string {
	return "number"
}

func (self *AwsAccount) UpsertUpdate() string {
	return "unit_id=excluded.unit_id, environment=excluded.environment, name=excluded.name, label=excluded.label"
}

// TableName returns named table for AwsAccount - units
//
// Interfaces:
//   - dbs.Table
//   - dbs.CreateableTable
//   - dbs.Insertable
func (self *AwsAccount) TableName() string {
	return "aws_accounts"
}

// Columns returns a map of all of the columns on the table - used for creation
//
// Interfaces:
//   - dbs.Createable
//   - dbs.CreateableTable
func (self *AwsAccount) Columns() map[string]string {
	return map[string]string{
		"id":          "INTEGER PRIMARY KEY",
		"ts":          "TEXT NOT NULL",
		"number":      "TEXT NOT NULL UNIQUE",
		"name":        "TEXT NOT NULL",
		"label":       "TEXT NOT NULL",
		"environment": "TEXT NOT NULL",
		"unit_id":     "INTEGER",
	}
}

// Indexes returns a map contains the indexes to create on the this. This map should
// be formed with key being the name of the index and the []string containg the
// names of the columns to use.
//
// Interfaces:
//   - dbs.Createable
//   - dbs.CreateableTable
func (self *AwsAccount) Indexes() map[string][]string {
	return map[string][]string{
		"awsacc_unit_idx": {"unit_id"},
	}
}

// InsetColumns returns the columns that should be used to insert a record into this table.
//
// Interfaces:
//   - dbs.Insertable
//   - dbs.InsertableRow
//   - dbs.Record
func (self *AwsAccount) InsertColumns() []string {
	return []string{
		"ts",
		"number",
		"name",
		"label",
		"environment",
		"unit_id",
	}
}

// GetID simply returns the current ID value for this row
//
// Interfaces:
//   - dbs.Row
//   - dbs.Insertable
//   - dbs.InsertableRow
//   - dbs.Record
func (self *AwsAccount) GetID() int {
	return self.ID
}

// SetID allows setting the ID of this row - normally used within the insert calls
// to update the original data passed in with the new id
//
// Interfaces:
//   - dbs.Row
//   - dbs.Insertable
//   - dbs.InsertableRow
//   - dbs.Record
func (self *AwsAccount) SetID(id int) {
	self.ID = id
}

// New is used by fakermany to return and empty instance of itself
// in an easier method
//
// Interfaces:
//   - dbs.Cloneable
func (self *AwsAccount) New() dbs.Cloneable {
	return &AwsAccount{}
}

// AwsAccounts is to be used on the struct that needs to pull in
// the accounts via a many to many join select statement and provides
// the Scan method so sqlx will handle the result correctly
//
// Interfaces:
//   - sql.Scanner
type AwsAccounts []*AwsAccount

// Scan converts the json aggregate result from a select statement into
// a series of GitHubTeams attached to the main struct and will be called
// directly by sqlx
//
// Interfaces:
//   - sql.Scanner
func (self *AwsAccounts) Scan(src interface{}) (err error) {
	switch src.(type) {
	case []byte:
		err = structs.Unmarshal(src.([]byte), self)
	case string:
		err = structs.Unmarshal([]byte(src.(string)), self)
	default:
		err = fmt.Errorf("unsupported scan src type")
	}
	return
}

// AwsAccountForeignKey is to be used on the struct that needs to pull in
// the repo via one to many join (being used on the `one` side).
//
// To swap a AwsAccount to a AwsAccountForeignKey:
//
//	var join = models.AwsAccountForeignKey(&AwsAccount{})
//
// Interfaces:
//   - sql.Scanner
type AwsAccountForeignKey AwsAccount

// UniqueValue returns the value representing the value of
// UniqueField
//
// Interfaces:
//   - dbs.Row
func (self *AwsAccountForeignKey) UniqueValue() string {
	return self.Number
}

func (self *AwsAccountForeignKey) Scan(src interface{}) (err error) {
	switch src.(type) {
	case []byte:

		err = structs.Unmarshal(src.([]byte), self)
	case string:
		err = structs.Unmarshal([]byte(src.(string)), self)
	default:
		err = fmt.Errorf("unsupported scan src type")
	}
	return
}