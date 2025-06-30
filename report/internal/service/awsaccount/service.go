package awsaccount

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/interfaces"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/sqlr"
)

type Service[T interfaces.Model] struct {
	ctx   context.Context
	log   *slog.Logger
	conf  *config.Config
	store *sqlr.RepositoryWithSelect[T]
}

// Close function to do any clean up
func (self *Service[T]) Close() (err error) {
	return
}

// GetAllAccounts returns all accounts as a slice from the database
func (self *Service[T]) GetAllAccounts() (accounts []T, err error) {
	var selectStmt = &sqlr.BoundStatement{Statement: stmtSelectAll}
	var log = self.log.With("operation", "GetAllAccounts")

	accounts = []T{}
	log.Debug("getting all awsaccounts from database...")

	// cast the data back to struct
	if err = self.store.Select(selectStmt); err == nil {
		accounts = selectStmt.Returned.([]T)
	}

	return
}

// NewService creates a service using the values passed
func NewService[T interfaces.Model](ctx context.Context, log *slog.Logger, conf *config.Config, store *sqlr.RepositoryWithSelect[T]) (srv *Service[T], err error) {
	srv = &Service[T]{}
	if log == nil {
		err = fmt.Errorf("no logger passed for awsaccount service")
		return
	}
	if conf == nil {
		err = fmt.Errorf("no config passed for awsaccount service")
		return
	}
	if store == nil {
		err = fmt.Errorf("no repository passed for awsaccount service")
		return
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
//
// If there is an error creating the service, then nil is returned
func Default[T interfaces.Model](ctx context.Context, log *slog.Logger, conf *config.Config) (srv *Service[T]) {

	store, err := sqlr.NewWithSelect[T](ctx, log, conf)
	if err != nil {
		log.Error("error creating sqlr repository", "error", err.Error())
		return nil
	}
	srv, err = NewService[T](ctx, log, conf, store)
	if err != nil {
		log.Error("error creating awsaccount service", "error", err.Error())
		return nil
	}

	return
}
