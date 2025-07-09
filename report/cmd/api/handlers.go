package main

import (
	"context"
	"log/slog"

	"github.com/danielgtaylor/huma/v2"
	"github.com/ministryofjustice/opg-reports/report/cmd/api/awsaccounts"
	"github.com/ministryofjustice/opg-reports/report/cmd/api/awscosts"
	"github.com/ministryofjustice/opg-reports/report/cmd/api/home"
	"github.com/ministryofjustice/opg-reports/report/cmd/api/teams"
	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/sqlr"
	"github.com/ministryofjustice/opg-reports/report/internal/service/api"
)

// RegisterHandlers attaches all the known functions to the api.
//
// To allow for service injection, each is called directly, so need to be manually added
func RegisterHandlers(ctx context.Context, log *slog.Logger, conf *config.Config, humaapi huma.API) {
	var (
		teamStore          = sqlr.DefaultWithSelect[*api.Team](ctx, log, conf)
		teamService        = api.Default[*api.Team](ctx, log, conf)
		awsAccountsStore   = sqlr.DefaultWithSelect[*api.AwsAccount](ctx, log, conf)
		awsAccountService  = api.Default[*api.AwsAccount](ctx, log, conf)
		awsCostsStore      = sqlr.DefaultWithSelect[*api.AwsCost](ctx, log, conf)
		awsCostsService    = api.Default[*api.AwsCost](ctx, log, conf)
		awsCostsStoreGroup = sqlr.DefaultWithSelect[*api.AwsCostGrouped](ctx, log, conf)
		awsCostsSrvGroup   = api.Default[*api.AwsCostGrouped](ctx, log, conf)
	)
	// HOME
	home.RegisterGetHomepage(log, conf, humaapi)
	// TEAMS
	teams.RegisterGetTeamsAll(log, conf, humaapi, teamService, teamStore)
	// AWS ACCOUNTS
	awsaccounts.RegisterGetAwsAccountsAll(log, conf, humaapi, awsAccountService, awsAccountsStore)
	// AWS COSTS
	awscosts.RegisterGetAwsCostsTop20(log, conf, humaapi, awsCostsService, awsCostsStore)
	awscosts.RegisterGetAwsGroupedCosts(log, conf, humaapi, awsCostsSrvGroup, awsCostsStoreGroup)
}
