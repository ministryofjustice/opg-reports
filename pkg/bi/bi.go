// Package bi provides build details that are replaced using ldflags
// during build that change values within the applications
package bi

import (
	"fmt"
	"log/slog"
)

var Semver string = "v0.0.1"                       // Semver tag used in build
var Commit string = "0"                            // Git commit hash used in the build
var Timestamp string = "2024-10-01T01:02:03Z00:00" // Time of the build
var ApiVersion string = "v1"                       // The api prefix to use for outbound calls from sfront
var Organisation string = "OPG"                    // The organisation name to use in the front end
var Mode string = "simple"                         // The mode is used to select which api handlers and navigation is used

func Dump() {
	slog.Info("Build info",
		slog.String("ApiVersion", ApiVersion),
		slog.String("Commit", Commit),
		slog.String("Mode", Mode),
		slog.String("Organisation", Organisation),
		slog.String("Semver", Semver),
		slog.String("Timestamp", Timestamp),
	)
}

// Signature generates the a string from build details.
// Currently:
//
//	`<semver> [<timestamp>] (<git-sha>)`
func Signature() string {
	return fmt.Sprintf("%s [%s] (%s)", Semver, Timestamp, Commit)
}
