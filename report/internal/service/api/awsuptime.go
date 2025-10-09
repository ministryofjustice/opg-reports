package api

import (
	"opg-reports/report/internal/repository/sqlr"
)

type AwsUptimeGetter[T Model] interface {
	Closer
	GetAllAwsUptime(store sqlr.RepositoryReader) (data []T, err error)
}

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

// AwsUptime
type AwsUptime struct {
	Date        string `json:"date,omitempty" db:"date" example:"2019-08-01"`       // The data the cost was incurred - provided from the cost explorer result
	Average     string `json:"average,omitempty" db:"average" example:"99.9501"`    // The average uptime percentage
	Granularity string `json:"granularity,omitempty" db:"granularity" example:"60"` // The time period accuracy in seconds
	// Joins
	AwsAccountID string            `json:"aws_account_id,omitempty" db:"aws_account_id"` // AwsAccount join key
	AwsAccount   *hasOneAwsAccount `json:"aws_account,omitempty" db:"aws_account"`       // AwsAccount join model via sql
	Team         *hasOneAwsAccount `json:"team,omitempty" db:"team"`                     // Team join model via sql
}

func (self *Service[T]) GetAllAwsUptime(store sqlr.RepositoryReader) (data []T, err error) {
	var selectStmt = &sqlr.BoundStatement{Statement: stmtAwsUptimeSelectAll}
	var log = self.log.With("operation", "GetAllAwsUptime")

	data = []T{}
	log.Debug("getting all awsuptime from database ... ")

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
//	  "date": "2025-05-17",
//	  "average": "98.99"
//	  "granularity": "60",
//	  "aws_account_id": "01011"
//	}]
//
// Note: Dont expose to the api endpoints
func (self *Service[T]) PutAwsUptime(store sqlr.RepositoryWriter, data []T) (results []*sqlr.BoundStatement, err error) {
	return self.Put(store, stmtAwsUptimeInsert, data)
}
