package seed

import (
	"context"
	"fmt"
	"log/slog"

	"opg-reports/report/config"
	"opg-reports/report/internal/repository/sqlr"
)

const label string = "seed-service"

type Service struct {
	ctx  context.Context
	log  *slog.Logger
	conf *config.Config
}

// All runs all seed commands in sequence
func (self *Service) All(sqlStore sqlr.RepositoryWriter) (results []*sqlr.BoundStatement, err error) {
	var r []*sqlr.BoundStatement
	// TEAMS
	r, err = self.Teams(sqlStore)
	if err != nil {
		return
	}
	results = append(results, r...)
	// AWS ACCOUNTS
	r, err = self.AwsAccounts(sqlStore)
	if err != nil {
		return
	}
	results = append(results, r...)
	// AWS COSTS
	r, err = self.AwsCosts(sqlStore)
	if err != nil {
		return
	}
	results = append(results, r...)
	// AWS UPTIME
	r, err = self.AwsUptime(sqlStore)
	if err != nil {
		return
	}
	results = append(results, r...)

	return
}

// New tries to create a version of this service using the context, logger and config values
// passed along.
//
// If logger / config is not passed an error is returned
func New(ctx context.Context, log *slog.Logger, conf *config.Config) (srv *Service, err error) {
	if log == nil {
		return nil, fmt.Errorf("no logger passed for %s", label)
	}
	if conf == nil {
		return nil, fmt.Errorf("no config passed for %s", label)
	}

	srv = &Service{
		ctx:  ctx,
		log:  log.With("service", label),
		conf: conf,
	}
	return
}

// Default creates a service by calling `New` and swallowing any errors.
//
// Errors are logged, but not shared
func Default(ctx context.Context, log *slog.Logger, conf *config.Config) (srv *Service) {
	srv, err := New(ctx, log, conf)
	if err != nil {
		log.Error("error with default", "err", err.Error())
	}
	return
}
