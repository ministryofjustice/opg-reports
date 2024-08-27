package seeder

import (
	"context"
	"database/sql"
	"log/slog"
	"sync"

	"github.com/ministryofjustice/opg-reports/datastore/aws_costs/awsc"
	"github.com/ministryofjustice/opg-reports/datastore/github_standards/ghs"
	"github.com/ministryofjustice/opg-reports/shared/convert"
	"github.com/ministryofjustice/opg-reports/shared/testhelpers"
)

var INSERT_FUNCTIONS map[string]insertF = map[string]insertF{
	"github_standards": func(ctx context.Context, fileContent []byte, db *sql.DB) (err error) {
		// unmarshal that content
		which := "github_standards"
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

		slog.Debug("inserting content", slog.String("for", which), slog.Int("count", len(repos)))
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
		slog.Debug("insert complete",
			slog.Int("count", len(repos)),
			slog.String("seconds", tick.Stop().Seconds()),
			slog.String("for", which))

		return tx.Commit()
	},

	"aws_costs": func(ctx context.Context, fileContent []byte, db *sql.DB) (err error) {
		// unmarshal that content
		which := "aws_costs"
		items := []*awsc.AwsCost{}
		items, err = convert.Unmarshals[*awsc.AwsCost](fileContent)
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
		q := awsc.New(db)
		qtx := q.WithTx(tx)

		slog.Debug("inserting content", slog.String("for", which), slog.Int("count", len(items)))
		tick := testhelpers.T()
		for _, item := range items {
			wg.Add(1)
			go func(g *awsc.AwsCost) {
				mu.Lock()
				qtx.Insert(ctx, g.Insertable())
				mu.Unlock()
				wg.Done()
			}(item)
		}

		wg.Wait()
		slog.Debug("insert complete",
			slog.Int("count", len(items)),
			slog.String("seconds", tick.Stop().Seconds()),
			slog.String("for", which))

		return tx.Commit()
	},
}
