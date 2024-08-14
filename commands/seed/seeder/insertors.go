package seeder

import (
	"context"
	"database/sql"
	"log/slog"
	"sync"

	"github.com/ministryofjustice/opg-reports/datastore/github_standards/ghs"
	"github.com/ministryofjustice/opg-reports/shared/convert"
	"github.com/ministryofjustice/opg-reports/shared/testhelpers"
)

type insertF func(ctx context.Context, fileContent []byte, db *sql.DB) error

var inserts map[string]insertF = map[string]insertF{
	"github_standards": func(ctx context.Context, fileContent []byte, db *sql.DB) (err error) {
		// unmarshal that content
		repos := []*ghs.GithubStandard{}
		repos, err = convert.Unmarshals[*ghs.GithubStandard](fileContent)
		if err != nil {
			return
		}
		// -- setup the transaction
		mu := &sync.Mutex{}
		wg := sync.WaitGroup{}
		tx, err := db.Begin()
		if err != nil {
			return
		}
		defer tx.Rollback()
		q := ghs.New(db)
		qtx := q.WithTx(tx)

		slog.Info("starting insertion", slog.String("for", "github_standards"), slog.Int("count", len(repos)))
		tick := testhelpers.T()
		for _, item := range repos {
			wg.Add(1)
			go func(g *ghs.GithubStandard) {
				mu.Lock()
				qtx.Insert(ctx, g.Insertable())
				mu.Unlock()
				wg.Done()
			}(item)
		}

		wg.Wait()
		slog.Info("generation complete", slog.String("seconds", tick.Stop().Seconds()), slog.String("for", "github_standards"))

		return tx.Commit()
	},
}
