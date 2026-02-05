package accountseeds

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/accounts/accountimports"
	"opg-reports/report/internal/domain/accounts/accountmodels"

	"github.com/jmoiron/sqlx"
)

var ErrSeedFailed = errors.New("seed account call failed with an error.")

// GetSeeds will generate 96 accounts (24 * 4 =A-Z with 4 environments)
func GetSeeds() (data []*accountmodels.Account) {
	var limit = 24 // so we dont loop around the captial letters
	var envs map[string]string = map[string]string{"A": "development", "B": "preproduction", "C": "integration", "D": "production"}
	var ch rune = 'B'
	data = []*accountmodels.Account{}

	// generate 180 accounts for testing purposes
	for i := 1; i <= limit; i++ {
		var prefix string = fmt.Sprintf("%03d", i)
		var team string = fmt.Sprintf("TEAM-%s", string(ch))
		for _, letter := range []string{"A", "B", "C", "D"} {
			data = append(data, &accountmodels.Account{
				ID:          fmt.Sprintf("%s%s", prefix, letter),
				Name:        fmt.Sprintf("Account %s%s", prefix, letter),
				Label:       letter,
				Environment: envs[letter],
				TeamName:    team,
			})
		}
		ch = getNextRune(ch)
	}
	return
}

// Seed assumes the database already exists and the inserts pre-determined data
// into the database via the import
func Seed(ctx context.Context, log *slog.Logger, db *sqlx.DB) (statements []*dbstmts.Insert[*accountmodels.Account, string], err error) {
	var seeds []*accountmodels.Account = GetSeeds()
	var lg *slog.Logger = log.With("func", "domain.accounts.accountseeds.Seed")

	lg.Debug("starting ...")
	statements, err = accountimports.Import(ctx, log, db, seeds)
	if err != nil {
		lg.Error("error with seed import.", "err", err.Error())
		err = errors.Join(ErrSeedFailed, err)
		return
	}
	lg.Debug("complete.")
	return

}

func getNextRune(ch rune) rune {
	return (ch+1-'A')%('Z'-'A') + 'A'
}
