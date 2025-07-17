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
    AND date <= :end_date
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
	Team *costHasOneTeam `json:"team,omitempty" db:"team"`
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
	AccountName string `json:"aws_account_name,omitempty" db:"aws_account_name"`
	Label       string `json:"aws_account_label,omitempty" db:"aws_account_label"`
	Environment string `json:"aws_account_environment,omitempty" db:"aws_account_environment"`
}

// GetGroupedCostsOptions contains a series of values that determines
// what fields are used within the sql statement to allow for easier
// handling of multiple, similar sql queries that differ by which
// columns are grouped or filtered
type GetGroupedCostsOptions struct {
	StartDate  string `json:"-"`
	EndDate    string `json:"-"`
	DateFormat string `json:"-"`

	Team        utils.TrueOrFilter `json:"team"`
	Region      utils.TrueOrFilter `json:"region"`
	Service     utils.TrueOrFilter `json:"service"`
	Account     utils.TrueOrFilter `json:"account"`
	AccountName utils.TrueOrFilter `json:"account_name"`
	Label       utils.TrueOrFilter `json:"account_label"`
	Environment utils.TrueOrFilter `json:"environment"`
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

type whereableInfo struct {
	Param  func()
	Option utils.TrueOrFilter
}

// Statement converts the configured options to a bound statement and provides the
// values and :params for `stmGroupedCosts`.
//
// It returns the bound statement and generated data object
func (self *GetGroupedCostsOptions) Statement() (bound *sqlr.BoundStatement, params *sqlParams) {
	var (
		stmt     string                        = stmtAwsCostsGrouped
		selected string                        = ""
		where    string                        = ""
		orderby  string                        = ""
		groupby  string                        = ""
		selects  map[string]utils.TrueOrFilter = map[string]utils.TrueOrFilter{}
		groups   map[string]utils.TrueOrFilter = map[string]utils.TrueOrFilter{}
		orders   map[string]utils.TrueOrFilter = map[string]utils.TrueOrFilter{}
		wheres   map[string]*whereableInfo     = map[string]*whereableInfo{}
	)
	// setup the default data values
	params = &sqlParams{
		StartDate:  self.StartDate,
		EndDate:    self.EndDate,
		DateFormat: self.DateFormat,
	}

	// generate the map between select alias and the option item
	selects = map[string]utils.TrueOrFilter{
		"region":                                self.Region,
		"service":                               self.Service,
		"aws_account_id":                        self.Account,
		"aws_accounts.team_name as team_name":   self.Team,
		"aws_accounts.name as aws_account_name": self.AccountName,
		"aws_accounts.environment as aws_account_environment": self.Environment,
		"aws_accounts.label as aws_account_label":             self.Label,
	}
	// add them to the select if they are marked as selectable
	for selectAs, tf := range selects {
		if tf.Selectable() {
			selected += fmt.Sprintf("%s, ", selectAs)
		}
	}

	groups = map[string]utils.TrueOrFilter{
		"region":                   self.Region,
		"service":                  self.Service,
		"aws_account_id":           self.Account,
		"aws_accounts.team_name":   self.Team,
		"aws_accounts.name":        self.AccountName,
		"aws_accounts.environment": self.Environment,
		"aws_accounts.label":       self.Label,
	}
	// add them to the select if they are marked as selectable
	for groupAs, tf := range groups {
		if tf.Groupable() {
			groupby += fmt.Sprintf("%s, ", groupAs)
		}
	}

	orders = map[string]utils.TrueOrFilter{
		"region":                   self.Region,
		"service":                  self.Service,
		"aws_account_id":           self.Account,
		"aws_accounts.team_name":   self.Team,
		"aws_accounts.name":        self.AccountName,
		"aws_accounts.environment": self.Environment,
		"aws_accounts.label":       self.Label,
	}
	// add them to the select if they are marked as selectable
	for orderAs, tf := range orders {
		if tf.Orderable() {
			orderby += fmt.Sprintf("%s ASC,", orderAs)
		}
	}
	// filtering / where
	wheres = map[string]*whereableInfo{
		"region=:region":                                    {Param: func() { params.Region = string(self.Region) }, Option: self.Region},
		"service=:service":                                  {Param: func() { params.Service = string(self.Service) }, Option: self.Service},
		"aws_account_id=:aws_account_id":                    {Param: func() { params.Account = string(self.Account) }, Option: self.Account},
		"aws_accounts.name=:aws_account_name":               {Param: func() { params.AccountName = string(self.AccountName) }, Option: self.AccountName},
		"aws_accounts.environment=:aws_account_environment": {Param: func() { params.Environment = string(self.Environment) }, Option: self.Environment},
		"aws_accounts.label=:aws_account_label":             {Param: func() { params.Label = string(self.Label) }, Option: self.Label},
		"lower(aws_accounts.team_name)=lower(:team_name)":   {Param: func() { params.Team = string(self.Team) }, Option: self.Team},
	}

	for whereAs, wi := range wheres {
		if wi.Option.Whereable() {
			wi.Param()
			where += fmt.Sprintf("%s AND ", whereAs)
		}
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
func (self *Service[T]) GetAllAwsCosts(store sqlr.RepositoryReader) (data []T, err error) {
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
