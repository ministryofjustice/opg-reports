package seeder

import (
	"context"
	"database/sql"
	"time"

	"github.com/ministryofjustice/opg-reports/datastore/aws_costs/awsc"
	"github.com/ministryofjustice/opg-reports/datastore/github_standards/ghs"
	"github.com/ministryofjustice/opg-reports/shared/dates"
)

var TRACKER_FUNCTIONS map[string]trackerF = map[string]trackerF{
	"github_standards": func(ctx context.Context, ts time.Time, db *sql.DB) (err error) {
		q := ghs.New(db)
		err = q.Track(ctx, ts.Format(dates.Format))
		return
	},
	"aws_costs": func(ctx context.Context, ts time.Time, db *sql.DB) (err error) {
		q := awsc.New(db)
		err = q.Track(ctx, ts.Format(dates.Format))
		return
	},
}
