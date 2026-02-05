package uptimeseeds

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/uptime/uptimeimports"
	"opg-reports/report/internal/domain/uptime/uptimemodels"
	"opg-reports/report/internal/utils/times"
	"time"

	"github.com/jmoiron/sqlx"
)

var ErrSeedFailed = errors.New("seed uptime call failed with an error.")

var (
	date  string
	seeds []*uptimemodels.Uptime
)

func init() {
	date = times.AsYMDString(time.Now())
	seeds = []*uptimemodels.Uptime{
		{Date: date, Average: "99.9901", AccountID: "001A"},
		{Date: date, Average: "99.9801", AccountID: "001B"},
		{Date: date, Average: "99.9801", AccountID: "001C"},
	}
}

// Seed assumes the database already exists and the inserts pre-determined data
// into the database via the import
func Seed(ctx context.Context, log *slog.Logger, db *sqlx.DB) (statements []*dbstmts.Insert[*uptimemodels.Uptime, int], err error) {
	var lg *slog.Logger = log.With("func", "domain.uptime.uptimeseeds.Seed")

	lg.Debug("starting ...")
	statements, err = uptimeimports.Import(ctx, log, db, seeds)
	if err != nil {
		lg.Error("error with seed import", "err", err.Error())
		err = errors.Join(ErrSeedFailed, err)
		return
	}
	lg.Debug("complete.")
	return

}
