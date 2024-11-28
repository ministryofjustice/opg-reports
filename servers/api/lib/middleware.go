package lib

import (
	"github.com/danielgtaylor/huma/v2"
)

const CTX_DB_KEY string = "db-path"

// AddMiddleware adds the standard middleware information and process
// for each api segment
// Currently - adds database path as a value to the context
func AddMiddleware(api huma.API, segments map[string]*ApiSegment, dbPath string) {
	//
	api.UseMiddleware(func(ctx huma.Context, next func(huma.Context)) {
		for segment, cfg := range segments {
			ctx = huma.WithValue(ctx, segment, cfg.DbFile)
		}
		next(ctx)
		// new version of tracking db location
		var key = CTX_DB_KEY
		ctx = huma.WithValue(ctx, key, dbPath)
		next(ctx)
	})

}