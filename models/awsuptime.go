package models

import (
	"github.com/ministryofjustice/opg-reports/internal/dbs"
)

// AwsUptime
//
// Interfaces:
//   - dbs.Table
//   - dbs.CreateableTable
//   - dbs.Insertable
//   - dbs.Row
//   - dbs.InsertableRow
//   - dbs.Record
//   - dbs.Cloneable
type AwsUptime struct {
	ID           int                   `json:"id,omitempty" db:"id" faker:"-"`
	Ts           string                `json:"ts,omitempty" db:"ts" faker:"time_string" doc:"Time the record was created."`
	Date         string                `json:"date,omitempty" db:"date" faker:"date_string" doc:"Date this cost was generated."` // The interval date for when this uptime was logged
	Average      float64               `json:"average,omitempty" db:"average" doc:"Percentage uptime average for this record period."`
	UnitID       int                   `json:"unit_id" db:"unit_id" faker:"-"`
	AwsAccountID int                   `json:"aws_account_id" db:"aws_account_id" faker:"-"`
	Unit         *UnitForeignKey       `json:"unit" db:"unit" faker:"-"`
	AwsAccount   *AwsAccountForeignKey `json:"aws_account" db:"aws_account" faker:"-"`
}

// TableName returns named table for AwsUptime - units
//
// Interfaces:
//   - dbs.Table
//   - dbs.CreateableTable
//   - dbs.Insertable
func (self *AwsUptime) TableName() string {
	return "aws_uptime"
}

// Columns returns a map of all of the columns on the table - used for creation
//
// Interfaces:
//   - dbs.Createable
//   - dbs.CreateableTable
func (self *AwsUptime) Columns() map[string]string {
	return map[string]string{
		"id":             "INTEGER PRIMARY KEY",
		"ts":             "TEXT NOT NULL",
		"date":           "TEXT NOT NULL",
		"average":        "REAL NOT NULL",
		"unit_id":        "INTEGER",
		"aws_account_id": "INTEGER",
	}
}

// Indexes returns a map contains the indexes to create on the this. This map should
// be formed with key being the name of the index and the []string containg the
// names of the columns to use.
//
// Interfaces:
//   - dbs.Createable
//   - dbs.CreateableTable
func (self *AwsUptime) Indexes() map[string][]string {
	return map[string][]string{
		"awsup_unit_idx": {"unit_id"},
		"awsup_acc_idx":  {"aws_account_id"},
	}
}

// InsetColumns returns the columns that should be used to insert a record into this table.
//
// Interfaces:
//   - dbs.Insertable
//   - dbs.InsertableRow
//   - dbs.Record
func (self *AwsUptime) InsertColumns() []string {
	return []string{
		"ts",
		"date",
		"average",
		"unit_id",
		"aws_account_id",
	}
}

// GetID simply returns the current ID value for this row
//
// Interfaces:
//   - dbs.Row
//   - dbs.Insertable
//   - dbs.InsertableRow
//   - dbs.Record
func (self *AwsUptime) GetID() int {
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
func (self *AwsUptime) SetID(id int) {
	self.ID = id
}

// New is used by fakermany to return and empty instance of itself
// in an easier method
//
// Interfaces:
//   - dbs.Cloneable
func (self *AwsUptime) New() dbs.Cloneable {
	return &AwsUptime{}
}
