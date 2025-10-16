package api

import (
	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/utils"
)

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

// awsCostsSqlParams is used in the GetAwsCostsGroupedOptions.Statement method
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

// awsCostsGroupedSqlFields is used to generate sql statement from the dynamic input values
//
// Base sql statements used for most cost database calls
// that filters out Tax and groups values by at least the date column.
//
// Each Field contains the information required for each part of the select
// statement
//
// Uses `:value` placeholders that are mapped out by the statement and
// relate to the `db` attribute on `awsCostsSqlParams`
func awsCostsGroupedSqlFields(options *GetAwsCostsGroupedOptions) []*Field {
	return []*Field{
		// exclude tax
		&Field{
			Key:   "tax",
			Where: "service != 'Tax'",
		},
		&Field{
			Key:     "date",
			Select:  "strftime(:date_format, date) as date",
			Where:   "(date >= :start_date AND date <= :end_date)",
			GroupBy: "strftime(:date_format, date)",
			OrderBy: "strftime(:date_format, date) DESC",
		},
		&Field{
			Key:     "cost",
			Select:  "coalesce(SUM(cost), 0) as cost",
			OrderBy: "CAST(coalesce(SUM(cost), 0) as REAL) DESC",
		},
		// Region
		&Field{
			Key:     "region",
			Select:  "region",
			Where:   "region=:region",
			GroupBy: "region",
			OrderBy: "region ASC",
			Value:   utils.Ptr(options.Region),
		},
		// Service
		&Field{
			Key:     "service",
			Select:  "service",
			Where:   "service=:service",
			GroupBy: "service",
			OrderBy: "service ASC",
			Value:   utils.Ptr(options.Service),
		},
		// AWS account id
		&Field{
			Key:     "aws_account_id",
			Select:  "aws_account_id",
			Where:   "aws_account_id=:aws_account_id",
			GroupBy: "aws_account_id",
			OrderBy: "aws_account_id ASC",
			Value:   utils.Ptr(options.Account),
		},
		// AWS account name
		&Field{
			Key:     "name",
			Select:  "aws_accounts.name as aws_account_name",
			Where:   "aws_accounts.name=:aws_account_name",
			GroupBy: "aws_accounts.name",
			OrderBy: "aws_accounts.name ASC",
			Value:   utils.Ptr(options.AccountName),
		},
		// AWS team name
		&Field{
			Key:     "team",
			Select:  "aws_accounts.team_name as team_name",
			Where:   "lower(aws_accounts.team_name)=lower(:team_name)",
			GroupBy: "aws_accounts.team_name",
			OrderBy: "aws_accounts.team_name ASC",
			Value:   utils.Ptr(options.Team),
		},
		// AWS environment
		&Field{
			Key:     "environment",
			Select:  "aws_accounts.environment as aws_account_environment",
			Where:   "aws_accounts.environment=:aws_account_environment",
			GroupBy: "aws_accounts.environment",
			OrderBy: "aws_accounts.environment ASC",
			Value:   utils.Ptr(options.Environment),
		},
		// AWS label
		&Field{
			Key:     "label",
			Select:  "aws_accounts.label as aws_account_label",
			Where:   "aws_accounts.label=:aws_account_label",
			GroupBy: "aws_accounts.label",
			OrderBy: "aws_accounts.label ASC",
			Value:   utils.Ptr(options.Label),
		},
	}
}

// awsCostsGroupedSqlStatement converts the configured options to a bound statement and provides the
// values and :params for `stmGroupedCosts`.
//
// It returns the bound statement and generated data object
func awsCostsGroupedSqlStatement(options *GetAwsCostsGroupedOptions) (bound *sqlr.BoundStatement, params *awsCostsSqlParams) {
	var (
		fields []*Field = awsCostsGroupedSqlFields(options)
		table  string   = "aws_costs"
		join   string   = `LEFT JOIN aws_accounts ON aws_accounts.id = aws_costs.aws_account_id`
		stmt   string   = BuildSelectFromFields(table, join, fields...)
	)
	params = &awsCostsSqlParams{
		StartDate:   options.StartDate,
		EndDate:     options.EndDate,
		DateFormat:  options.DateFormat,
		Team:        options.Team,
		Region:      options.Region,
		Service:     options.Service,
		Account:     options.Account,
		AccountName: options.AccountName,
		Label:       options.Label,
		Environment: options.Environment,
	}
	bound = &sqlr.BoundStatement{Data: params, Statement: stmt}
	return
}
