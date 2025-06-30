package teams

import (
	"context"
	"log/slog"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/interfaces"
	"github.com/ministryofjustice/opg-reports/report/internal/service/team"
)

// Service is a small helper that fetches the service for team.Team related calls
// and returns that
func Service[T interfaces.Model](ctx context.Context, log *slog.Logger, conf *config.Config) (srv *team.Service[T]) {
	srv = team.Default[T](ctx, log, conf)
	return
}
