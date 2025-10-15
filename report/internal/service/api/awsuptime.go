package api

import (
	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/utils"
)

// stmtAwsUptimeInsert used to insert records into the database the PutX functions
const stmtAwsUptimeInsert string = `
INSERT INTO aws_uptime (
	date,
	average,
	granularity,
	aws_account_id
) VALUES (
	:date,
	:average,
	:granularity,
	:aws_account_id
) ON CONFLICT (aws_account_id,date)
 	DO UPDATE SET average=excluded.average, granularity=excluded.granularity
RETURNING id;
`

// stmtAwsUptimeSelectAll fetches all records and orders them by cost.
//
// This should never be exposed to the api layer as table
// has millions of rows
const stmtAwsUptimeSelectAll string = `
SELECT
	aws_uptime.date,
	aws_uptime.average,
	aws_uptime.granularity,
	json_object(
		'id', aws_accounts.id,
		'name', aws_accounts.name,
		'label', aws_accounts.label,
		'environment', aws_accounts.environment
	) as aws_account,
	json_object(
		'name', aws_accounts.team_name
	) as team
FROM aws_uptime
LEFT JOIN aws_accounts on aws_accounts.id = aws_uptime.aws_account_id
GROUP BY aws_uptime.id
ORDER BY
	CAST(aws_uptime.average as REAL) DESC,
	aws_uptime.date DESC,
	aws_accounts.team_name ASC,
	aws_accounts.name ASC,
	aws_accounts.environment ASC;
`

// awsUptimeSqlParams is used in the Statement method
// to generate the parameters to bind to the sql
type awsUptimeSqlParams struct {
	StartDate  string `json:"start_date,omitempty" db:"start_date"`
	EndDate    string `json:"end_date,omitempty" db:"end_date"`
	DateFormat string `json:"date_format,omitempty" db:"date_format"`
	Team       string `json:"team_name,omitempty" db:"team_name"`
}

// awsUptimeGroupedSqlFields is used to generate sql statement from the dynamic input values
//
// Base sql statements used for most cost database calls
// that filters out Tax and groups values by at least the date column.
//
// Each Field contains the information required for each part of the select
// statement
//
// Uses `:value` placeholders that are mapped out by the statement and
// relate to the `db` attribute on `awsUptimeSqlParams`
func awsUptimeGroupedSqlFields(options *GetAwsUptineGroupedOptions) []*Field {
	return []*Field{
		&Field{
			Key:     "average",
			Select:  "AVG(average) as average",
			OrderBy: "CAST(AVG(average) as REAL) DESC",
		},
		&Field{
			Key:     "date",
			Select:  "strftime(:date_format, date) as date",
			Where:   "(date >= :start_date AND date <= :end_date)",
			GroupBy: "strftime(:date_format, date)",
			OrderBy: "strftime(:date_format, date) ASC",
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
	}
}

// awsUptimeGroupedSqlStatement converts the configured options to a bound statement and provides the
// values and :params for `stmGroupedCosts`.
//
// It returns the bound statement and generated data object
func awsUptimeGroupedSqlStatement(options *GetAwsUptineGroupedOptions) (bound *sqlr.BoundStatement, params *awsCostsSqlParams) {
	var (
		fields []*Field = awsUptimeGroupedSqlFields(options)
		table  string   = "aws_costs"
		join   string   = `LEFT JOIN aws_accounts ON aws_accounts.id = aws_costs.aws_account_id`
		stmt   string   = BuildSelectFromFields(table, join, fields...)
	)
	params = &awsCostsSqlParams{
		StartDate:  options.StartDate,
		EndDate:    options.EndDate,
		DateFormat: options.DateFormat,
		Team:       options.Team,
	}
	bound = &sqlr.BoundStatement{Data: params, Statement: stmt}
	return
}
