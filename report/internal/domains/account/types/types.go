package types

import (
	"opg-reports/report/packages/convert"
	"opg-reports/report/packages/types/interfaces"
	"opg-reports/report/packages/types/models"
	"strings"
)

// ImportAccount is used in the importing of data from
// raw account files.
//
// Has a different set of json field names to align with
// the original source files.
//
// interfaces.Insertable
type ImportAccount struct {
	ID          string `json:"id,omitempty"`            // This is the Account ID as a string - they can have leading 0
	Name        string `json:"name,omitempty" `         // account name as used internally
	Label       string `json:"label,omitempty" `        // internal label
	Environment string `json:"environment,omitempty" `  // environment type
	TeamName    string `json:"billing_unit,omitempty" ` // team associated with the account; uses builling_unit due to the source data in opg-metadata
}

// Map returns a map of all fields on this struct
func (self *ImportAccount) Map() (m map[string]interface{}) {
	m = map[string]interface{}{}
	convert.Between(self, &m)
	return
}

// Account is used in the api result and contain
// the result of the database select query
//
// interfaces.Selectable
// interfaces.Result
type Account struct {
	ID          string `json:"id"`
	Name        string `json:"name" `
	Label       string `json:"label" `
	Environment string `json:"environment" `
	TeamName    string `json:"team"`
}

// Sequence is used to return the columns in the order they are selected.
func (self *Account) Sequence() []any {
	return []any{
		&self.ID,
		&self.Name,
		&self.Label,
		&self.Environment,
		&self.TeamName,
	}
}
func (self *Account) Result() interfaces.Result {
	var m = map[string]interface{}{}
	convert.Between(self, &m)
	return m
}

// Select is the select statement used for pull data from the db.
//
// interfaces.Statement
type Select struct {
	Statement string
}

// SQL updates the sql statement to adjust for team filtering
func (self *Select) SQL(filter interfaces.Filterable) (s string) {
	var model = filter.(*models.Filter)
	var teamFilter = ""
	// if the team filter has been set, then update the sql
	if model.Team != "" {
		teamFilter = `team_name = :team AND`
	}
	s = strings.ReplaceAll(self.Statement, `{TEAM_FILTER}`, teamFilter)

	return
}
