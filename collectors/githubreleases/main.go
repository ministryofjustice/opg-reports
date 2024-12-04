/*
githubreleases fetches release data (using mix of workflow runs and merges to main as a proxy).

Usage:

	githubreleases [flags]

The flags are:

	-organisation=<organisation>
		The name of the github organisation.
		Default: `ministryofjustice`
	-team=<github-team>
		Team slug for whose repos to check.
		Default: `opg`
	-start=<yyyy-mm-dd>
		Start date to fetch data for.
	-end=<yyyy-mm-dd>
		End date to fetch data for.
	-output=<path-pattern>
		Path (with magic values) to the output file
		Default: `./data/{day}_github_releases.json`

The command presumes an active, authorised session that can connect
to GitHub.
*/
package main

import (
	"fmt"
	"log/slog"

	"github.com/ministryofjustice/opg-reports/collectors/githubreleases/lib"
)

var (
	args = &lib.Arguments{}
)

func main() {
	var err error
	lib.SetupArgs(args)

	slog.Info("[githubreleases] starting ...")
	slog.Debug("[githubreleases]", slog.String("args", fmt.Sprintf("%+v", args)))

	err = lib.Run(args)
	if err != nil {
		panic(err)
	}

	slog.Info("[githubreleases] done.")

}
