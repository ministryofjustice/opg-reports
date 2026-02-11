package infracostseeds

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/accounts/accountseeds"
	"opg-reports/report/internal/domain/infracosts/infracostimports"
	"opg-reports/report/internal/domain/infracosts/infracostmodels"
	"opg-reports/report/internal/utils/times"
	"time"

	"github.com/jmoiron/sqlx"
)

var ErrSeedImportFailed = errors.New("seed costs call failed with an error.")

// GetSeeds generates a large amount of cost data to allow better testing
func GetSeeds() (data []*infracostmodels.Cost) {
	var (
		accounts  = accountseeds.GetSeeds()
		startdate = times.ResetYear(times.Add(time.Now(), -2, times.YEAR))
		enddate   = times.ResetMonth(times.Add(time.Now(), -1, times.MONTH))
		months    = times.Months(startdate, enddate)
		regions   = []string{"eu-west-1", "eu-west-2", "NoRegion"}
		services  = []string{
			"Amazon Relational Database Service",
			"Amazon Simple Storage Service",
			"AmazonCloudWatch",
			"Amazon Elastic Load Balancing",
			"AWS Shield",
			"AWS Config",
			"AWS CloudTrail",
			"AWS Key Management Service",
			"Amazon Virtual Private Cloud",
			"Amazon Elastic Container Service",
		}
	)
	data = []*infracostmodels.Cost{}
	for _, month := range months {
		for _, account := range accounts {
			for _, service := range services {
				for _, region := range regions {
					var price float64 = (-1000.0) + (rand.Float64() * (1000 - -1000.0)) // 95-100%
					data = append(data, &infracostmodels.Cost{
						Region:    region,
						Service:   service,
						Date:      times.AsYMDString(month),
						Cost:      fmt.Sprintf("%g", price),
						AccountID: account.ID,
					})
				}
			}
		}
	}

	return
}

// Seed assumes the database already exists and the inserts pre-determined data
// into the database via the import
func Seed(ctx context.Context, log *slog.Logger, db *sqlx.DB) (statements []*dbstmts.Insert[*infracostmodels.Cost, int], err error) {
	var lg *slog.Logger = log.With("func", "infracostseeds.Seed")
	var seeds []*infracostmodels.Cost = GetSeeds()
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
