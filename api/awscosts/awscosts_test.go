package awscosts_test

// import (
// 	"testing"

// 	"github.com/danielgtaylor/huma/v2"
// 	"github.com/danielgtaylor/huma/v2/humatest"
// 	"github.com/ministryofjustice/opg-reports/api/awscosts"
// )

// var testDBFile = "./test.db"

// // testApiAwsCostsMiddleware pushes in the location of the dummy database
// func testApiAwsCostsMiddleware(ctx huma.Context, next func(huma.Context)) {
// 	ctx = huma.WithValue(ctx, awscosts.Segment, testDBFile)
// 	next(ctx)
// }

// func TestApiAwsCostsTotal(t *testing.T) {

// 	_, api := humatest.New(t)
// 	api.UseMiddleware(testApiAwsCostsMiddleware)

// }
