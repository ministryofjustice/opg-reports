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

// Statement converts the configured options to a bound statement and provides the
// values and :params for `stmGroupedCosts`.
//
// It returns the bound statement and generated data object
func (self *GetGroupedCostsOptions) Statement() (bound *sqldb.BoundStatement, data map[string]string) {
	var (
		stmt            = stmGroupedCosts
		selected string = ""
		where    string = ""
		orderby  string = ""
		groupby  string = ""
	)
	// setup the default data values
	data = map[string]string{
		"start_date":  self.StartDate,
		"end_date":    self.EndDate,
		"date_format": self.DateFormat,
	}

	// check the team, account, env, region and service values and update the
	// sql

	// Team
	if self.Team.Selectable() {
		selected += fmt.Sprintf("%s,", "team_name")
	}
	if self.Team.Whereable() {
		data["team_name"] = string(self.Team)
		where += fmt.Sprintf("%s AND ", "team_name=:team_name")
	}
	if self.Team.Groupable() {
		data["g1"] = "team_name"
		groupby += fmt.Sprintf("%s,", ":g1")
	}
	if self.Team.Orderable() {
		data["o1"] = "team_name"
		orderby += fmt.Sprintf("%s ASC,", ":o1")
	}

	// Region
	if self.Region.Selectable() {
		selected += fmt.Sprintf("%s,", "region")
	}
	if self.Region.Whereable() {
		data["region"] = string(self.Region)
		where += fmt.Sprintf("%s AND ", "region=:region")
	}
	if self.Region.Groupable() {
		data["g2"] = "region"
		groupby += fmt.Sprintf("%s,", ":g2")
	}
	if self.Region.Orderable() {
		data["o2"] = "region"
		orderby += fmt.Sprintf("%s ASC,", ":o2")
	}
	// Service
	if self.Service.Selectable() {
		selected += fmt.Sprintf("%s,", "service")
	}
	if self.Service.Whereable() {
		data["service"] = string(self.Service)
		where += fmt.Sprintf("%s AND ", "service=:service")
	}
	if self.Service.Groupable() {
		data["g3"] = "service"
		groupby += fmt.Sprintf("%s,", ":g3")
	}
	if self.Service.Orderable() {
		data["o3"] = "service"
		orderby += fmt.Sprintf("%s ASC,", ":o4")
	}
	// Account
	if self.Account.Selectable() {
		selected += fmt.Sprintf("%s,", "aws_account_id")
	}
	if self.Account.Whereable() {
		data["aws_account_id"] = string(self.Account)
		where += fmt.Sprintf("%s AND ", "aws_account_id=:aws_account_id")
	}
	if self.Account.Groupable() {
		data["g4"] = "aws_account_id"
		groupby += fmt.Sprintf("%s,", ":g4")
	}
	if self.Account.Orderable() {
		data["o4"] = "aws_account_id"
		orderby += fmt.Sprintf("%s ASC,", ":o4")
	}
	// Environment
	if self.Environment.Selectable() {
		selected += fmt.Sprintf("%s,", "aws_accounts.environment as environment")
	}
	if self.Environment.Whereable() {
		data["environment"] = string(self.Environment)
		where += fmt.Sprintf("%s AND ", "aws_accounts.environment=:environment")
	}
	if self.Environment.Groupable() {
		data["g5"] = "aws_accounts.environment"
		groupby += fmt.Sprintf("%s,", ":g5")
	}
	if self.Environment.Orderable() {
		data["o5"] = "aws_accounts.environment"
		orderby += fmt.Sprintf("%s ASC,", ":o5")
	}

	// Replace the placeholders with the real values
	stmt = strings.ReplaceAll(stmt, "{SELECT}", selected)
	stmt = strings.ReplaceAll(stmt, "{WHERE}", where)
	stmt = strings.ReplaceAll(stmt, "{GROUP_BY}", groupby)
	stmt = strings.ReplaceAll(stmt, "{ORDER_BY}", orderby)

	bound = &sqldb.BoundStatement{Data: data, Statement: stmt}
	return
}
