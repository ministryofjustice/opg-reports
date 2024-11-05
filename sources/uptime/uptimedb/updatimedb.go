package uptimedb

import "github.com/ministryofjustice/opg-reports/pkg/datastore"

const CreateUptimeTable datastore.CreateStatement = `
CREATE TABLE uptime (
    id INTEGER PRIMARY KEY,
    ts TEXT NOT NULL,
    unit TEXT NOT NULL,
    date TEXT NOT NULL,
	average REAL NOT NULL
) STRICT
;`

const CreateUptimeTableDateIndex datastore.CreateStatement = `CREATE INDEX uptime_date_idx ON uptime(date);`
const CreateUptimeTableUnitDateIndex datastore.CreateStatement = `CREATE INDEX uptime_unit_date_idx ON uptime(unit,date);`

const InsertUptime datastore.InsertStatement = `
INSERT INTO uptime(
    ts,
    unit,
    date,
    average
) VALUES (
    :ts,
	:unit,
	:date,
	:average
) RETURNING id
;`

// RowCount returns the total number of records within the database
const RowCount datastore.SelectStatement = `
SELECT
	count(*) as row_count
FROM uptime
LIMIT 1
`

const Uptime datastore.NamedSelectStatement = `
SELECT
    (coalesce(SUM(average), 0) / count(*) ) as average,
    strftime(:date_format, date) as date
FROM costs
WHERE
    date >= :start_date
    AND date < :end_date
GROUP BY strftime(:date_format, date)
ORDER by strftime(:date_format, date) ASC
`

const UptimePerUnit datastore.NamedSelectStatement = `
SELECT
    (coalesce(SUM(average), 0) / count(*) ) as average,
    strftime(:date_format, date) as date
FROM costs
WHERE
    date >= :start_date
    AND date < :end_date
	AND unit = :unit
GROUP BY strftime(:date_format, date), unit
ORDER by strftime(:date_format, date), unit ASC
`
