package api

import (
	"fmt"
	"log/slog"
	"strings"

	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/utils"
)

type AwsCostsGetter[T Model] interface {
	Closer
	GetAllAwsCosts(store sqlr.Reader) (data []T, err error)
}
type AwsCostsTop20Getter[T Model] interface {
	Closer
	GetTop20AwsCosts(store sqlr.Reader) (data []T, err error)
}
type AwsCostsGroupedGetter[T Model] interface {
	Closer
	GetGroupedAwsCosts(store sqlr.Reader, options *GetGroupedCostsOptions) (data []T, err error)
}

// stmtAwsCostsInsert used to insert records into the database the PutX functions
const stmtAwsCostsInsert string = `
INSERT INTO aws_costs (
	region,
	service,
	date,
	cost,
	aws_account_id
) VALUES (
	:region,
	:service,
	:date,
	:cost,
	:aws_account_id
) ON CONFLICT (aws_account_id,date,region,service)
 	DO UPDATE SET cost=excluded.cost
RETURNING id;
`

// stmtAwsCostsSelectAll fetches all records and orders them by cost.
//
// This should never be exposed to the api layer as table
// has millions of rows
const stmtAwsCostsSelectAll string = `
SELECT
	aws_costs.region,
	aws_costs.service,
	aws_costs.date,
	aws_costs.cost,
	json_object(
		'id', aws_accounts.id,
		'name', aws_accounts.name,
		'label', aws_accounts.label,
		'environment', aws_accounts.environment
	) as aws_account,
	json_object(
		'name', aws_accounts.team_name
	) as team
FROM aws_costs
LEFT JOIN aws_accounts on aws_accounts.id = aws_costs.aws_account_id
GROUP BY aws_costs.id
ORDER BY
	CAST(aws_costs.cost as REAL) DESC,
	aws_accounts.team_name ASC,
	aws_accounts.name ASC,
	aws_accounts.environment ASC,
	aws_costs.region ASC,
	aws_costs.service ASC;
`

// stmtSelectTop20 returns the top20 most expensive costs stored.
//
// Excludes tax. as that is worked out on a single day for the
// month so would always fill this list.
const stmtAwsCostsSelectTop20 string = `
SELECT
	aws_costs.region,
	aws_costs.service,
	aws_costs.date,
	aws_costs.cost,
	json_object(
		'id', aws_accounts.id,
		'name', aws_accounts.name,
		'label', aws_accounts.label,
		'environment', aws_accounts.environment
	) as aws_account,
	json_object(
		'name', aws_accounts.team_name
	) as team
FROM aws_costs
LEFT JOIN aws_accounts on aws_accounts.id = aws_costs.aws_account_id
WHERE
	aws_costs.service != 'Tax'
GROUP BY aws_costs.id
ORDER BY
	CAST(aws_costs.cost as REAL) DESC,
	aws_accounts.team_name ASC,
	aws_accounts.name ASC,
	aws_accounts.environment ASC,
	aws_costs.region ASC,
	aws_costs.service ASC
LIMIT 20;
`

// stmtAwsCostsGrouped is the base sql statements used for most cost database calls
// that filters out Tax and groups values by at least the date column.
//
// It contains extra :params to allow extension of the query and typically
// generated from an api input dataset
//
// :date_format = used for date time grouping via strftime on the date column
// :start_date	= lower bound on the date where
// :end_date 	= upper bound on the date where
//
// {SELECT} 	= generated extra columns
// {WHERE} 		= generated where clauses
// {GROUP_BY}	= generated group by
// {ORDER_BY}	= generated order by
const stmtAwsCostsGrouped string = `
SELECT
	{SELECT}
    coalesce(SUM(cost), 0) as cost,
    strftime(:date_format, date) as date
FROM aws_costs
LEFT JOIN aws_accounts ON aws_accounts.id = aws_costs.aws_account_id
WHERE
	{WHERE}
    date >= :start_date
    AND date < :end_date
	AND service != 'Tax'
GROUP BY
	{GROUP_BY}
	strftime(:date_format, date)
ORDER BY
	CAST(aws_costs.cost as REAL) DESC,
	{ORDER_BY}
	strftime(:date_format, date) ASC
;
`

// AwsCost
type AwsCost struct {
	Region  string `json:"region,omitempty" db:"region" example:"eu-west-1|eu-west-2|NoRegion"` // The AWS region
	Service string `json:"service,omitempty" db:"service" example:"AWS service name"`           // The AWS service name
	Date    string `json:"date,omitempty" db:"date" example:"2019-08-24"`                       // The data the cost was incurred - provided from the cost explorer result
	Cost    string `json:"cost,omitempty" db:"cost" example:"-10.537"`                          // The actual cost value as a string - without an currency, but is USD by default

	// Joins
	// AwsAccount joins
	AwsAccountID string            `json:"aws_account_id,omitempty" db:"aws_account_id"`
	AwsAccount   *hasOneAwsAccount `json:"aws_account,omitempty" db:"aws_account"`
	// Team joins
	Team *costHasOneTeam `json:"team,omitempty" db:"team"`
}

type AwsCostGrouped struct {
	Region  string `json:"region,omitempty" db:"region" example:"eu-west-1|eu-west-2|NoRegion"` // The AWS region
	Service string `json:"service,omitempty" db:"service" example:"AWS service name"`           // The AWS service name
	Date    string `json:"date,omitempty" db:"date" example:"2019-08-24"`                       // The data the cost was incurred - provided from the cost explorer result
	Cost    string `json:"cost,omitempty" db:"cost" example:"-10.537"`                          // The actual cost value as a string - without an currency, but is USD by default
	// Fields captured via joins in the sql
	TeamName              string `json:"team_name,omitempty" db:"team_name"`
	AwsAccountID          string `json:"aws_account_id,omitempty" db:"aws_account_id"`
	AwsAccountName        string `json:"aws_account_name,omitempty" db:"account_name"`
	AwsAccountLabel       string `json:"aws_account_label,omitempty" db:"account_label"`
	AwsAccountEnvironment string `json:"aws_account_environment,omitempty" db:"environment"`
}

// awsAccount is used to capture sql join data
type awsAccount struct {
	ID          string `json:"id,omitempty" db:"id" example:"012345678910"` // This is the AWS Account ID as a string
	Name        string `json:"name,omitempty" db:"name" example:"Public API"`
	Label       string `json:"label,omitempty" db:"label" example:"aurora-cluster"`
	Environment string `json:"environment,omitempty" db:"environment" example:"development|preproduction|production"`
}

// hasOneAwsAccount used for the join to awsAccount to handle the scaning into a struct
type hasOneAwsAccount awsAccount

// Scan handles the processing of the join data
func (self *hasOneAwsAccount) Scan(src interface{}) (err error) {
	switch src.(type) {
	case []byte:
		err = utils.Unmarshal(src.([]byte), self)
	case string:
		err = utils.Unmarshal([]byte(src.(string)), self)
	default:
		err = fmt.Errorf("unsupported scan src type")
	}
	return
}

// costTeam maps to the team model
type costTeam struct {
	Name string `json:"name,omitempty" db:"name" example:"SRE"`
}

type costHasOneTeam costTeam

// Scan handles the processing of the join data
func (self *costHasOneTeam) Scan(src interface{}) (err error) {
	switch src.(type) {
	case []byte:
		err = utils.Unmarshal(src.([]byte), self)
	case string:
		err = utils.Unmarshal([]byte(src.(string)), self)
	default:
		err = fmt.Errorf("unsupported scan src type")
	}
	return
}

// sqlParams is used in the GetGroupedCostsOptions.Statement method
// to generate the parameters to bind to the sql
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
func (self *GetGroupedCostsOptions) Statement() (bound *sqlr.BoundStatement, params *sqlParams) {
	var (
		stmt            = stmtAwsCostsGrouped
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
		selected += fmt.Sprintf("%s,", "aws_accounts.team_name as team_name")
	}
	if self.Team.Whereable() {
		params.Team = string(self.Team)
		where += fmt.Sprintf("%s AND ", "aws_accounts.team_name=:team_name")
	}
	if self.Team.Groupable() {
		groupby += fmt.Sprintf("%s,", "aws_accounts.team_name")
	}
	if self.Team.Orderable() {
		orderby += fmt.Sprintf("%s ASC,", "aws_accounts.team_name")
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
		groupby += fmt.Sprintf("%s,", "aws_accounts.environment")
	}
	if self.Environment.Orderable() {
		orderby += fmt.Sprintf("%s ASC,", "aws_accounts.environment")
	}

	// Replace the placeholders with the real values
	stmt = strings.ReplaceAll(stmt, "{SELECT}", selected)
	stmt = strings.ReplaceAll(stmt, "{WHERE}", where)
	stmt = strings.ReplaceAll(stmt, "{GROUP_BY}", groupby)
	stmt = strings.ReplaceAll(stmt, "{ORDER_BY}", orderby)

	bound = &sqlr.BoundStatement{Data: params, Statement: stmt}
	return
}

// GetAllCosts will return all records
//
// Using this is generally a bad idea as this table will contain millions of rows
func (self *Service[T]) GetAllAwsCosts(store sqlr.Reader) (data []T, err error) {
	var selectStmt = &sqlr.BoundStatement{Statement: stmtAwsCostsSelectAll}
	var log = self.log.With("operation", "GetAllCosts")

	data = []T{}
	log.Debug("getting all awscosts from database ... ")

	// cast the data back to struct
	if err = store.Select(selectStmt); err == nil {
		data = selectStmt.Returned.([]T)
	}

	return
}

// GetTop20Costs returns top 20 most expensive costs store in the database
func (self *Service[T]) GetTop20AwsCosts(store sqlr.Reader) (data []T, err error) {
	var selectStmt = &sqlr.BoundStatement{Statement: stmtAwsCostsSelectTop20}
	var log = self.log.With("operation", "GetTop20Costs")

	data = []T{}
	log.Debug("getting top20 awscosts from database ... ")

	// cast the data back to struct
	if err = store.Select(selectStmt); err == nil {
		data = selectStmt.Returned.([]T)
	}
	return
}

// GetGroupedCosts uses a set of options to generate the sql statement that will select, filter,
// group and order by the data set between provided dates.
func (self *Service[T]) GetGroupedAwsCosts(store sqlr.Reader, options *GetGroupedCostsOptions) (data []T, err error) {
	var selectStmt, _ = options.Statement()
	var log = self.log.With("operation", "GetGroupedCosts")

	data = []T{}
	log.Debug("getting grouped awscosts from database ...")

	// cast the data back to struct
	if err = store.Select(selectStmt); err == nil {
		data = selectStmt.Returned.([]T)
	}
	return
}

// PutAwsCosts inserts new cost records into the table.
//
// Expects data to be like:
//
//	[{
//	  "cost": "1.152",
//	  "date": "2025-05-31",
//	  "region": "eu-west-2",
//	  "service": "Amazon Virtual Private Cloud"
//	  "aws_account_id": "01011"
//	}]
//
// Note: Dont expose to the api endpoints
func (self *Service[T]) PutAwsCosts(store sqlr.Writer, data []T) (results []*sqlr.BoundStatement, err error) {
	var (
		inserts []*sqlr.BoundStatement = []*sqlr.BoundStatement{}
		log     *slog.Logger           = self.log.With("operation", "PutAwsCosts")
	)
	results = []*sqlr.BoundStatement{}

	log.Debug("generating insert statements for aws costs")
	// for each cost item generate the insert
	for _, row := range data {
		inserts = append(inserts, &sqlr.BoundStatement{Data: row, Statement: stmtAwsCostsInsert})
	}
	log.With("count", len(inserts)).Debug("inserting records from file ...")

	// run inserts
	if err = store.Insert(inserts...); err != nil {
		log.Error("error inserting", "err", err.Error())
		return
	}
	// only merge in the items with return values
	for _, in := range inserts {
		if in.Returned != nil {
			results = append(results, in)
		}
	}
	if len(results) != len(data) {
		err = fmt.Errorf("not all costs were inserted; expected [%d] actual [%d]", len(data), len(results))
		return
	}

	log.With("inserted", len(results)).Info("inserting records successful")
	return
}
