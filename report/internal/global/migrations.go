package global

import (
	"context"
	"opg-reports/report/internal/cost/costmigrations"
	"opg-reports/report/internal/global/migrations"
	"opg-reports/report/internal/team/teammigrations"

	_ "github.com/mattn/go-sqlite3"
)

// MigrateAll is a wrapper around migrating all known migrations
func MigrateAll(ctx context.Context, flags *migrations.Args) (err error) {

	if err = migrations.Run(ctx, flags, costmigrations.Migrations); err != nil {
		return
	}

	if err = migrations.Run(ctx, flags, teammigrations.Migrations); err != nil {
		return
	}

	return
}
