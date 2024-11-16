package models

import (
	"github.com/ministryofjustice/opg-reports/internal/dbs"
)

// AwsCost represents a changed for service from AWS
// Interfaces:
//   - dbs.Table
//   - dbs.CreateableTable
//   - dbs.Insertable
//   - dbs.Row
//   - dbs.InsertableRow
//   - dbs.Record
//   - dbs.Cloneable
type AwsCost struct {
	ID           int                   `json:"id,omitempty" db:"id" faker:"-"`
	Ts           string                `json:"ts,omitempty" db:"ts"  faker:"time_string" doc:"Time the record was created."`                                                            // TS is timestamp when the record was created
	Region       string                `json:"region,omitempty" db:"region" faker:"oneof: NoRegion, eu-west-1, eu-west-2, us-east-2" doc:"Region this cost was generated within."`      // From the cost data, this is the region the service cost aws generated in
	Service      string                `json:"service,omitempty" db:"service" faker:"oneof: Tax, ecs, ec2, s3, sqs, waf, ses, rds" doc:"Name of the service that generated this cost."` // The AWS service name
	Date         string                `json:"date,omitempty" db:"date" faker:"date_string" doc:"Date this cost was generated."`                                                        // The data the cost was incurred - provided from the cost explorer result
	Cost         string                `json:"cost,omitempty" db:"cost" faker:"float_string" doc:"Cost value."`                                                                         // The actual cost value as a string - without an currency, but is USD by default
	AwsAccountID int                   `json:"aws_account_id,omitempty" db:"aws_account_id" faker:"-"`                                                                                  // AwsCost has many accounts, account has one unit
	AwsAccount   *AwsAccountForeignKey `json:"aws_account,omitempty" db:"aws_account" faker:"-"`                                                                                        // AwsCost has many accounts, account has one unit
	Unit         *UnitForeignKey       `json:"unit" db:"unit" faker:"-"`                                                                                                                // fetched from join on account
}

// TableName returns named table for AwsCost - units
//
// Interfaces:
//   - dbs.Table
//   - dbs.CreateableTable
//   - dbs.Insertable
func (self *AwsCost) TableName() string {
	return "aws_costs"
}

// Columns returns a map of all of the columns on the table - used for creation
//
// Interfaces:
//   - dbs.Createable
//   - dbs.CreateableTable
func (self *AwsCost) Columns() map[string]string {
	return map[string]string{
		"id":             "INTEGER PRIMARY KEY",
		"ts":             "TEXT NOT NULL",
		"region":         "TEXT NOT NULL",
		"service":        "TEXT NOT NULL",
		"date":           "TEXT NOT NULL",
		"cost":           "TEXT NOT NULL",
		"aws_account_id": "INTEGER NOT NULL",
	}
}

// Indexes returns a map contains the indexes to create on the this. This map should
// be formed with key being the name of the index and the []string containg the
// names of the columns to use.
//
// Interfaces:
//   - dbs.Createable
//   - dbs.CreateableTable
func (self *AwsCost) Indexes() map[string][]string {
	return map[string][]string{
		"awscosts_date_idx":         {"date"},
		"awscosts_date_account_idx": {"date", "aws_account_id"},
	}
}

// InsetColumns returns the columns that should be used to insert a record into this table.
//
// Interfaces:
//   - dbs.Insertable
//   - dbs.InsertableRow
//   - dbs.Record
func (self *AwsCost) InsertColumns() []string {
	return []string{
		"ts",
		"region",
		"service",
		"date",
		"cost",
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
func (self *AwsCost) GetID() int {
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
func (self *AwsCost) SetID(id int) {
	self.ID = id
}

// New is used by fakermany to return and empty instance of itself
// in an easier method
//
// Interfaces:
//   - dbs.Cloneable
func (self *AwsCost) New() dbs.Cloneable {
	return &AwsCost{}
}