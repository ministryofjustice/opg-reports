package main

import (
	"context"
	"opg-reports/report/internal/cost/costimport"
	"opg-reports/report/package/awsclients"
	"opg-reports/report/package/awsid"
	"opg-reports/report/package/times"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
)

func importCosts(ctx context.Context, flags *cli) (err error) {
	// aws client
	client, err := awsclients.New[*costexplorer.Client](ctx, flags.Region)
	if err != nil {
		return
	}
	// run import
	err = costimport.Import(ctx, client, &costimport.Input{
		DB:        flags.DB,
		Driver:    flags.Driver,
		Params:    flags.Params,
		DateStart: times.MustFromString(flags.DateStart),
		DateEnd:   times.MustFromString(flags.DateEnd),
		AccountID: awsid.AccountID(ctx, flags.Region),
	})
	return

}
