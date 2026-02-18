package imports

import (
	"context"
	"opg-reports/report/internal/cost/costimport"
	"opg-reports/report/internal/global"
	"opg-reports/report/internal/global/migrations"
	"opg-reports/report/package/awsclients"
	"opg-reports/report/package/awsid"
	"opg-reports/report/package/times"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
)

type Args struct {
	// Database related
	DB            string `json:"db"`             // --db
	Driver        string `json:"driver"`         // --driver
	Params        string `json:"params"`         // --params
	MigrationFile string `json:"migration_file"` // --file
	Region        string `json:"region"`         // --region; aws region ; AWS
	DateStart     string `json:"date_start"`     // --start
	DateEnd       string `json:"date_end"`       // --end
	SrcFile       string `json:"src-file"`       // --src-file ; used for file based imports

}

func ImportCosts(ctx context.Context, flags *Args) (err error) {
	// run the migrations
	err = global.MigrateAll(ctx, &migrations.Args{
		DB:            flags.DB,
		Driver:        flags.Driver,
		Params:        flags.Params,
		MigrationFile: flags.MigrationFile,
	})
	if err != nil {
		return
	}
	// aws client
	client, err := awsclients.New[*costexplorer.Client](ctx, flags.Region)
	if err != nil {
		return
	}
	// run import
	err = costimport.Import(ctx, client, &costimport.Args{
		DB:        flags.DB,
		Driver:    flags.Driver,
		Params:    flags.Params,
		DateStart: times.MustFromString(flags.DateStart),
		DateEnd:   times.MustFromString(flags.DateEnd),
		AccountID: awsid.AccountID(ctx, flags.Region),
	})
	return

}
