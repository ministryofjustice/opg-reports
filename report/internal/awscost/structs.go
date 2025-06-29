package awscost

import (
	"fmt"
	"strings"

	"github.com/ministryofjustice/opg-reports/report/internal/sqldb"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

// GetGroupedCostsOptions contains a series of values that determines
// what fields are used within the sql statement to allow for easier
// handling of multiple, similar sql queries that differ by which
// columns are grouped or filtered
type GetGroupedCostsOptions struct {
	StartDate  string
	EndDate    string
	DateFormat string

	Team        utils.TrueOrFilter
	Region      utils.TrueOrFilter
	Service     utils.TrueOrFilter
	Account     utils.TrueOrFilter
	Environment utils.TrueOrFilter
}

type sqlParams struct {
	StartDate   string `json:"start_date,omitempty" db:"start_date"`
	EndDate     string `json:"end_date,omitempty" db:"end_date"`
	DateFormat  string `json:"date_format,omitempty" db:"date_format"`
	Region      string `json:"region,omitempty" db:"region"`
	Service     string `json:"service,omitempty" db:"service"`
	Team        string `json:"team_name,omitempty" db:"team_name"`
	Account     string `json:"aws_account_id,omitempty" db:"aws_account_id"`
	Environment string `json:"environment,omitempty" db:"environment"`
}

// Statement converts the configured options to a bound statement and provides the
// values and :params for `stmGroupedCosts`.
//
// It returns the bound statement and generated data object
func (self *GetGroupedCostsOptions) Statement() (bound *sqldb.BoundStatement, params *sqlParams) {
	var (
		stmt            = stmGroupedCosts
		selected string = ""
		where    string = ""
		orderby  string = ""
		groupby  string = ""
	)
	// setup the default data values
	params = &sqlParams{
		StartDate:  self.StartDate,
		EndDate:    self.EndDate,
		DateFormat: self.DateFormat,
	}

	// check the team, account, env, region and service values and update the
	// sql

	// Team
	if self.Team.Selectable() {
		selected += fmt.Sprintf("%s,", "team_name")
	}
	if self.Team.Whereable() {
		params.Team = string(self.Team)
		where += fmt.Sprintf("%s AND ", "team_name=:team_name")
	}
	if self.Team.Groupable() {
		groupby += fmt.Sprintf("%s,", "team_name")
	}
	if self.Team.Orderable() {
		orderby += fmt.Sprintf("%s ASC,", "team_name")
	}

	// Region
	if self.Region.Selectable() {
		selected += fmt.Sprintf("%s,", "region")
	}
	if self.Region.Whereable() {
		params.Region = string(self.Region)
		where += fmt.Sprintf("%s AND ", "region=:region")
	}
	if self.Region.Groupable() {
		groupby += fmt.Sprintf("%s,", "region")
	}
	if self.Region.Orderable() {
		orderby += fmt.Sprintf("%s ASC,", "region")
	}
	// Service
	if self.Service.Selectable() {
		selected += fmt.Sprintf("%s,", "service")
	}
	if self.Service.Whereable() {
		params.Service = string(self.Service)
		where += fmt.Sprintf("%s AND ", "service=:service")
	}
	if self.Service.Groupable() {
		groupby += fmt.Sprintf("%s,", "service")
	}
	if self.Service.Orderable() {
		orderby += fmt.Sprintf("%s ASC,", "service")
	}
	// Account - tag name & label as well, the account id is unique
	if self.Account.Selectable() {
		selected += fmt.Sprintf("%s, %s, %s,", "aws_account_id", "aws_accounts.name as account_name", "aws_accounts.label as account_label")
	}
	if self.Account.Whereable() {
		params.Account = string(self.Account)
		where += fmt.Sprintf("%s AND ", "aws_account_id=:aws_account_id")
	}
	if self.Account.Groupable() {
		groupby += fmt.Sprintf("%s,", "aws_account_id")
	}
	if self.Account.Orderable() {
		orderby += fmt.Sprintf("%s ASC,", "aws_account_id")
	}
	// Environment
	if self.Environment.Selectable() {
		selected += fmt.Sprintf("%s,", "aws_accounts.environment as environment")
	}
	if self.Environment.Whereable() {
		params.Environment = string(self.Environment)
		where += fmt.Sprintf("%s AND ", "aws_accounts.environment=:environment")
	}
	if self.Environment.Groupable() {
		groupby += fmt.Sprintf("%s,", "aws_accounts.environmen")
	}
	if self.Environment.Orderable() {
		orderby += fmt.Sprintf("%s ASC,", "aws_accounts.environmen")
	}

	// Replace the placeholders with the real values
	stmt = strings.ReplaceAll(stmt, "{SELECT}", selected)
	stmt = strings.ReplaceAll(stmt, "{WHERE}", where)
	stmt = strings.ReplaceAll(stmt, "{GROUP_BY}", groupby)
	stmt = strings.ReplaceAll(stmt, "{ORDER_BY}", orderby)

	bound = &sqldb.BoundStatement{Data: params, Statement: stmt}
	return
}
