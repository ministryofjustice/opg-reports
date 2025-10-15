package api

import "opg-reports/report/internal/repository/sqlr"

type AwsUptimeGetter[T Model] interface {
	Closer
	GetAllAwsUptime(store sqlr.RepositoryReader) (data []T, err error)
}

type AwsUptimeGroupedGetter[T Model] interface {
	Closer
	GetAllAwsUptime(store sqlr.RepositoryReader) (data []T, err error)
}

// GetAwsUptineGroupedOptions contains a series of values that determines
// what fields are used within the sql statement to allow for easier
// handling of multiple, similar sql queries that differ by which
// columns are grouped or filtered via the incoming api request.
//
// Should be very similar to `awsUptimeSqlParams` which are the bound
// versions of this data
type GetAwsUptineGroupedOptions struct {
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	DateFormat string `json:"date_format"`
	Team       string `json:"team"`
}

type AwsUptimeGrouped struct {
	Date        string `json:"date,omitempty" db:"date" example:"2019-08-01"`       // The data the cost was incurred - provided from the cost explorer result
	Average     string `json:"average,omitempty" db:"average" example:"99.9501"`    // The average uptime percentage
	Granularity string `json:"granularity,omitempty" db:"granularity" example:"60"` // The time period accuracy in seconds
	// Fields captured via joins in the sql
	TeamName              string `json:"team,omitempty" db:"team_name"`
	AwsAccountID          string `json:"account,omitempty" db:"aws_account_id"`
	AwsAccountName        string `json:"account_name,omitempty" db:"aws_account_name"`
	AwsAccountLabel       string `json:"account_label,omitempty" db:"aws_account_label"`
	AwsAccountEnvironment string `json:"environment,omitempty" db:"aws_account_environment"`
}

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

// GetAllAwsUptime returns all uptime data - please dont use on the api for real as it will be very heavy on the DB
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

// GetGroupedAwsUptime uses a set of options to generate the sql statement that will select, filter,
// group and order by the data set between provided dates.
func (self *Service[T]) GetGroupedAwsUptime(store sqlr.RepositoryReader, options *GetAwsUptineGroupedOptions) (data []T, err error) {
	var selectStmt, _ = awsUptimeGroupedSqlStatement(options)
	var log = self.log.With("operation", "GetGroupedAwsUptime")

	data = []T{}
	log.Debug("getting grouped awsuptime from database ...")

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
