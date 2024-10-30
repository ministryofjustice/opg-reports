// Package bi provides build details that are replaced using ldflags
// during build that change values within the applications
package bi

import "log/slog"

var Semver string = "v0.0.1"
var Commit string = "0"
var Timestamp string = "2024-10-01T01:02:03Z00:00"
var Organisation string = "OPG"
var ApiVersion string = "v1"

func Dump() {
	slog.Info("Build info",
		slog.String("ApiVersion", ApiVersion),
		slog.String("Commit", Commit),
		slog.String("Organisation", Organisation),
		slog.String("Semver", Semver),
		slog.String("Timestamp", Timestamp),
	)
}
