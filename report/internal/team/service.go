package team

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

// Seed is used to insert baseline data for teams
// TODO: set up base seed with units in data
func (self *Service[T]) Seed() (err error) {
	var (
		now   = time.Now().UTC().Format(time.RFC3339)
		seeds = []*repository.BoundStatement{
			{Statement: stmtInsert, Data: &Team{Name: "Sirius", CreatedAt: now}},
		}
	)
	err = self.store.Insert(seeds...)
	return
}

// GetAllTeams returns all teams as a slice from the database
func (self *Service[T]) GetAllTeams() (teams []*Team) {
	teams = []*Team{}
	return
}

func NewService[T interfaces.Model](ctx context.Context, log *slog.Logger, store *repository.Repository[T]) (srv *Service[T], err error) {
	if log == nil {
		return nil, fmt.Errorf("no logger passed for team service")
	}
	if store == nil {
		return nil, fmt.Errorf("no repository passed for team service")
	}

	srv = &Service[T]{
		ctx:   ctx,
		log:   log.With("service", "team"),
		store: store,
	}
	return
}
