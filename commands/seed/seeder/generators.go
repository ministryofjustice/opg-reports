package seeder

import (
	"context"
	"database/sql"
	"log/slog"
	"sync"

	"github.com/ministryofjustice/opg-reports/datastore/aws_costs/awsc"
	"github.com/ministryofjustice/opg-reports/datastore/aws_uptime/awsu"
	"github.com/ministryofjustice/opg-reports/datastore/github_standards/ghs"
	"github.com/ministryofjustice/opg-reports/shared/fake"
	"github.com/ministryofjustice/opg-reports/shared/testhelpers"
)

// map of funcs that inset data from files
var GENERATOR_FUNCTIONS map[string]generatorF = map[string]generatorF{
	"github_standards": githubStandardsGenerator,
	"aws_costs":        awsCostsGenerator,
	"aws_uptime":       awsUptimeGenerator,
}

func awsUptimeGenerator(ctx context.Context, num int, db *sql.DB) (err error) {
	mu := &sync.Mutex{}
	wg := sync.WaitGroup{}

	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()
	q := awsu.New(db)
	qtx := q.WithTx(tx)

	slog.Info("starting generation", slog.String("for", "aws_uptime"))
	tick := testhelpers.T()
	for x := 0; x < num; x++ {
		wg.Add(1)
		go func(i int) {
			mu.Lock()
			c := awsu.Fake()
			qtx.Insert(ctx, c.Insertable())
			mu.Unlock()
			wg.Done()
		}(x)
	}
	wg.Wait()
	slog.Info("generation complete",
		slog.Int("n", num),
		slog.String("seconds", tick.Stop().Seconds()),
		slog.String("for", "aws_uptime"))

	return tx.Commit()
}

func awsCostsGenerator(ctx context.Context, num int, db *sql.DB) (err error) {
	mu := &sync.Mutex{}
	wg := sync.WaitGroup{}

	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()
	q := awsc.New(db)
	qtx := q.WithTx(tx)

	slog.Info("starting generation", slog.String("for", "aws_costs"))
	tick := testhelpers.T()
	for x := 0; x < num; x++ {
		wg.Add(1)
		go func(i int) {
			mu.Lock()
			c := awsc.Fake()
			qtx.Insert(ctx, c.Insertable())
			mu.Unlock()
			wg.Done()
		}(x)
	}
	wg.Wait()
	slog.Info("generation complete", slog.Int("n", num), slog.String("seconds", tick.Stop().Seconds()), slog.String("for", "aws_costs"))

	return tx.Commit()
}

func githubStandardsGenerator(ctx context.Context, num int, db *sql.DB) (err error) {
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
	slog.Info("generation complete", slog.Int("n", num), slog.String("seconds", tick.Stop().Seconds()), slog.String("for", "github_standards"))

	return tx.Commit()
}
