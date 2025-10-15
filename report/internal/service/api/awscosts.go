package api

import (
	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/utils"
)

type AwsCostsGetter[T Model] interface {
	Closer
	GetAllAwsCosts(store sqlr.RepositoryReader) (data []T, err error)
}
type AwsCostsTop20Getter[T Model] interface {
	Closer
	GetTop20AwsCosts(store sqlr.RepositoryReader) (data []T, err error)
}
type AwsCostsGroupedGetter[T Model] interface {
	Closer
	GetGroupedAwsCosts(store sqlr.RepositoryReader, options *GetGroupedCostsOptions) (data []T, err error)
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

type AwsCostGrouped struct {
	Region  string `json:"region,omitempty" db:"region" example:"eu-west-1|eu-west-2|NoRegion"` // The AWS region
	Service string `json:"service,omitempty" db:"service" example:"AWS service name"`           // The AWS service name
	Date    string `json:"date,omitempty" db:"date" example:"2019-08-24"`                       // The data the cost was incurred - provided from the cost explorer result
	Cost    string `json:"cost,omitempty" db:"cost" example:"-10.537"`                          // The actual cost value as a string - without an currency, but is USD by default
	// Fields captured via joins in the sql
	TeamName              string `json:"team,omitempty" db:"team_name"`
	AwsAccountID          string `json:"account,omitempty" db:"aws_account_id"`
	AwsAccountName        string `json:"account_name,omitempty" db:"aws_account_name"`
	AwsAccountLabel       string `json:"account_label,omitempty" db:"aws_account_label"`
	AwsAccountEnvironment string `json:"environment,omitempty" db:"aws_account_environment"`
}

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
	Team *hasOneTeam `json:"team,omitempty" db:"team"`
}

// awsCostsSqlParams is used in the GetGroupedCostsOptions.Statement method
// to generate the parameters to bind to the sql
type awsCostsSqlParams struct {
	StartDate   string `json:"start_date,omitempty" db:"start_date"`
	EndDate     string `json:"end_date,omitempty" db:"end_date"`
	DateFormat  string `json:"date_format,omitempty" db:"date_format"`
	Region      string `json:"region,omitempty" db:"region"`
	Service     string `json:"service,omitempty" db:"service"`
	Team        string `json:"team_name,omitempty" db:"team_name"`
	Account     string `json:"aws_account_id,omitempty" db:"aws_account_id"`
	AccountName string `json:"aws_account_name,omitempty" db:"aws_account_name"`
	Label       string `json:"aws_account_label,omitempty" db:"aws_account_label"`
	Environment string `json:"aws_account_environment,omitempty" db:"aws_account_environment"`
}

// GetGroupedCostsOptions contains a series of values that determines
// what fields are used within the sql statement to allow for easier
// handling of multiple, similar sql queries that differ by which
// columns are grouped or filtered
type GetGroupedCostsOptions struct {
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	DateFormat  string `json:"date_format"`
	Team        string `json:"team"`
	Region      string `json:"region"`
	Service     string `json:"service"`
	Account     string `json:"account"`
	AccountName string `json:"account_name"`
	Label       string `json:"account_label"`
	Environment string `json:"environment"`
}

// Fields is used to generate sql statement from the dynamic input values
//
// Base sql statements used for most cost database calls
// that filters out Tax and groups values by at least the date column.
//
// Each Field contains the information required for each part of the select
// statement
//
// Uses `:value` placeholders that are mapped out by the statement and
// relate to the `db` attribute on `awsCostsSqlParams`
func (self *GetGroupedCostsOptions) Fields() []*Field {
	return []*Field{
		// exclude tax
		&Field{
			Key:   "tax",
			Where: "service != 'Tax'",
		},
		&Field{
			Key:     "cost",
			Select:  "coalesce(SUM(cost), 0) as cost",
			OrderBy: "CAST(coalesce(SUM(cost), 0) as REAL) DESC",
		},
		&Field{
			Key:     "date",
			Select:  "strftime(:date_format, date) as date",
			Where:   "(date >= :start_date AND date <= :end_date)",
			GroupBy: "strftime(:date_format, date)",
			OrderBy: "strftime(:date_format, date) ASC",
		},
		// Region
		&Field{
			Key:     "region",
			Select:  "region",
			Where:   "region=:region",
			GroupBy: "region",
			OrderBy: "region ASC",
			Value:   utils.Ptr(self.Region),
		},
		// Service
		&Field{
			Key:     "service",
			Select:  "service",
			Where:   "service=:service",
			GroupBy: "service",
			OrderBy: "service ASC",
			Value:   utils.Ptr(self.Service),
		},
		// AWS account id
		&Field{
			Key:     "aws_account_id",
			Select:  "aws_account_id",
			Where:   "aws_account_id=:aws_account_id",
			GroupBy: "aws_account_id",
			OrderBy: "aws_account_id ASC",
			Value:   utils.Ptr(self.Account),
		},
		// AWS account name
		&Field{
			Key:     "name",
			Select:  "aws_accounts.name as aws_account_name",
			Where:   "aws_accounts.name=:aws_account_name",
			GroupBy: "aws_accounts.name",
			OrderBy: "aws_accounts.name ASC",
			Value:   utils.Ptr(self.AccountName),
		},
		// AWS team name
		&Field{
			Key:     "team",
			Select:  "aws_accounts.team_name as team_name",
			Where:   "lower(aws_accounts.team_name)=lower(:team_name)",
			GroupBy: "aws_accounts.team_name",
			OrderBy: "aws_accounts.team_name ASC",
			Value:   utils.Ptr(self.Team),
		},
		// AWS environment
		&Field{
			Key:     "environment",
			Select:  "aws_accounts.environment as aws_account_environment",
			Where:   "aws_accounts.environment=:aws_account_environment",
			GroupBy: "aws_accounts.environment",
			OrderBy: "aws_accounts.environment ASC",
			Value:   utils.Ptr(self.Environment),
		},
		// AWS label
		&Field{
			Key:     "label",
			Select:  "aws_accounts.label as aws_account_label",
			Where:   "aws_accounts.label=:aws_account_label",
			GroupBy: "aws_accounts.label",
			OrderBy: "aws_accounts.label ASC",
			Value:   utils.Ptr(self.Label),
		},
	}
}

// Groups is used to show how the costs where grouped together. Typically
// this is used by api recievers to generate table headers and so on
//
// Order does matter: Team, Environment, Account, Service, Region
func (self *GetGroupedCostsOptions) Groups() (groups []string) {
	groups = []string{}

	mapped := map[string]string{}
	utils.Convert(self, &mapped)

	for k, v := range mapped {
		if v == "true" {
			groups = append(groups, k)
		}
	}

	return
}

// Statement converts the configured options to a bound statement and provides the
// values and :params for `stmGroupedCosts`.
//
// It returns the bound statement and generated data object
func (self *GetGroupedCostsOptions) Statement() (bound *sqlr.BoundStatement, params *awsCostsSqlParams) {
	var (
		fields []*Field = self.Fields()
		join   string   = `LEFT JOIN aws_accounts ON aws_accounts.id = aws_costs.aws_account_id`
		stmt   string   = BuildSelectFromFields("aws_costs", join, fields...)
	)
	params = &awsCostsSqlParams{
		StartDate:   self.StartDate,
		EndDate:     self.EndDate,
		DateFormat:  self.DateFormat,
		Team:        self.Team,
		Region:      self.Region,
		Service:     self.Service,
		Account:     self.Account,
		AccountName: self.AccountName,
		Label:       self.Label,
		Environment: self.Environment,
	}
	bound = &sqlr.BoundStatement{Data: params, Statement: stmt}
	return
}

// GetAllCosts will return all records
//
// Using this is generally a bad idea as this table will contain millions of rows
func (self *Service[T]) GetAllAwsCosts(store sqlr.RepositoryReader) (data []T, err error) {
	var selectStmt = &sqlr.BoundStatement{Statement: stmtAwsCostsSelectAll}
	var log = self.log.With("operation", "GetAllAwsCosts")

	data = []T{}
	log.Debug("getting all awscosts from database ... ")

	// cast the data back to struct
	if err = store.Select(selectStmt); err == nil {
		data = selectStmt.Returned.([]T)
	}

	return
}

// GetTop20Costs returns top 20 most expensive costs store in the database
func (self *Service[T]) GetTop20AwsCosts(store sqlr.RepositoryReader) (data []T, err error) {
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
func (self *Service[T]) GetGroupedAwsCosts(store sqlr.RepositoryReader, options *GetGroupedCostsOptions) (data []T, err error) {
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
func (self *Service[T]) PutAwsCosts(store sqlr.RepositoryWriter, data []T) (results []*sqlr.BoundStatement, err error) {

	return self.Put(store, stmtAwsCostsInsert, data)
}
