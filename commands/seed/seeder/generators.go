package seeder

import (
	"context"
	"database/sql"
	"log/slog"
	"sync"

	"github.com/ministryofjustice/opg-reports/datastore/github_standards/ghs"
	"github.com/ministryofjustice/opg-reports/shared/fake"
	"github.com/ministryofjustice/opg-reports/shared/testhelpers"
)

type generatorF func(ctx context.Context, num int, db *sql.DB) error

// map of funcs that inset data from files
var GENERATOR_FUNCTIONS map[string]generatorF = map[string]generatorF{
	// -- githu standards
	// use go concurrency for speed
	"github_standards": func(ctx context.Context, num int, db *sql.DB) (err error) {
		mu := &sync.Mutex{}
		wg := sync.WaitGroup{}

		owner := fake.String(12)
		tx, err := db.Begin()
		if err != nil {
			return
		}
		defer tx.Rollback()
		q := ghs.New(db)
		qtx := q.WithTx(tx)

		slog.Info("starting generation", slog.String("for", "github_standards"))
		tick := testhelpers.T()
		for x := 0; x < num; x++ {
			wg.Add(1)
			go func(i int) {
				mu.Lock()
				g := ghs.Fake(nil, &owner)
				qtx.Insert(ctx, g.Insertable())
				mu.Unlock()
				wg.Done()
			}(x)
		}

		wg.Wait()
		slog.Info("generation complete", slog.String("seconds", tick.Stop().Seconds()), slog.String("for", "github_standards"))

		return tx.Commit()
	},
}
