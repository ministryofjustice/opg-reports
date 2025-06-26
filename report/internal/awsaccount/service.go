package awsaccount

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

// GetAllAccounts returns all accounts as a slice from the database
func (self *Service[T]) GetAllAccounts() (teams []T, err error) {
	var selectStmt = &sqldb.BoundStatement{Statement: stmtSelectAll}
	var log = self.log.With("operation", "GetAllAccounts")

	teams = []T{}
	log.Debug("getting all awsaccounts from database...")

	if err = self.store.Select(selectStmt); err == nil {
		// cast the data back to struct
		teams = selectStmt.Returned.([]T)
	}

	return
}

// NewService creates a service using the values passed
func NewService[T interfaces.Model](ctx context.Context, log *slog.Logger, conf *config.Config, store *sqldb.Repository[T]) (srv *Service[T], err error) {
	if log == nil {
		return nil, fmt.Errorf("no logger passed for awsaccount service")
	}
	if conf == nil {
		return nil, fmt.Errorf("no config passed for awsaccount service")
	}
	if store == nil {
		return nil, fmt.Errorf("no repository passed for awsaccount service")
	}

	srv = &Service[T]{
		ctx:   ctx,
		log:   log.With("service", "awsaccount"),
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
