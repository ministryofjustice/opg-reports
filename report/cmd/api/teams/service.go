package teams

import (
	"context"
	"log/slog"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/team"
)

// Service is a small helper that fetches the service for team.Team related calls
// and returns that
func Service(ctx context.Context, log *slog.Logger, conf *config.Config) (srv *team.Service[*team.Team], err error) {

	srv, err = team.Default[*team.Team](ctx, log, conf)
	return
}
