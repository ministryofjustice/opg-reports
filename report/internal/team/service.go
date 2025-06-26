package team

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/interfaces"
	"github.com/ministryofjustice/opg-reports/report/internal/sqldb"
)

type Service[T interfaces.Model] struct {
	ctx   context.Context
	log   *slog.Logger
	conf  *config.Config
	store *sqldb.Repository[T]
}

// Close function to do any clean up
func (self *Service[T]) Close() (err error) {
	return
}

// GetAllTeams returns all teams as a slice from the database
// Calls the database
func (self *Service[T]) GetAllTeams() (teams []T, err error) {
	var selectStmt = &sqldb.BoundStatement{Statement: stmtSelectAll}
	var log = self.log.With("operation", "GetAllTeams")

	teams = []T{}
	log.Debug("getting all teams from database...")

	if err = self.store.Select(selectStmt); err == nil {
		// cast the data back to struct
		teams = selectStmt.Returned.([]T)
	}

	return
}

func NewService[T interfaces.Model](ctx context.Context, log *slog.Logger, conf *config.Config, store *sqldb.Repository[T]) (srv *Service[T], err error) {
	srv = &Service[T]{}

	if log == nil {
		err = fmt.Errorf("no logger passed for team service")
		return
	}
	if conf == nil {
		err = fmt.Errorf("no config passed for team service")
		return
	}
	if store == nil {
		err = fmt.Errorf("no repository passed for team service")
		return
	}

	srv = &Service[T]{
		ctx:   ctx,
		log:   log.With("service", "team"),
		conf:  conf,
		store: store,
	}
	return
}

// Default generates the default repository and then the service
func Default[T interfaces.Model](ctx context.Context, log *slog.Logger, conf *config.Config) (srv *Service[T], err error) {

	store, err := sqldb.New[T](ctx, log, conf)
	if err != nil {
		return
	}
	srv, err = NewService[T](ctx, log, conf, store)

	return
}
