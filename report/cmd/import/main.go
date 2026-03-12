package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

// main root command
var root *cobra.Command = &cobra.Command{
	Use:               "import",
	Short:             `import data`,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

func init() {
	// var now = time.Now().UTC()

}

func main() {
	var err error
	var log *slog.Logger
	var ctx = context.Background()

	root.AddCommand()

	err = root.ExecuteContext(ctx)
	if err != nil {
		log.Error("error with command", "err", err.Error())
		panic("error")
		os.Exit(1)
	}

}
