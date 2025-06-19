package team

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ministryofjustice/opg-reports/report/internal/interfaces"
	"github.com/ministryofjustice/opg-reports/report/internal/sqldb"
)

type Service[T interfaces.Model] struct {
	ctx   context.Context
	log   *slog.Logger
	store *sqldb.Repository[T]
}

// Seed is used to insert test data to the team table
func (self *Service[T]) Import() (err error) {
	var (
		now   = time.Now().UTC().Format(time.RFC3339)
		seeds = []*sqldb.BoundStatement{
			{Statement: stmtInsert, Data: &Team{Name: "TeamA", CreatedAt: now}},
			{Statement: stmtInsert, Data: &Team{Name: "TeamB", CreatedAt: now}},
			{Statement: stmtInsert, Data: &Team{Name: "TeamC", CreatedAt: now}},
		}
	)
	err = self.store.Insert(seeds...)

	return
}

// GetAllTeams returns all teams as a slice from the database
// Calls the database
func (self *Service[T]) GetAllTeams() (teams []*Team, err error) {
	var selectStmt = &sqldb.BoundStatement{Statement: stmtSelectAll}
	teams = []*Team{}

	err = self.store.Select(selectStmt)
	// cast the data back to struct
	if err == nil {
		teams = selectStmt.Returned.([]*Team)
	}

	return
}

func NewService[T interfaces.Model](ctx context.Context, log *slog.Logger, store *sqldb.Repository[T]) (srv *Service[T], err error) {
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
