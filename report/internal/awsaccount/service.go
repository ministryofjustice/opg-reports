package awsaccount

import (
	"context"
	"fmt"
	"log/slog"
	"time"

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

// Seed is used to insert test data to the table, so for now we create 3 dummy versions
func (self *Service[T]) Seed() (err error) {
	var (
		log   = self.log.With("operation", "Seed")
		now   = time.Now().UTC().Format(time.RFC3339)
		seeds = []*sqldb.BoundStatement{
			{Statement: stmtInsert, Data: &AwsAccount{ID: "001A", Name: "Acc01A", Label: "A", Environment: "development", CreatedAt: now}},
			{Statement: stmtInsert, Data: &AwsAccount{ID: "001B", Name: "Acc01B", Label: "B", Environment: "production", CreatedAt: now}},
			{Statement: stmtInsert, Data: &AwsAccount{ID: "002A", Name: "Acc02A", Label: "A", Environment: "production", CreatedAt: now}},
		}
	)
	log.Info("inserting seed data ...")
	err = self.store.Insert(seeds...)

	return
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
