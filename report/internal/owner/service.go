package owner

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ministryofjustice/opg-reports/report/internal/repository"
)

type Service struct {
	ctx   context.Context
	log   *slog.Logger
	store *repository.Repository
}

// Seed is used to insert baseline data for owners
func (self *Service) Seed() {

}

func (self *Service) GetAllOwners() (owners []*Owner) {
	owners = []*Owner{}
	return
}

func NewService(ctx context.Context, log *slog.Logger, store *repository.Repository) (srv *Service, err error) {
	if log == nil {
		return nil, fmt.Errorf("no logger passed for owner service")
	}
	if store == nil {
		return nil, fmt.Errorf("no repository passed for owner service")
	}

	srv = &Service{
		ctx:   ctx,
		log:   log.With("service", "owner"),
		store: store,
	}
	return
}
