package lib

import (
	"context"
	"log/slog"

	"github.com/danielgtaylor/huma/v2"
)

// SetupSegments runs the setup functions for all the segments
func SetupSegments(ctx context.Context, segments map[string]*ApiSegment) {
	for name, segment := range segments {
		slog.Info("[api] setup segments", slog.String("segment", name))
		segment.SetupFunc(ctx, segment.DbFile, true)
	}
}

// RegisterSegments calls the registration function for each of the api segments
// allowing them to attach their own routes to the api
func RegisterSegments(api huma.API, segments map[string]*ApiSegment) {
	slog.Info("[api] registering segments ...")
	for name, segment := range segments {
		slog.Info("[api] register segment", slog.String("segment", name))
		segment.RegisterFunc(api)
	}
	slog.Info("[api] segments registered.")
}
