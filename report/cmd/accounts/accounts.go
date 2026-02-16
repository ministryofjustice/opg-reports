package main

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/domain/accounts/accountcli"
	"opg-reports/report/internal/domain/accounts/accountgetter"
	"opg-reports/report/internal/domain/accounts/accountmodels"
	"opg-reports/report/internal/utils/ghclients"
	"opg-reports/report/internal/utils/logger"
	"opg-reports/report/internal/utils/marshal"
	"os"
	"path/filepath"

	"github.com/google/go-github/v81/github"
	"github.com/spf13/cobra"
)

// config items
var (
	ctx context.Context // default context
	log *slog.Logger    // default logger
)

// the command flags used on the import cli tool
type cli struct {
	Tag  string // represents --tag
	File string // --file
}

var ErrGitHubTokenMissing = errors.New("missing github token value.")

// default command options
var flags *cli = &cli{
	Tag:  "v0.1.26",
	File: "aws.accounts.json",
}

var rootCmd *cobra.Command = &cobra.Command{
	Use:   "accounts",
	Short: `accounts fetches the account json data from opg-metadata and write to a file`,
	RunE:  runCmd,
}

// runCmd main runner
func runCmd(cmd *cobra.Command, args []string) (err error) {
	var (
		dir    string
		client *github.Client
		data   []*accountmodels.Account = []*accountmodels.Account{}
		token  string                   = os.Getenv("GITHUB_TOKEN")
	)
	if token == "" {
		err = ErrGitHubTokenMissing
		return
	}
	// github client
	client, err = ghclients.New(ctx, log, token)
	if err != nil {
		return
	}
	dir, _ = os.MkdirTemp("", "__import-accounts-*")
	// run the data getter command
	data, err = accountgetter.GetAwsAccountData(ctx, log, client.Repositories, &accountgetter.Options{
		Tag:           flags.Tag,
		DataDirectory: dir,
	})
	if err != nil {
		return
	}

	os.MkdirAll(filepath.Dir(flags.File), os.ModePerm)
	os.WriteFile(flags.File, marshal.MustMarshal(data), os.ModePerm)

	return

}

// setup default values for config and logging
func init() {
	ctx = context.Background()
	log = logger.New(os.Getenv("LOG_LEVEL"), os.Getenv("LOG_TYPE"))

	rootCmd.Flags().StringVar(&flags.Tag, "tag", flags.Tag, "Release tag")
	rootCmd.Flags().StringVar(&flags.File, "file", flags.File, "Write accouts to this file")
}

func main() {
	var err error

	rootCmd.AddCommand(
		accountcli.CMD(ctx, log),
	)

	err = rootCmd.ExecuteContext(ctx)
	if err != nil {
		log.Error("error running command", "err", err.Error())
		os.Exit(1)
	}
}
