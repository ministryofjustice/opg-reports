package teams

import (
	"context"
	"log/slog"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/sqldb"
	"github.com/ministryofjustice/opg-reports/report/internal/team"
)

// Service is a small helper that fetches the service for team.Team related calls
// and returns that
func Service(ctx context.Context, log *slog.Logger, conf *config.Config) (srv *team.Service[*team.Team], err error) {
	var datastore *sqldb.Repository[*team.Team]

	datastore, err = sqldb.New[*team.Team](ctx, log, conf)
	if err != nil {
		return
	}

	srv, err = team.NewService(ctx, log, conf, datastore)
	return
}
