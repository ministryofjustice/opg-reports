package component

import (
	"context"
	"log/slog"
	"opg-reports/report/config"
	"opg-reports/report/internal/repository/restr"
	"opg-reports/report/internal/service/front"
)

type Parser[T front.Response, R front.Result] func(response T) (R, error)

type Component[T front.Response, R front.Result] struct {
	ctx  context.Context
	log  *slog.Logger
	conf *config.Config

	Parser Parser[T, R]
}

// Call handles fetching the data from the API, so Component is tightly coupled
// to the front module, more than would like!
func (self *Component[T, R]) Call(client restr.RepositoryRestGetter, endpoint string) (R, error) {
	var service = front.Default[T, R](self.ctx, self.log, self.conf)
	self.log.With("endpoint", endpoint).Info("calling api ... ")
	return service.GetFromAPI(client, endpoint, self.Parser)

}

func New[T front.Response, R front.Result](
	ctx context.Context,
	log *slog.Logger,
	conf *config.Config,
	parser Parser[T, R],
) *Component[T, R] {
	return &Component[T, R]{
		ctx:    ctx,
		log:    log,
		conf:   conf,
		Parser: parser,
	}
}
