package migrations

import (
	"context"
	"log/slog"
	"opg-reports/report/package/cntxt"
)

const convert_aws_accounts string = `
INSERT INTO accounts (id, name, label, environment, uptime_tracking, team_name)
	SELECT id, name, label, environment, uptime_tracking, team_name FROM aws_accounts;
`

const delete_aws_accounts string = `
DROP INDEX IF EXISTS aws_accounts_id_idx;
DROP TABLE IF EXISTS aws_accounts;
`

// convert old aws cost data into new form
const convert_aws_costs string = `
INSERT INTO costs (created_at, region, service, month, cost, account_id)
	SELECT created_at, region, service, strftime("%Y-%m",date) as date, cost, aws_account_id FROM aws_costs
;`

const delete_aws_costs string = `
DROP INDEX IF EXISTS aws_costs_date_idx;
DROP INDEX IF EXISTS aws_costs_date_account_idx;
DROP INDEX IF EXISTS aws_costs_unique_idx;
DROP TABLE IF EXISTS aws_costs;
`

const convert_aws_uptime string = `
INSERT INTO uptime (month, account_id, average, granularity)
	SELECT strftime("%Y-%m", date) as date, aws_account_id, AVG(average), granularity FROM aws_uptime GROUP BY strftime("%Y-%m", date), aws_account_id
;
`
const delete_aws_uptime string = `
DROP INDEX IF EXISTS aws_uptime_date_idx;
DROP INDEX IF EXISTS aws_uptime_account_date_idx;
DROP TABLE IF EXISTS aws_uptime;
`

var conversions = []*Migration{
	{Key: "convert_aws_accounts", Stmt: convert_aws_accounts},
	{Key: "delete_aws_accounts", Stmt: delete_aws_accounts},
	{Key: "convert_aws_costs", Stmt: convert_aws_costs},
	{Key: "delete_aws_costs", Stmt: delete_aws_costs},
	{Key: "convert_aws_uptime", Stmt: convert_aws_uptime},
	{Key: "delete_aws_uptime", Stmt: delete_aws_uptime},
}

// Convert runs conversions against older versions of the DB to keep data
func Convert(ctx context.Context, flags *Args) (err error) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "global", "func", "Convert")
	log.Info("starting ... ")
	err = runMigrations(ctx, flags, conversions)
	log.Info("complete.")
	return
}
