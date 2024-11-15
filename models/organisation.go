package models

import (
	"github.com/ministryofjustice/opg-reports/internal/dbs"
)

// Organisation represents are organisational structure and can point to
// multiple git hub teams, services and accounts.
//
// Interfaces:
//   - dbs.Table
//   - dbs.CreateableTable
//   - dbs.Insertable
//   - dbs.Row
//   - dbs.InsertableRow
//   - dbs.Record
//   - dbs.Cloneable
type Organisation struct {
	ID    int    `json:"id,omitempty" db:"id" faker:"-"`
	Ts    string `json:"ts,omitempty" db:"ts"  faker:"time_string" doc:"Time the record was created."` // TS is timestamp when the record was created
	Name  string `json:"name,omitempty" db:"name" faker:"unique, oneof: opg,laa,moj,cica"`
	Units Units  `json:"units,omitempty" db:"units" faker:"-"`
}

// TableName returns named table for Organisation - units
//
// Interfaces:
//   - dbs.Table
//   - dbs.CreateableTable
//   - dbs.Insertable
func (self *Organisation) TableName() string {
	return "organisations"
}

// Columns returns a map of all of the columns on the table - used for creation
//
// Interfaces:
//   - dbs.Createable
//   - dbs.CreateableTable
func (self *Organisation) Columns() map[string]string {
	return map[string]string{
		"id":   "INTEGER PRIMARY KEY",
		"ts":   "TEXT NOT NULL",
		"name": "TEXT NOT NULL UNIQUE",
	}
}

// Indexes returns a map contains the indexes to create on the this. This map should
// be formed with key being the name of the index and the []string containg the
// names of the columns to use.
//
// Interfaces:
//   - dbs.Createable
//   - dbs.CreateableTable
func (self *Organisation) Indexes() map[string][]string {
	return map[string][]string{
		"org_name_idx": {"name"},
	}
}

// InsetColumns returns the columns that should be used to insert a record into this table.
//
// Interfaces:
//   - dbs.Insertable
//   - dbs.InsertableRow
//   - dbs.Record
func (self *Organisation) InsertColumns() []string {
	return []string{"ts", "name"}
}

// GetID simply returns the current ID value for this row
//
// Interfaces:
//   - dbs.Row
//   - dbs.Insertable
//   - dbs.InsertableRow
//   - dbs.Record
func (self *Organisation) GetID() int {
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
func (self *Organisation) SetID(id int) {
	self.ID = id
}

// New is used by fakermany to return and empty instance of itself
// in an easier method
//
// Interfaces:
//   - dbs.Cloneable
func (self *Organisation) New() dbs.Cloneable {
	return &Organisation{}
}
