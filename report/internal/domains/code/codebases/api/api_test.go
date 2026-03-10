package api

// func TestInternalDomainsCodebasesApi(t *testing.T) {
// 	var (
// 		err    error
// 		ctx    = t.Context()
// 		dir    = t.TempDir()
// 		driver = "sqlite3"
// 		dbpath = filepath.Join(dir, "test-api.db")
// 		opts   = args.Default[*args.API](time.Now().UTC())
// 		db     = &args.DB{Driver: driver, DB: dbpath}
// 	)
// 	opts.DB = db
// 	// seed database with dummy data
// 	_, err = seeds.Seed(ctx, db)
// 	if err != nil {
// 		t.Errorf("unexpected error seeding data: %s", err.Error())
// 		t.FailNow()
// 	}

// 	// mock an api call
// 	url := "/test/"
// 	req := httptest.NewRequest(http.MethodGet, url, nil)
// 	writer := httptest.NewRecorder()
// 	mux := http.NewServeMux()
// 	// register
// 	handler.RegisterAPI(ctx, mux, Handler(opts), "/test/{$}")
// 	// call
// 	mux.ServeHTTP(writer, req)
// 	// result
// 	res := writer.Result()
// 	resp := &models.APIResponse[*types.Codebase, *types.Codebase, *models.Request, *models.Filter]{}
// 	err = convert.Response(res, &resp)
// 	if err != nil {
// 		t.Errorf("unexpected error: %s", err.Error())
// 	}

// 	// check count - against the seed number
// 	if len(resp.Data) != 50 {
// 		t.Errorf("incorrect number of records returned.")
// 	}

// }
