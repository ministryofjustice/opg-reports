package info

import (
	"fmt"
	"log/slog"
)

// Values that are replaced at build relating to versioning details
const (
	Commit    string = "eb5aa62ac8848d48c6ee37ab629f52a3c70e2f55" // Current commit
	Timestamp string = "2024-06-09T10:07:33Z00:00"                // Current timestamp
	Semver    string = "0.0.1"                                    // Current semver
)

// Values that are replaced at build time which relate to configuration
const (
	Organisation string = "OPG"  // Name of the organisation - used to display in the front end
	Dataset      string = "real" // Used by the api init to decide if we should download real or generate seeded database
	Fixtures     string = "full" // Used by front and api to determine if we're using all areas of just standards
)

// Log outputs the build and config details via slog.Info
func Log() {
	slog.Info("Build info",
		slog.String("Semver", Semver),
		slog.String("Commit", Commit),
		slog.String("Timestamp", Timestamp))
	slog.Info("Config info",
		slog.String("Organisation", Organisation),
		slog.String("Dataset", Dataset),
		slog.String("Fixtures", Fixtures))
}

// BuildInfo returns a formatted string containing all of the build constants
// and is normally used by api and front to display data about the build
func BuildInfo() string {
	return fmt.Sprintf("%s [%s] (%s)", Semver, Commit, Timestamp)
}

// ConfigInfo returns string containing all of the configuration values
func ConfigInfo() string {
	return fmt.Sprintf("%s [%s] (%s)", Organisation, Dataset, Fixtures)
}
