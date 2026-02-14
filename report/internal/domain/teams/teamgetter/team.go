// Package teamgetter handles import teams from the opg-metadata repository which can then be used to populate the database
package teamgetter

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/domain/teams/teammodels"
	"opg-reports/report/internal/utils/files"
	"opg-reports/report/internal/utils/unmarshal"
	"opg-reports/report/internal/utils/zips"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v81/github"
)

var (
	ErrFailedtoFindRelease       = errors.New("failed to find release matching requested options.")
	ErrNoAssetsInRelease         = errors.New("no assets attached to release.")
	ErrNoMatchingAssetsInRelease = errors.New("no matching assets attached to release.")
	ErrGithubAssetDownloadFailed = errors.New("failed to download github asset with error.")
	ErrNoTeamsDatafile           = errors.New("no team.json data file found.")
)

// GitHubClient wrapper around methods needed to fetch info from github to download the metadata release
//
// Wrapper for: *github.RepositoriesService
type GitHubClient interface {
	GetReleaseByTag(ctx context.Context, owner, repo, tag string) (*github.RepositoryRelease, *github.Response, error)
	DownloadReleaseAsset(ctx context.Context, owner, repo string, id int64, followRedirectsClient *http.Client) (rc io.ReadCloser, redirectURL string, err error)
}

type Options struct {
	Tag           string
	DataDirectory string // tmp directory to write data into - will get deleted at the end
}

// fixed data about team sources - pulled from github
const (
	owner      string = "ministryofjustice" // org
	repository string = "opg-metadata"      // repo name
	assetName  string = "metadata.zip"      // the attached asset name in the release
	teamFile   string = "teams.json"        // the file to use in the release
)

// GetTeamData[T GitHubClient] connects to github and checks the releases for `opg-metadata` for the tag set within the options and then fetches
// the metadata.zip from that release, extracting and parsing the teams.json file, which is then returned as a slice.
//
// Note: opg-metadata is private, so suitable permissions are required on the github client (and its token).
func GetTeamData[T GitHubClient](ctx context.Context, log *slog.Logger, gh T, options *Options) (teams []*teammodels.Team, err error) {
	var (
		lg           *slog.Logger = log.With("func", "team.GetTeamData")
		release      *github.RepositoryRelease
		asset        *github.ReleaseAsset
		metadataFile string
	)
	// clean out the directory content
	defer func() {
		os.RemoveAll(options.DataDirectory)
	}()

	lg.With("options", options).Debug("starting ...")
	// find the release data
	lg.Debug("getting release ...")
	release, err = getRelease(ctx, log, gh, options)
	if err != nil {
		return
	}
	// find the metadata asset
	lg.Debug("getting release asset ...")
	asset, err = getReleaseAsset(ctx, log, release)
	if err != nil {
		return
	}
	// download the asset to a extract and parse the teams.json file
	lg.Debug("downloading asset ...")
	metadataFile, err = downloadAsset(ctx, log, gh, asset, options)
	if err != nil {
		return
	}
	// extract the zip and get the team data
	lg.Debug("extracting and converting ...")
	teams, err = extractAndGetTeams(ctx, log, metadataFile, options)
	if err != nil {
		return
	}

	lg.With("count", len(teams)).Debug("complete.")
	return
}

// extractAndGetTeams extract the file and parse the file
func extractAndGetTeams(ctx context.Context, log *slog.Logger, metadataZip string, options *Options) (teams []*teammodels.Team, err error) {
	var (
		lg        *slog.Logger = log.With("func", "team.extractAndGetTeams")
		extractTo string       = filepath.Join(options.DataDirectory, "ex")
		teamsFile string       = filepath.Join(extractTo, teamFile)
	)
	teams = []*teammodels.Team{}
	lg.Debug("starting ...")
	lg.Debug("extracting zip ...")
	_, err = zips.Extract(metadataZip, extractTo)
	if err != nil {
		return
	}
	if !files.Exists(teamsFile) {
		err = ErrNoTeamsDatafile
		return
	}
	lg.Debug("umarshal from file ...")
	err = unmarshal.FromFile(teamsFile, &teams)
	if err != nil {
		return
	}
	lg.Debug("lowercase names ...")
	for _, team := range teams {
		team.Name = strings.ToLower(team.Name)
	}
	lg.Debug("complete.")
	return
}

// downloadAsset downloads zip locally
func downloadAsset[T GitHubClient](ctx context.Context, log *slog.Logger, gh T, asset *github.ReleaseAsset, options *Options) (src string, err error) {
	var (
		buff     io.ReadCloser
		lg       *slog.Logger = log.With("func", "team.downloadAsset")
		dlTo     string       = filepath.Join(options.DataDirectory, "dl")
		fileDest string       = filepath.Join(dlTo, *asset.Name)
	)
	lg.Debug("starting ...")
	// try to download to buffer
	lg.Debug("downloading release asset ...")
	buff, _, err = gh.DownloadReleaseAsset(ctx, owner, repository, *asset.ID, http.DefaultClient)
	if err != nil {
		lg.Error("failed to download github release asset", "err", err.Error())
		err = errors.Join(ErrGithubAssetDownloadFailed, err)
		return
	}
	defer buff.Close()
	// write to local folder and file
	os.MkdirAll(dlTo, os.ModePerm)
	// copy to local dir
	lg.With("dest", fileDest).Debug("copying to local file ...")
	err = files.Copy(buff, fileDest)
	if err != nil {
		lg.Error("error downloading metadata release asset", "err", err.Error())
		return
	}
	src = fileDest
	lg.Debug("compelte.")
	return
}

// getReleaseAsset finds the named asset from the attached assets
func getReleaseAsset(ctx context.Context, log *slog.Logger, release *github.RepositoryRelease) (asset *github.ReleaseAsset, err error) {
	var lg *slog.Logger = log.With("func", "team.getReleaseAsset")
	lg.Debug("starting ...")
	// no assets, so a failure
	if len(release.Assets) <= 0 {
		err = ErrNoAssetsInRelease
		lg.Error("no assets found", "err", err.Error())
		return
	}
	// look for the asset name
	for _, item := range release.Assets {
		if *item.Name == assetName {
			asset = item
			break
		}
	}
	if asset == nil {
		err = ErrNoMatchingAssetsInRelease
		lg.Error("no metadata.zip asset found", "err", err.Error())
		return
	}
	lg.Debug("complete.")
	return
}

// getRelease finds the tagged release
func getRelease[T GitHubClient](ctx context.Context, log *slog.Logger, gh T, options *Options) (release *github.RepositoryRelease, err error) {
	var lg *slog.Logger = log.With("func", "team.getRelease")
	lg.Debug("starting ...")
	release, _, err = gh.GetReleaseByTag(ctx, owner, repository, options.Tag)
	if err != nil {
		lg.Error("error finding metadata release.", "err", err.Error())
		err = errors.Join(ErrFailedtoFindRelease, err)
		return
	}
	lg.Debug("complete.")
	return
}
