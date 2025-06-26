package awscost

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

// GetAll will return all records
//
// Using this is generally a bad idea as this table will contain millions of rows
func (self *Service[T]) GetAll() (accounts []T, err error) {
	var selectStmt = &sqldb.BoundStatement{Statement: stmtSelectAll}
	var log = self.log.With("operation", "GetAll")

	accounts = []T{}
	log.Debug("getting all awscosts from database...")

	// cast the data back to struct
	if err = self.store.Select(selectStmt); err == nil {
		accounts = selectStmt.Returned.([]T)
	}

	return
}

// NewService creates a service using the values passed
func NewService[T interfaces.Model](ctx context.Context, log *slog.Logger, conf *config.Config, store *sqldb.Repository[T]) (srv *Service[T], err error) {
	srv = &Service[T]{}
	if log == nil {
		err = fmt.Errorf("no logger passed for awscost service")
		return
	}
	if conf == nil {
		err = fmt.Errorf("no config passed for awscost service")
		return
	}
	if store == nil {
		err = fmt.Errorf("no repository passed for awscost service")
		return
	}

	srv = &Service[T]{
		ctx:   ctx,
		log:   log.With("service", "awscost"),
		conf:  conf,
		store: store,
	}
	return
}

// Default generates the default repository as and then the service
func Default[T interfaces.Model](ctx context.Context, log *slog.Logger, conf *config.Config) (srv *Service[T], err error) {

	store, err := sqldb.New[T](ctx, log, conf)
	if err != nil {
		return
	}
	srv, err = NewService[T](ctx, log, conf, store)

	return
}
