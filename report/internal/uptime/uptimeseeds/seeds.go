package uptimeseeds

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/db/dbstatements"
	"opg-reports/report/internal/uptime/uptimeimports"
	"opg-reports/report/internal/uptime/uptimemodels"
	"opg-reports/report/internal/utils/times"
	"time"

	"github.com/jmoiron/sqlx"
)

var (
	date  string
	seeds []*uptimemodels.Uptime
)

func init() {
	date = times.AsYMDString(time.Now())
	seeds = []*uptimemodels.Uptime{
		{Date: date, Average: "99.9901", AccountID: "001A"},
		{Date: date, Average: "99.9801", AccountID: "001A"},
		{Date: date, Average: "99.9801", AccountID: "001A"},
	}
}

// Seed assumes the database already exists and the inserts pre-determined data
// into the database via the import
func Seed(ctx context.Context, log *slog.Logger, db *sqlx.DB) (statements []*dbstatements.DataStatement[*uptimemodels.Uptime, int], err error) {

	log = log.With("package", "uptime", "func", "Seed")
	log.Debug("starting ...")

	statements, err = uptimeimports.Import(ctx, log, db, seeds)
	if err != nil {
		log.Error("error with seed import", "err", err.Error())
		err = errors.Join(ErrSeedFailed, err)
		return
	}
	log.Debug("complete")
	return

}
