package infracostseeds

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/infracosts/infracostimports"
	"opg-reports/report/internal/domain/infracosts/infracostmodels"
	"opg-reports/report/internal/utils/times"
	"time"

	"github.com/jmoiron/sqlx"
)

var ErrSeedImportFailed = errors.New("seed costs call failed with an error.")

var (
	costDate string
	seeds    []*infracostmodels.Cost
)

func init() {
	costDate = times.AsString(times.Add(time.Now(), -1, times.MONTH), times.YMD)

	seeds = []*infracostmodels.Cost{
		// account 001A
		{Region: "eu-west-1", Service: "ECS", Date: costDate, Cost: "-0.01", AccountID: "001A"},
		{Region: "eu-west-1", Service: "S3", Date: costDate, Cost: "10.10", AccountID: "001A"},
		{Region: "eu-west-1", Service: "RDS", Date: costDate, Cost: "100.57", AccountID: "001A"},
		{Region: "eu-west-1", Service: "SQS", Date: costDate, Cost: "23.01", AccountID: "001A"},
		{Region: "eu-west-2", Service: "IAM", Date: costDate, Cost: "0.002", AccountID: "001A"},
		// account 001B
		{Region: "eu-west-1", Service: "ECS", Date: costDate, Cost: "-50.21", AccountID: "001B"},
		{Region: "eu-west-2", Service: "S3", Date: costDate, Cost: "603.15", AccountID: "001B"},
		{Region: "eu-west-1", Service: "RDS", Date: costDate, Cost: "105.7", AccountID: "001B"},
		{Region: "eu-west-1", Service: "R53", Date: costDate, Cost: "1.7001", AccountID: "001B"},
		{Region: "us-west-1", Service: "EKS", Date: costDate, Cost: "27501.88", AccountID: "001B"},
		// account 002A
		{Region: "eu-west-1", Service: "ECS", Date: costDate, Cost: "1.02", AccountID: "002A"},
		{Region: "eu-west-2", Service: "S3", Date: costDate, Cost: "37.00", AccountID: "002A"},
		{Region: "eu-west-1", Service: "RDS", Date: costDate, Cost: "-300.68", AccountID: "002A"},
		{Region: "eu-west-1", Service: "SNS", Date: costDate, Cost: "103.51", AccountID: "002A"},
		{Region: "eu-west-2", Service: "RDS", Date: costDate, Cost: "502.44", AccountID: "002A"},
		// account 003A
		{Region: "eu-west-1", Service: "ECS", Date: costDate, Cost: "102.44", AccountID: "003A"},
		{Region: "eu-west-2", Service: "S3", Date: costDate, Cost: "7.0012", AccountID: "003A"},
		{Region: "eu-west-1", Service: "S3", Date: costDate, Cost: "96.35", AccountID: "003A"},
		{Region: "eu-west-1", Service: "SNS", Date: costDate, Cost: "18.19", AccountID: "003A"},
		{Region: "us-west-1", Service: "S3", Date: costDate, Cost: "2.4474", AccountID: "003A"},
		// account 004A
		{Region: "us-west-1", Service: "S3", Date: costDate, Cost: "102.7409", AccountID: "004A"},
	}
}

// Seed assumes the database already exists and the inserts pre-determined data
// into the database via the import
func Seed(ctx context.Context, log *slog.Logger, db *sqlx.DB) (statements []*dbstmts.Insert[*infracostmodels.Cost, int], err error) {
	var lg *slog.Logger = log.With("func", "domain.infracosts.infracostseeds.Seed")

	lg.Debug("starting ...")
	statements, err = infracostimports.Import(ctx, log, db, seeds)
	if err != nil {
		lg.Error("error with seed import.", "err", err.Error())
		err = errors.Join(ErrSeedImportFailed, err)
		return
	}
	lg.Debug("complete.")
	return

}
