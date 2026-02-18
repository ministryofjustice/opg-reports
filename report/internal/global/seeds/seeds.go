package seeds

import (
	"context"
	"opg-reports/report/internal/global/migrations"
	"opg-reports/report/internal/team/teamimport"
	"opg-reports/report/package/dbx"
)

var teamList []string = []string{
	"team-a",
	"team-b",
	"team-c",
	"team-d",
	"team-e",
	"team-f",
}

// environments are shared account environment names to use in seeded data
var environmentList []string = []string{
	"development",
	"pre-production",
	"integrations",
	"production",
}

// fix regions to used for seeds
var regionList []string = []string{
	"eu-west-1",
	"eu-west-2",
	"us-east-1",
	"NoRegion",
}

// aws service that we use for seeding
var serviceList []string = []string{
	"Amazon Relational Database Service",
	"Amazon Simple Storage Service",
	"AmazonCloudWatch",
	"Amazon Elastic Load Balancing",
	"AWS Shield",
	"AWS Config",
	"AWS CloudTrail",
	"AWS Key Management Service",
	"Amazon Virtual Private Cloud",
	"Amazon Elastic Container Service",
	"EC2 - Other",
}

type Results struct {
	Teams []*teamimport.Model
}

type Args struct {
	DB            string `json:"db"`             // --db
	Driver        string `json:"driver"`         // --driver
	Params        string `json:"params"`         // --params
	MigrationFile string `json:"migration_file"` // --file
}

func SeedAll(ctx context.Context, in *Args) (results *Results, err error) {
	var args = &dbx.InsertArgs{
		DB:     in.DB,
		Driver: in.Driver,
		Params: in.Params,
	}
	results = &Results{}
	// run the migrations
	err = migrations.MigrateAll(ctx, &migrations.Args{
		DB:            in.DB,
		Driver:        in.Driver,
		Params:        in.Params,
		MigrationFile: in.MigrationFile,
	})
	if err != nil {
		return
	}
	// seed teams
	results.Teams, err = seedTeams(ctx, args)
	if err != nil {
		return
	}

	return
}

func seedTeams(ctx context.Context, in *dbx.InsertArgs) (insert []*teamimport.Model, err error) {

	insert = []*teamimport.Model{}
	for _, team := range teamList {
		insert = append(insert, &teamimport.Model{Name: team})
	}
	err = dbx.Insert(ctx, teamimport.InsertStatement, insert, in)

	return
}
