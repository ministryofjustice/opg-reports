package lib

import (
	"log/slog"

	"github.com/danielgtaylor/huma/v2"
)

// RegisterHandlers
func RegisterHandlers(api huma.API, handlers map[string]RegisterHandlerFunc) {
	slog.Info("[api] registering handlers ...")
	for name, regFunc := range handlers {
		slog.Info("[api] register handler", slog.String("handler", name))
		regFunc(api)
	}
	slog.Info("[api] handlers registered.")
}
