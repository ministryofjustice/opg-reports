package infracosts

import (
	"opg-reports/report/internal/utils/query"
	"opg-reports/report/internal/utils/times"
)

// fixed values for this endpoint, used by the operation setup for huma
const (
	ENDPOINT      string = `/v1/infracosts/range/{date_range}/`
	opID          string = `infracosts-by-month`
	opSummary     string = `Return costs grouped by the month and other filter options.`
	opDescription string = `Returns a table of costs`
)

type InfracostRequest struct {
	DateRange   string `json:"date_range" path:"date_range" required:"true" doc:"Date range to use." example:"2025-01..2025-02" pattern:"([0-9]{4}-[0-9]{2}..[0-9]{4}-[0-9]{2})"` // required - date range input
	Team        string `query:"team" json:"team"`
	Account     string `query:"account" json:"account"`
	Environment string `query:"environment" json:"environment"`
	Service     string `query:"service" json:"service"`
}

// Months returns all months between dates
func (self *InfracostRequest) Months() (months []string) {
	months = times.AsYMStrings(times.FromMonthRangeString(self.DateRange))
	return
}

type filter struct {
	Months      string `db:"months" json:"months"`
	Team        string `db:"team" json:"team"`
	Environment string `db:"environment" json:"environment"`
	Service     string `db:"service" json:"service"`
}

var selectQuery = &query.Select{
	// table name & alias
	From: "infracost as base",
	// left join on to accounts
	Joins: "LEFT JOIN accounts ON accounts.id = base.account_id",
	// segment data that is used to build the select statement from the request
	Segments: map[string][]*query.Segment{
		"date_range": []*query.Segment{
			{Select: `strftime("%Y-%m", base.date) as date`},
			{Select: `CAST(COALESCE(SUM(cost), 0) as REAL) as cost`},
			{Where: `base.service != 'Tax'`},
			{Where: `strftime("%Y-%m", base.date) IN (:months)`},
			{GroupBy: `strftime("%Y-%m", base.date)`},
		},
		"team": []*query.Segment{
			{Select: `accounts.team_name as team`},
			{Where: `accounts.team_name = :team`},
			{GroupBy: `accounts.team_name`},
		},
		"account": []*query.Segment{
			{Select: `accounts.name as account_name`},
			{Select: `accounts.id as account_id`},
			{GroupBy: `accounts.name`},
		},
		"environment": []*query.Segment{
			{Select: `accounts.environment as account_environment`},
			{Where: `accounts.environment = :environment`},
			{GroupBy: `accounts.environment`},
		},
		"service": []*query.Segment{
			{Select: `base.service as service`},
			{Where: `base.service = :service`},
			{GroupBy: `base.service`},
		},
	}}
