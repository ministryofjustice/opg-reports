// Package awscosts provides struct and database methods for handling cost explorer data
// that is then used by the api
package awscosts

// const Segment string = "awscosts"
// const Tag string = "AWS Costs"

// var API api = api{
// 	Register: register,
// }

// var SampleCost = &Cost{Organisation: "foobar"}

// func Setup(ctx context.Context, dbFilepath string) {
// 	var err error
// 	var db *sqlx.DB
// 	var isNew bool = false
// 	var n int = 15000

// 	db, isNew, err = datastore.NewDB(ctx, datastore.Sqlite, dbFilepath)
// 	defer db.Close()

// 	if err != nil {
// 		panic(err)
// 	}

// 	datastore.Create(ctx, db, DB.Create)
// 	if isNew {
// 		_, err = datastore.InsertMany(ctx, db, DB.Insert, Fakes(n, SampleCost))
// 	}
// 	if err != nil {
// 		panic(err)
// 	}
// }
