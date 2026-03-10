package main

import (
	codebases "opg-reports/report/internal/domains/code/codebases/importer"
	"opg-reports/report/packages/clients"
	"os"

	"github.com/google/go-github/v84/github"
	"github.com/spf13/cobra"
)

var importCodebases = &cobra.Command{
	Use:   `codebases`,
	Short: `import codebases from github.`,
	RunE:  importCodebasesF,
}

// importCodebasesF
func importCodebasesF(cmd *cobra.Command, args []string) (err error) {
	var client *github.Client
	// get the client
	client, err = clients.New[*github.Client](cmd.Context(), os.Getenv("GITHUB_TOKEN"))
	if err != nil {
		return
	}
	runner(
		cmd.Context(),
		client,
		cliFlags,
		codebases.InsertStatement,
		codebases.Filter,
		codebases.Transform,
		codebases.Get,
	)
	return
}
