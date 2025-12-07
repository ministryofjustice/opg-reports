package main

import (
	"context"
	"log/slog"

	"opg-reports/report/cmd/api/awsaccounts"
	"opg-reports/report/cmd/api/awscosts"
	"opg-reports/report/cmd/api/awsuptime"
	"opg-reports/report/cmd/api/githubcodeowners"
	"opg-reports/report/cmd/api/home"
	"opg-reports/report/cmd/api/teams"
	"opg-reports/report/config"
	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/service/api"

	"github.com/danielgtaylor/huma/v2"
)

// RegisterHandlers attaches all the known functions to the api.
//
// To allow for service injection, each is called directly, so need to be manually added
func RegisterHandlers(ctx context.Context, log *slog.Logger, conf *config.Config, humaapi huma.API) {
	var (
		// TEAMS
		teamStore   = sqlr.DefaultWithSelect[*api.Team](ctx, log, conf)
		teamService = api.Default[*api.Team](ctx, log, conf)
		// ACCOUNTS
		awsAccountsStore  = sqlr.DefaultWithSelect[*api.AwsAccount](ctx, log, conf)
		awsAccountService = api.Default[*api.AwsAccount](ctx, log, conf)
		// COSTS
		awsCostsStore          = sqlr.DefaultWithSelect[*api.AwsCost](ctx, log, conf)
		awsCostsService        = api.Default[*api.AwsCost](ctx, log, conf)
		awsCostsGroupedStore   = sqlr.DefaultWithSelect[*api.AwsCostGrouped](ctx, log, conf)
		awsCostsGroupedService = api.Default[*api.AwsCostGrouped](ctx, log, conf)
		// UPTIME
		awsUptimeStore          = sqlr.DefaultWithSelect[*api.AwsUptime](ctx, log, conf)
		awsUptimeService        = api.Default[*api.AwsUptime](ctx, log, conf)
		awsUptimeGroupedStore   = sqlr.DefaultWithSelect[*api.AwsUptimeGrouped](ctx, log, conf)
		awsUptimeGroupedService = api.Default[*api.AwsUptimeGrouped](ctx, log, conf)
		// GITHUB CODEOWNERS
		githubCodeOwnerStore   = sqlr.DefaultWithSelect[*api.GithubCodeOwner](ctx, log, conf)
		githubCodeOwnerService = api.Default[*api.GithubCodeOwner](ctx, log, conf)
	)
	// HOME
	home.RegisterGetHomepage(log, conf, humaapi)
	// TEAMS
	teams.RegisterGetTeamsAll(log, conf, humaapi, teamService, teamStore)
	// AWS ACCOUNTS
	awsaccounts.RegisterGetAwsAccountsAll(log, conf, humaapi, awsAccountService, awsAccountsStore)
	// AWS COSTS
	awscosts.RegisterGetAwsCostsTop20(log, conf, humaapi, awsCostsService, awsCostsStore)
	awscosts.RegisterGetAwsCostsGrouped(log, conf, humaapi, awsCostsGroupedService, awsCostsGroupedStore)
	// AWS UPTIME
	awsuptime.RegisterGetAwsUptimeAll(log, conf, humaapi, awsUptimeService, awsUptimeStore)
	awsuptime.RegisterGetAwsUptimeGrouped(log, conf, humaapi, awsUptimeGroupedService, awsUptimeGroupedStore)
	// GITHUB CODEOWNERS
	githubcodeowners.RegisterGetGithubCodeOwnersAll(log, conf, humaapi, githubCodeOwnerService, githubCodeOwnerStore)
	githubcodeowners.RegisterGetGithubCodeOwnersForTeam(log, conf, humaapi, githubCodeOwnerService, githubCodeOwnerStore)

}
