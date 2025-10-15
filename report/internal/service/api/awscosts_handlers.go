package api

import "opg-reports/report/internal/repository/sqlr"

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
	GetGroupedAwsCosts(store sqlr.RepositoryReader, options *GetAwsCostsGroupedOptions) (data []T, err error)
}

// GetAwsCostsGroupedOptions contains a series of values that determines
// what fields are used within the sql statement to allow for easier
// handling of multiple, similar sql queries that differ by which
// columns are grouped or filtered
type GetAwsCostsGroupedOptions struct {
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

// GetGroupedAwsCosts uses a set of options to generate the sql statement that will select, filter,
// group and order by the data set between provided dates.
func (self *Service[T]) GetGroupedAwsCosts(store sqlr.RepositoryReader, options *GetAwsCostsGroupedOptions) (data []T, err error) {
	var selectStmt, _ = awsCostsGroupedSqlStatement(options)
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
