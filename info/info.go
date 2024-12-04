package info

import (
	"fmt"
	"log/slog"
)

var (
	Commit    string = "eb5aa62ac8848d48c6ee37ab629f52a3c70e2f55" // Current commit
	Timestamp string = "2024-06-09T10:07:33Z00:00"                // Current timestamp
	Semver    string = "0.0.1"                                    // Current semver
)

var (
	Organisation string = "OPG"                     // Name of the organisation - used to display in the front end
	Dataset      string = "real"                    // Used by the api init to decide if we should download real or generate seeded database
	Fixtures     string = "full"                    // Used by front and api to determine if we're using all areas of just standards
	BucketName   string = "report-data-development" // Name of the bucket where all data is stored
)

// Server default locations
const (
	ServerDefaultFrontAddr string = "localhost:8080"
	ServerDefaultApiAddr   string = "localhost:8081"
)

// CostsBillingDay is the day in the month where updated billing data has been
// fetched and we can use the month prior to the current.
//
// Example - 10th April would be looking at Feb costs, 15th April can see March costs
const AwsBillingDay int = 15

// GovUK asset version to download
const GovUKFrontendVersion string = "5.7.1"

// Log outputs the build and config details via slog.Info
func Log() {
	slog.Info("Build info",
		slog.String("Semver", Semver),
		slog.String("Commit", Commit),
		slog.String("Timestamp", Timestamp),
	)
	slog.Info("Config info",
		slog.String("Organisation", Organisation),
		slog.String("Dataset", Dataset),
		slog.String("Fixtures", Fixtures),
		slog.String("BucketName", BucketName),
	)
}

// BuildInfo returns a formatted string containing all of the build constants
// and is normally used by api and front to display data about the build
func BuildInfo() string {
	return fmt.Sprintf("%s [%s] (%s)", Semver, Commit, Timestamp)
}

// ConfigInfo returns string containing all of the configuration values
func ConfigInfo() string {
	return fmt.Sprintf("%s [dataset:%s] [fixtures:%s] [bucketname:%s]", Organisation, Dataset, Fixtures, BucketName)
}
