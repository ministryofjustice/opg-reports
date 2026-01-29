package accountselects

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/db/dbselects"
	"opg-reports/report/internal/db/dbstatements"
	"opg-reports/report/internal/domain/accounts/accountmodels"

	"github.com/jmoiron/sqlx"
)

const selectAllStmt string = `
SELECT
	id,
	name,
	label,
	environment,
	team_name
FROM accounts
ORDER BY name, environment ASC
`

// empty is used when no filtering / substitution on the sql statement, such
// as a `select *` or `select count(*)`
type empty struct{}

func All(ctx context.Context, log *slog.Logger, db *sqlx.DB) (data []*accountmodels.Account, err error) {
	var (
		lg       *slog.Logger = log.With("func", "domain.accounts.accountselects.All")
		selector *dbstatements.SelectStatement[*empty, *accountmodels.Account]
	)
	data = []*accountmodels.Account{}
	// setup the select
	selector = &dbstatements.SelectStatement[*empty, *accountmodels.Account]{
		Statement: selectAllStmt,
		Data:      &empty{},
	}

	lg.Debug("starting ...")
	err = dbselects.Select(ctx, log, db, selector)
	if err != nil {
		lg.Error("error with select.", "err", err.Error())
		return
	}
	// setup output
	for _, row := range selector.Returned {
		data = append(data, row)
	}
	lg.Debug("complete.")
	return
}
