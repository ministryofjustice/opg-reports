package owner

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ministryofjustice/opg-reports/report/internal/interfaces"
	"github.com/ministryofjustice/opg-reports/report/internal/repository"
)

type Service[T interfaces.Model] struct {
	ctx   context.Context
	log   *slog.Logger
	store *repository.Repository[T]
}

// Seed is used to insert baseline data for owners
// TODO: set up base seed with units in data
func (self *Service[T]) Seed() (err error) {
	var (
		now   = time.Now().UTC().Format(time.RFC3339)
		seeds = []*repository.BoundStatement{
			{Statement: stmtInsert, Data: &Owner{Name: "Sirius", CreatedAt: now}},
		}
	)
	err = self.store.Insert(seeds...)
	return
}

func (self *Service[T]) GetAllOwners() (owners []*Owner) {
	owners = []*Owner{}
	return
}

func NewService[T interfaces.Model](ctx context.Context, log *slog.Logger, store *repository.Repository[T]) (srv *Service[T], err error) {
	if log == nil {
		return nil, fmt.Errorf("no logger passed for owner service")
	}
	if store == nil {
		return nil, fmt.Errorf("no repository passed for owner service")
	}

	srv = &Service[T]{
		ctx:   ctx,
		log:   log.With("service", "owner"),
		store: store,
	}
	return
}
