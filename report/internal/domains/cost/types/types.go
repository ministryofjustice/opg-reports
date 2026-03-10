package types

import (
	"opg-reports/report/packages/convert"
	"opg-reports/report/packages/types/interfaces"
	"opg-reports/report/packages/types/models"
	"strings"
)

// ImportCost is used in the importing of data from aws costexplorer api
//
// interfaces.Insertable
type ImportCost struct {
	Region    string `json:"region,omitempty"`      // AWS Region
	Service   string `json:"service,omitempty"`     // The AWS service name
	Month     string `json:"month,omitempty"`       // The data the cost was incurred - provided from the cost explorer result
	Cost      string `json:"cost,omitempty"`        // The actual cost value as a string - without an currency, but is USD by default
	AccountID string `json:"account_id,omityempty"` // the actual account id - string as it can have leading zeros. Use in joins as well
}

// Map returns a map of all fields on this struct
func (self *ImportCost) Map() (m map[string]interface{}) {
	m = map[string]interface{}{}
	convert.Between(self, &m)
	return
}

// CostByTeam used in the byteam api endpoint
//
// interfaces.Selectablw
type CostByTeam struct {
	Month string  `json:"month"`
	Cost  float64 `json:"cost"`
	Team  string  `json:"team"`
}

// Sequence is used to return the columns in the order they are selected
func (self *CostByTeam) Sequence() []any {
	return []any{
		&self.Month,
		&self.Cost,
		&self.Team,
	}
}

// Select is the select statement used for pull data from the db.
//
// interfaces.Statement
type Select struct {
	Statement string
}

// SQL updates the sql statement to adjust for team filtering
func (self *Select) SQL(filter interfaces.Filterable) (s string) {
	var model = &models.Filter{}
	var teamFilter = ""
	convert.Between(filter.Map(), &model)
	// if the team filter has been set, then update the sql
	if model.Team != "" {
		teamFilter = `accounts.team_name = :team AND`
	}
	s = strings.ReplaceAll(self.Statement, `{TEAM_FILTER}`, teamFilter)

	return
}
