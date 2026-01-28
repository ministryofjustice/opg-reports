// Package account handles import accounts from the opg-metadata repository which can then be used to populate the database
package account

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/domain/accounts/accountmodels"
	"opg-reports/report/internal/utils/files"
	"opg-reports/report/internal/utils/unmarshal"
	"opg-reports/report/internal/utils/zips"
	"os"
	"path/filepath"

	"github.com/google/go-github/v81/github"
)

// errors
var (
	ErrFailedtoFindRelease       = errors.New("failed to find release matching requested options.")
	ErrNoAssetsInRelease         = errors.New("no assets attached to release.")
	ErrNoMatchingAssetsInRelease = errors.New("no matching assets attached to release.")
	ErrGithubAssetDownloadFailed = errors.New("failed to download github asset with error.")
	ErrNoTeamsDatafile           = errors.New("no accounts.aws.json data file found.")
	ErrFailedtoUnmarshal         = errors.New("failed to unmarshal struct.")
)

// GitHubClient wrapper around methods needed to fetch info from github to download the metadata release
//
// Wrapper for: *github.RepositoriesService
type GitHubClient interface {
	GetReleaseByTag(ctx context.Context, owner, repo, tag string) (*github.RepositoryRelease, *github.Response, error)
	DownloadReleaseAsset(ctx context.Context, owner, repo string, id int64, followRedirectsClient *http.Client) (rc io.ReadCloser, redirectURL string, err error)
}

// GetAwsAccountDataOptions input struct for the command to specify some variable data
type GetAwsAccountDataOptions struct {
	Tag           string // the version tag on the release to utilise
	DataDirectory string // directory to write the data into, should be empty as it will be removed
}

// fixed data about team sources - pulled from github
const (
	owner          string = "ministryofjustice" // org
	repository     string = "opg-metadata"      // repo name
	assetName      string = "metadata.zip"      // the attached asset name in the release
	extractFile    string = "accounts.aws.json" // the file to use in the release
	downloadSubDir string = "dl"                // subdirectory to download data into
	extractSubDir  string = "ex"                // directory to extraxt zip to
)

// GetAccountData[T GitHubClient] connects to github and checks the releases for `opg-metadata` for the tag set within the options and then fetches
// the metadata.zip from that release, extracting and parsing the accounts.aws..json file, which is then returned as a slice.
//
// Return data as model used for importing to the database. This get call is used by import commands directly
//
// Note: opg-metadata is private, so suitable permissions are required on the github client (and its token).
func GetAwsAccountData[T GitHubClient](ctx context.Context, log *slog.Logger, gh T, options *GetAwsAccountDataOptions) (accounts []*accountmodels.AwsAccount, err error) {

	var (
		release      *github.RepositoryRelease
		asset        *github.ReleaseAsset
		metadataFile string
	)
	// clean out the directory content
	defer func() {
		os.RemoveAll(options.DataDirectory)
	}()

	log = log.With("package", "accounts", "func", "GetTeamData").With("options", options)
	log.Debug("starting ...")

	// find the release data
	log.Debug("getting release ...")
	release, err = getRelease(ctx, log, gh, options)
	if err != nil {
		return
	}
	// find the metadata asset
	log.Debug("getting release asset ...")
	asset, err = getReleaseAsset(ctx, log, release)
	if err != nil {
		return
	}
	// download the asset to a extract and parse the teams.json file
	log.Debug("downloading asset ...")
	metadataFile, err = downloadAsset(ctx, log, gh, asset, options)
	if err != nil {
		return
	}
	// extract the zip and get the team data
	log.Debug("extracting and converting to models ...")
	accounts, err = extractAndGet(ctx, log, metadataFile, options)
	if err != nil {
		return
	}

	log.With("count", len(accounts)).Debug("complete.")
	return
}

// extractAndGet extract the file and parse the file
func extractAndGet(ctx context.Context, log *slog.Logger, metadataZip string, options *GetAwsAccountDataOptions) (accounts []*accountmodels.AwsAccount, err error) {
	var (
		extractTo string = filepath.Join(options.DataDirectory, extractSubDir)
		file      string = filepath.Join(extractTo, extractFile)
	)
	accounts = []*accountmodels.AwsAccount{}

	_, err = zips.Extract(metadataZip, extractTo)
	if err != nil {
		return
	}
	if !files.Exists(file) {
		err = ErrNoTeamsDatafile
		return
	}
	err = unmarshal.FromFile(file, &accounts)
	if err != nil {
		err = errors.Join(ErrFailedtoUnmarshal, err)
		return
	}
	return
}

// downloadAsset downloads zip locally
func downloadAsset(ctx context.Context, log *slog.Logger, gh GitHubClient, asset *github.ReleaseAsset, options *GetAwsAccountDataOptions) (src string, err error) {
	var (
		buff     io.ReadCloser
		dlTo     string = filepath.Join(options.DataDirectory, downloadSubDir)
		fileDest string = filepath.Join(dlTo, *asset.Name)
	)
	// try to download
	// download to buffer
	buff, _, err = gh.DownloadReleaseAsset(ctx, owner, repository, *asset.ID, http.DefaultClient)
	if err != nil {
		log.Error("failed to download github release asset", "err", err.Error())
		err = errors.Join(ErrGithubAssetDownloadFailed, err)
		return
	}

	defer buff.Close()
	// write to local folder and file
	os.MkdirAll(dlTo, os.ModePerm)
	// copy to local dir
	err = files.Copy(buff, fileDest)
	if err != nil {
		log.Error("error downloading metadata release asset", "err", err.Error())
		return
	}
	src = fileDest
	return
}

// getReleaseAsset finds the named asset from the attached assets
func getReleaseAsset(ctx context.Context, log *slog.Logger, release *github.RepositoryRelease) (asset *github.ReleaseAsset, err error) {
	// no assets, so a failure
	if len(release.Assets) <= 0 {
		err = ErrNoAssetsInRelease
		log.Error("no assets found", "err", err.Error())
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
		log.Error("no metadata.zip asset found", "err", err.Error())
		return
	}
	return
}

// getRelease finds the tagged release
//
// Does not check for pre-release / draft status, just the tag name
func getRelease(ctx context.Context, log *slog.Logger, gh GitHubClient /**github.RepositoriesService*/, options *GetAwsAccountDataOptions) (release *github.RepositoryRelease, err error) {
	release, _, err = gh.GetReleaseByTag(ctx, owner, repository, options.Tag)
	if err != nil {
		log.Error("error finding metadata release", "err", err.Error())
		err = errors.Join(ErrFailedtoFindRelease, err)
		return
	}
	return
}
