package seed

import (
	"context"
	"fmt"
	"log/slog"

	"opg-reports/report/config"
	"opg-reports/report/internal/repository/sqlr"
)

const label string = "seed-service"

// SeedAllResults contains all the results from each step of the all function
type SeedAllResults struct {
	Teams            []*sqlr.BoundStatement
	AwsAccounts      []*sqlr.BoundStatement
	AwsCosts         []*sqlr.BoundStatement
	AwsUptime        []*sqlr.BoundStatement
	GithubCodeOwners []*sqlr.BoundStatement
}

type Service struct {
	ctx  context.Context
	log  *slog.Logger
	conf *config.Config
}

// All runs all seed commands in sequence
func (self *Service) All(sqlStore sqlr.RepositoryWriter) (results *SeedAllResults, err error) {
	// setup the results
	results = &SeedAllResults{}
	// TEAMS
	results.Teams, err = self.Teams(sqlStore)
	if err != nil {
		return
	}
	// AWS ACCOUNTS
	results.AwsAccounts, err = self.AwsAccounts(sqlStore)
	if err != nil {
		return
	}
	// AWS COSTS
	results.AwsCosts, err = self.AwsCosts(sqlStore)
	if err != nil {
		return
	}
	// AWS UPTIME
	results.AwsUptime, err = self.AwsUptime(sqlStore)
	if err != nil {
		return
	}
	// GITHUB CODE OWNERS
	results.GithubCodeOwners, err = self.GithubCodeOwners(sqlStore)
	if err != nil {
		return
	}

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
