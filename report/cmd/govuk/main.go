package main

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/conf"
	"opg-reports/report/internal/domain/govuk/govuk"
	"opg-reports/report/internal/utils/ghclients"
	"opg-reports/report/internal/utils/logger"
	"os"

	"github.com/google/go-github/v81/github"
	"github.com/spf13/cobra"
)

const (
	cmdName   string = "govuk" // root command name
	shortDesc string = `govuk downloads pre-built css and associated assets from a specificed release.`
	longDesc  string = `
govuk downloads pre-built css and associated assets from a specified release of the alphagov/govuk-frontend
repoistory and extracts the zip into the requested folder.
`
)

// config items
var (
	cfg *conf.Config    // default config
	ctx context.Context // default context
	log *slog.Logger    // default logger
)

var (
	rootCmd *cobra.Command = &cobra.Command{
		Use:   cmdName,
		Short: shortDesc,
		Long:  longDesc,
		RunE:  govUKRunE,
	}
)

var govUKReleaseTag = "v5.11.0"
var govUKDir string = "govuk"

var (
	ErrGitHubTokenMissing = errors.New("missing github token value.")
	ErrGitHubConnFailed   = errors.New("github client failed with error.")
)

// govUKRunE passed into the cobra command
func govUKRunE(cmd *cobra.Command, args []string) (err error) {
	var client *github.Client
	// get the github client
	client, err = ghclient()
	if err != nil {
		return
	}
	err = downloadAssets(ctx, log, client.Repositories, &govuk.Options{
		Tag:       govUKReleaseTag,
		Directory: govUKDir,
	})
	return
}

// downloadAssets run the download
func downloadAssets(ctx context.Context, log *slog.Logger, client govuk.GitHubClient, opts *govuk.Options) (err error) {
	var dir string
	var lg *slog.Logger = log.With("func", "govuk.downloadAssets")

	lg.Info("starting govuk download command ...")
	lg.With("opts", opts).Debug("govuk options ...")
	dir, err = govuk.DownloadFrontEnd(ctx, log, client, opts)
	if err != nil {
		return
	}
	lg.With("directory", dir).Info("complete.")
	return
}

func ghclient() (client *github.Client, err error) {
	// fail if there is no github token
	if cfg.GithubToken == "" {
		err = ErrGitHubTokenMissing
		return
	}
	// create client
	client, err = ghclients.New(ctx, log, cfg.GithubToken)
	if err != nil {
		log.Error("error connecting to client.", "err", err.Error())
		err = errors.Join(ErrGitHubConnFailed, err)
		return
	}
	return
}

func setup() {
	cfg = conf.New()
	ctx = context.Background()
	log = logger.New(cfg.Log.Level, cfg.Log.Type)
}

// setup default values for config and logging & add options
func init() {
	setup()
	rootCmd.Flags().StringVar(&govUKReleaseTag, "release-tag", govUKReleaseTag, "Release to fetch govuk assets from.")
	rootCmd.Flags().StringVar(&govUKDir, "directory", govUKDir, "Directory to extract zip into.")
}

func main() {
	var err error

	err = rootCmd.Execute()
	if err != nil {
		log.Error("error running command", "err", err.Error())
		os.Exit(1)
	}
}
