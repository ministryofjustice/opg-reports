package main

import (
	"context"
	"log/slog"
	"opg-reports/report/config"
	"opg-reports/report/internal/component"
)

type allComponents struct {
	TeamNavigation         *component.Component[*apiResponseTeams, []string]
	AwsCostsGroupedByMonth *component.Component[*apiResponseAwsCostsGrouped, []map[string]string]
}

// List of all components
var Components *allComponents

// initComponents creates all the components that the site uses for various parts of data
func initComponents(ctx context.Context, log *slog.Logger, conf *config.Config) {

	Components = &allComponents{
		TeamNavigation:         component.New(ctx, log, conf, parseAllTeamsForNavigation),
		AwsCostsGroupedByMonth: component.New(ctx, log, conf, parseAwsCostsGrouped),
	}

}
