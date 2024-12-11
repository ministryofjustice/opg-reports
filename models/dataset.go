package models

import (
	"strings"

	"github.com/ministryofjustice/opg-reports/internal/dbs"
)

// Dataset contains a single record to say if the data in use is real or not
//
// Interfaces:
//   - dbs.Table
//   - dbs.CreateableTable
//   - dbs.Insertable
//   - dbs.Row
//   - dbs.InsertableRow
//   - dbs.Record
//   - dbs.Cloneable
type Dataset struct {
	ID   int    `json:"id,omitempty" db:"id" faker:"-"`
	Ts   string `json:"ts,omitempty" db:"ts"  faker:"time_string" doc:"Time the record was created."` // TS is timestamp when the record was created
	Name string `json:"name,omitempty" db:"name" faker:"unique, oneof: real,fake"`
}

// UniqueValue returns the value representing the value of
// UniqueField
//
// Interfaces:
//   - dbs.Row
func (self *Dataset) UniqueValue() string {
	return strings.ToLower(self.Name)
}

// Interfaces:
//   - dbs.Insertable
func (self *Dataset) UniqueField() string {
	return "name"
}
func (self *Dataset) UpsertUpdate() string {
	return "name=excluded.name"
}

// TableName returns named table for Dataset - units
//
// Interfaces:
//   - dbs.Table
//   - dbs.CreateableTable
//   - dbs.Insertable
func (self *Dataset) TableName() string {
	return "dataset"
}

// Columns returns a map of all of the columns on the table - used for creation
//
// Interfaces:
//   - dbs.Createable
//   - dbs.CreateableTable
func (self *Dataset) Columns() map[string]string {
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
func (self *Dataset) Indexes() map[string][]string {
	return map[string][]string{
		"dataset_name_idx": {"name"},
	}
}

// InsetColumns returns the columns that should be used to insert a record into this table.
//
// Interfaces:
//   - dbs.Insertable
//   - dbs.InsertableRow
//   - dbs.Record
func (self *Dataset) InsertColumns() []string {
	return []string{
		"ts",
		"name",
	}
}

// GetID simply returns the current ID value for this row
//
// Interfaces:
//   - dbs.Row
//   - dbs.Insertable
//   - dbs.InsertableRow
//   - dbs.Record
func (self *Dataset) GetID() int {
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
func (self *Dataset) SetID(id int) {
	self.ID = id
}

// New is used by fakermany to return and empty instance of itself
// in an easier method
//
// Interfaces:
//   - dbs.Cloneable
func (self *Dataset) New() dbs.Cloneable {
	return &Dataset{}
}
