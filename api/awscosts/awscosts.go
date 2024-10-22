package awscosts

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/awscosts"
	"github.com/ministryofjustice/opg-reports/datastore"
	"github.com/ministryofjustice/opg-reports/fakes"
)

const segment string = "awscosts"         // the name of the api segment this package handles
const dbFile string = "./dbs/awscosts.db" // relative location of the database

var databaseConfig *datastore.Config = datastore.Sqlite

// Setup will download or generate a database
// for the api to operate from
func Setup(ctx context.Context) (err error) {
	var db *sqlx.DB
	var isNewDb bool = false

	db, isNewDb, err = datastore.New(ctx, databaseConfig, dbFile)
	if err != nil {
		return
	}
	awscosts.Create(ctx, db)
	// If this ia new database, then seed it with random data
	if isNewDb {
		err = seed(ctx, db)
	}
	return
}

func seed(ctx context.Context, db *sqlx.DB) (err error) {
	var (
		count int              = 100000
		org   string           = "seeded-org"
		seeds []*awscosts.Cost = []*awscosts.Cost{
			{Organisation: org, AccountID: "10000001", AccountName: "AA", Unit: "GroupA", Label: "A Prod", Environment: "production"},
			{Organisation: org, AccountID: "10000002", AccountName: "AB", Unit: "GroupA", Label: "A Dev", Environment: "development"},
			{Organisation: org, AccountID: "20000001", AccountName: "BA", Unit: "GroupB", Label: "B Prod", Environment: "production"},
			{Organisation: org, AccountID: "20000002", AccountName: "BB", Unit: "GroupB", Label: "B Dev", Environment: "development"},
			{Organisation: org, AccountID: "30000001", AccountName: "CA", Unit: "GroupC", Label: "C Prod", Environment: "production"},
			{Organisation: org, AccountID: "30000002", AccountName: "CB", Unit: "GroupC", Label: "C Dev", Environment: "development"},
			{Organisation: org, AccountID: "40000001", AccountName: "DA", Unit: "GroupD", Label: "D Prod", Environment: "production"},
			{Organisation: org, AccountID: "40000002", AccountName: "DB", Unit: "GroupD", Label: "D Dev", Environment: "development"},
			{Organisation: org, AccountID: "50000001", AccountName: "EA", Unit: "GroupE", Label: "D Prod", Environment: "production"},
			{Organisation: org, AccountID: "50000002", AccountName: "EB", Unit: "GroupE", Label: "D Dev", Environment: "development"},
		}
		insert []*awscosts.Cost = []*awscosts.Cost{}
	)

	for i := 0; i < count; i++ {
		base := fakes.Choice(seeds)
		insert = append(insert, awscosts.FakeFrom(base))
	}

	_, err = awscosts.InsertMany(ctx, db, insert)

	return
}
