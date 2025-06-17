package main

import (
	"log/slog"
	"os"

	"github.com/ministryofjustice/opg-reports/cmd/data/awscosts"
)

func cliRunner(logger *slog.Logger, dataSource string, args []string) (err error) {

	logger = logger.With("source", dataSource)

	switch dataSource {
	case "awscosts":
		awscosts.Run(logger, args)
	}

	return

}

// main
//
//	<data-source> --commands
func main() {
	var (
		logger            = slog.New(slog.NewTextHandler(os.Stdout, nil))
		dataSource string = ""
	)

	// check nubmer of args
	if len(os.Args) < 2 {
		logger.Error("No arguments passed")
		os.Exit(1)
	}

	dataSource = os.Args[1]

	cliRunner(logger, dataSource, os.Args[2:])
}
