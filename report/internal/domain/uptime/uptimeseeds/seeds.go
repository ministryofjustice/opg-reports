package uptimeseeds

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/accounts/accountseeds"
	"opg-reports/report/internal/domain/uptime/uptimeimports"
	"opg-reports/report/internal/domain/uptime/uptimemodels"
	"opg-reports/report/internal/utils/times"
	"time"

	"github.com/jmoiron/sqlx"
)

var ErrSeedFailed = errors.New("seed uptime call failed with an error.")

// GetSeeds returns the raw seeds to use for this package
// Should generate 76k entries
func GetSeeds() (data []*uptimemodels.Uptime) {
	var (
		accounts  = accountseeds.GetSeeds()
		startdate = times.ResetYear(times.Add(time.Now(), -2, times.YEAR))
		enddate   = times.ResetMonth(times.Add(time.Now(), -1, times.MONTH))
		days      = times.Days(startdate, enddate)
	)
	data = []*uptimemodels.Uptime{}
	// generate large number of random averages per day
	for _, day := range days {
		for _, acc := range accounts {
			var avg float64 = (95) + (rand.Float64() * (100 - 95)) // 95-100%
			data = append(data, &uptimemodels.Uptime{
				Date:      times.AsYMDString(day),
				AccountID: acc.ID,
				Average:   fmt.Sprintf("%g", avg),
			})
		}
	}
	return
}

// Seed assumes the database already exists and the inserts pre-determined data
// into the database via the import
func Seed(ctx context.Context, log *slog.Logger, db *sqlx.DB) (statements []*dbstmts.Insert[*uptimemodels.Uptime, int], err error) {
	var seeds []*uptimemodels.Uptime = GetSeeds()
	var lg *slog.Logger = log.With("func", "domain.uptime.uptimeseeds.Seed")

	lg.Debug("starting ...")
	statements, err = uptimeimports.Import(ctx, log, db, seeds)
	if err != nil {
		lg.Error("error with seed import", "err", err.Error())
		err = errors.Join(ErrSeedFailed, err)
		return
	}
	lg.With("count", len(statements)).Debug("complete.")
	return

}
