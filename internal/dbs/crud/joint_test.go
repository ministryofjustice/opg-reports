package crud_test

import "github.com/ministryofjustice/opg-reports/internal/dbs"

var _ dbs.Joiner = TeamList{}

type TeamList []*Team

func (self TeamList) Scan(src interface{}) (err error) {
	return nil
}

type Unit struct {
	ID   int    `json:"id,omitempty" db:"id" faker:"-"`
	Name string `json:"name,omitempty" db:"name" faker:"word"`
}

func (self *Unit) TableName() string {
	return "units"
}

func (self *Unit) GetID() int {
	return self.ID
}
func (self *Unit) SetID(id int) {
	self.ID = id
}
func (self *Unit) Columns() map[string]string {
	return map[string]string{"id": "INTEGER PRIMARY KEY", "name": "TEXT NOT NULL"}
}

func (self *Unit) Indexes() map[string][]string {
	return map[string][]string{
		"idx_name": {"name"},
	}
}
func (self *Unit) InsertColumns() []string {
	return []string{
		"name", "ts",
	}
}
func (self *Unit) New() dbs.Cloneable {
	return &Unit{}
}

type Team struct {
	ID   int    `json:"id,omitempty" db:"id" faker:"-"`
	Name string `json:"name,omitempty" db:"name" faker:"word"`
}

func (self *Team) TableName() string {
	return "teams"
}

func (self *Team) GetID() int {
	return self.ID
}
func (self *Team) SetID(id int) {
	self.ID = id
}
func (self *Team) Columns() map[string]string {
	return map[string]string{"id": "INTEGER PRIMARY KEY", "name": "TEXT NOT NULL"}
}
func (self *Team) Indexes() map[string][]string {
	return map[string][]string{
		"idx_name": {"name"},
	}
}
func (self *Team) InsertColumns() []string {
	return []string{
		"name", "ts",
	}
}
func (self *Team) New() dbs.Cloneable {
	return &Team{}
}
