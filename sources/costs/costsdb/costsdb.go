// Package costsdb provides all the database statements for
// working with costs table
package costsdb

import (
	"github.com/ministryofjustice/opg-reports/pkg/datastore"
)

// CreateCostTable is the create table statement for costs
const CreateCostTable datastore.CreateStatement = `
CREATE TABLE IF NOT EXISTS costs (
    id INTEGER PRIMARY KEY,
    ts TEXT NOT NULL,

    organisation TEXT NOT NULL,
    account_id TEXT NOT NULL,
    account_name TEXT NOT NULL,
    unit TEXT NOT NULL,
    label TEXT NOT NULL,
    environment TEXT NOT NULL,

	region TEXT NOT NULL,
    service TEXT NOT NULL,
    date TEXT NOT NULL,
    cost TEXT NOT NULL
) STRICT
;`

// CreateCostTableIndex is the index creation statements
const CreateCostTableIndex datastore.CreateStatement = `CREATE INDEX IF NOT EXISTS costs_date_idx ON costs(date);`

// InsertCosts is named parameter statement to insert a single entry to
// costs table with the new id being returned
const InsertCosts datastore.InsertStatement = `
INSERT INTO costs(
    ts,
    organisation,
    account_id,
    account_name,
    unit,
    label,
    environment,
    service,
    region,
    date,
    cost
) VALUES (
    :ts,
	:organisation,
	:account_id,
	:account_name,
	:unit,
	:label,
	:environment,
	:service,
	:region,
	:date,
	:cost
) RETURNING id
;`

// RowCount returns the total number of records within the database
const RowCount datastore.SelectStatement = `
SELECT
	count(*) as row_count
FROM costs
LIMIT 1
;`

// Total returns the sum of the cost field for all
// records with the date range passed (start_date, end_date)
// Excludes tax
const Total datastore.SelectStatement = `
SELECT
    coalesce(SUM(cost), 0) as total
FROM costs
WHERE
    date >= ?
	AND date < ?
	AND service != 'Tax'
LIMIT 1
;`

// TaxOverview is used to calculate the total costs
// within the given date range (>= :start_date, < :end_date) and
// splits that based on the `service` being 'Tax' or not
// Used to show top line numbers without VAT etc
const TaxOverview datastore.NamedSelectStatement = `
SELECT
    'Including Tax' as service,
    coalesce(SUM(cost), 0) as cost,
    strftime(:date_format, date) as date
FROM costs as incTax
WHERE
    incTax.date >= :start_date
    AND incTax.date < :end_date
GROUP BY strftime(:date_format, incTax.date)
UNION ALL
SELECT
    'Excluding Tax' as service,
    coalesce(SUM(cost), 0) as cost,
    strftime(:date_format, date) as date
FROM costs as excTax
WHERE
    excTax.date >= :start_date
    AND excTax.date < :end_date
	AND excTax.service != 'Tax'
GROUP BY strftime(:date_format, date)
ORDER by date ASC
;`

// PerUnit groups the cost data by the time period and unit
// and limits the data to the date range specfied
// (>= :start_date, < :end_date) returning the SUM costs
// for each grouping
// Excludes tax
const PerUnit datastore.NamedSelectStatement = `
SELECT
    unit,
    coalesce(SUM(cost), 0) as cost,
    strftime(:date_format, date) as date
FROM costs
WHERE
    date >= :start_date
    AND date < :end_date
	AND service != 'Tax'
GROUP BY strftime(:date_format, date), unit
ORDER by strftime(:date_format, date), unit ASC
`

// PerUnitForUnit operates like PerUnit but also filters
// the result on unit
// Excludes tax
const PerUnitForUnit datastore.NamedSelectStatement = `
SELECT
    unit,
    coalesce(SUM(cost), 0) as cost,
    strftime(:date_format, date) as date
FROM costs
WHERE
    date >= :start_date
    AND date < :end_date
	AND service != 'Tax'
	AND unit = :unit
GROUP BY strftime(:date_format, date), unit
ORDER by strftime(:date_format, date), unit ASC
`

// PerUnitEnvironment groups cost date by the date period, unit
// and environment values in the row and restricts the dataset to the
// date range passed (>= :start_date, < :end_date) returning the
// SUM of each grouping as `cost`
// If the environment field is "null" then we default to "production"
// as several accounts (like managment / identity ) have only one
// version
// Excludes tax
const PerUnitEnvironment datastore.NamedSelectStatement = `
SELECT
    unit,
	IIF(environment != "null", environment, "production") as environment,
    coalesce(SUM(cost), 0) as cost,
    strftime(:date_format, date) as date
FROM costs
WHERE
    date >= :start_date
    AND date < :end_date
	AND service != 'Tax'
GROUP BY strftime(:date_format, date), unit, environment
ORDER by strftime(:date_format, date), unit, environment ASC
`

// PerUnitEnvironmentForUnit expands PerUnitEnvironment by allowing
// filtering by unit
const PerUnitEnvironmentForUnit datastore.NamedSelectStatement = `
SELECT
    unit,
	IIF(environment != "null", environment, "production") as environment,
    coalesce(SUM(cost), 0) as cost,
    strftime(:date_format, date) as date
FROM costs
WHERE
    date >= :start_date
    AND date < :end_date
	AND service != 'Tax'
	AND unit = :unit
GROUP BY strftime(:date_format, date), unit, environment
ORDER by strftime(:date_format, date), unit, environment ASC
`

// Detailed is used to show the cost of each type of service per account and
// org for the time period passed along - allowing to track costs changes
// for s3 etc overtime at a granular level
// Limits the data to the date range expressed (>= :start_date, < :end_date)
// Excludes tax
const Detailed datastore.NamedSelectStatement = `
SELECT
    unit,
	IIF(environment != "null", environment, "production") as environment,
	account_id,
	account_name,
	label,
	service,
    coalesce(SUM(cost), 0) as cost,
    strftime(:date_format, date) as date
FROM costs
WHERE
    date >= :start_date
    AND date < :end_date
	AND service != 'Tax'
GROUP BY strftime(:date_format, date), unit, environment, account_id, service
ORDER by strftime(:date_format, date), unit, environment, account_id, service ASC
`

// DetailedForUnit is an extension of Detailed and further restricts the data set
// to match the unit passed
// Excludes tax
const DetailedForUnit datastore.NamedSelectStatement = `
SELECT
    unit,
	IIF(environment != "null", environment, "production") as environment,
	account_id,
	account_name,
	label,
	service,
    coalesce(SUM(cost), 0) as cost,
    strftime(:date_format, date) as date
FROM costs
WHERE
    date >= :start_date
    AND date < :end_date
	AND service != 'Tax'
	AND unit = :unit
GROUP BY strftime(:date_format, date), unit, environment, account_id, service
ORDER by strftime(:date_format, date), unit, environment, account_id, service  ASC
`
