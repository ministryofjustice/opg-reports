// Package govuk is used to download govuk assets from a known release version to a local directory.
//
// Used by front end of the reporting site to use built versions of css while avoiding introducing
// a significant js ecosystem.
package govuk

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/utils/files"
	"opg-reports/report/internal/utils/zips"
	"os"
	"path/filepath"

	"github.com/google/go-github/v81/github"
)

var (
	ErrFailedtoFindRelease       = errors.New("failed to find release matching requested options.")
	ErrNoAssetsInRelease         = errors.New("no assets attached to release.")
	ErrNoMatchingAssetsInRelease = errors.New("no matching assets attached to release.")
	ErrGithubAssetDownloadFailed = errors.New("failed to download github asset with error.")
	ErrZipExtractFailed          = errors.New("failed to extract zip in target directory.")
	ErrZipExtractMissingVersion  = errors.New("extracted zip did not contain version file.")
)

// GitHubClient wrapper around methods needed to fetch info from github to download release assets
//
// Wrapper for: *github.RepositoriesService
type GitHubClient interface {
	GetReleaseByTag(ctx context.Context, owner, repo, tag string) (*github.RepositoryRelease, *github.Response, error)
	DownloadReleaseAsset(ctx context.Context, owner, repo string, id int64, followRedirectsClient *http.Client) (rc io.ReadCloser, redirectURL string, err error)
}

type Options struct {
	Tag       string // Release tag to use for the download - v5.11.0
	Directory string // directory to extract the zip data into
}

// fixed data about data sources from github
const (
	owner      string = "alphagov"       // owner
	repository string = "govuk-frontend" // repo name
)

// DownloadFrontEnd[T GitHubClient] connects to github and checks the `alphagov/govuk-frontend` repository for the release tag specified
// and attempts to download the `release-${tag}.zip` asset attached to that release. It then extracts that zip file into the `Directory`
// that has been passed.
func DownloadFrontEnd[T GitHubClient](ctx context.Context, log *slog.Logger, gh T, options *Options) (extractedDir string, err error) {
	var (
		zip     string // zip file thats been downloaded
		release *github.RepositoryRelease
		asset   *github.ReleaseAsset
		lg      *slog.Logger = log.With("func", "domain.govuk.govuk.DownloadFrontEnd")
	)

	lg.With("options", options).Debug("starting ...")
	// find the release data
	lg.Debug("getting release ...")
	release, err = getRelease(ctx, log, gh, options)
	if err != nil {
		return
	}
	// find the metadata asset
	lg.Debug("getting release asset ...")
	asset, err = getReleaseAsset(ctx, log, release, options.Tag)
	if err != nil {
		return
	}
	// download the asset to a tmp path
	lg.Debug("downloading asset ...")
	zip, err = downloadAsset(ctx, log, gh, asset)
	if err != nil {
		return
	}
	// remove the tmp zip
	defer func() {
		os.Remove(zip)
	}()
	// extract the zip file
	lg.With("zip", zip).Debug("extracting zip file ...")
	extractedDir, err = extract(ctx, log, zip, options)
	if err != nil {
		return
	}
	lg.With("extractedDir", extractedDir).Debug("complete.")
	return
}

// extract the zip into the directory directly (no sub-folder) and confirm it worked
func extract(ctx context.Context, log *slog.Logger, zip string, options *Options) (extracted string, err error) {
	var lg *slog.Logger = log.With("func", "domain.govuk.govuk.extractAndGetTeams")
	var versionFile string = filepath.Join(options.Directory, "VERSION.txt")
	lg.Debug("starting ...")
	// make sure the dir exists
	os.MkdirAll(options.Directory, os.ModePerm)
	lg.With("extractTo", options.Directory).Debug("extracting zip ...")
	_, err = zips.Extract(zip, options.Directory)
	if err != nil {
		return
	}
	// check it extracted correctly by making sure directory exists and
	// version file is present
	if !files.DirExists(options.Directory) {
		err = ErrZipExtractFailed
		lg.Error("failed to extract directory", "dest", options.Directory)
		return
	}
	if !files.Exists(versionFile) {
		err = ErrZipExtractMissingVersion
		lg.Error("did not find versions.txt", "versionFile", versionFile)
		return
	}
	extracted = options.Directory
	lg.Debug("complete.")
	return
}

// downloadAsset downloads zip locally
func downloadAsset[T GitHubClient](ctx context.Context, log *slog.Logger, gh T, asset *github.ReleaseAsset) (src string, err error) {
	var (
		buff     io.ReadCloser
		dir      string
		fileDest string
		lg       *slog.Logger = log.With("func", "domain.govuk.govuk.downloadAsset")
	)
	lg.Debug("starting ...")

	dir, err = os.MkdirTemp("", "govuk-front-*")
	if err != nil {
		return
	}

	fileDest = filepath.Join(dir, *asset.Name)
	// try to download to buffer
	lg.Debug("downloading release asset ...")
	buff, _, err = gh.DownloadReleaseAsset(ctx, owner, repository, *asset.ID, http.DefaultClient)
	if err != nil {
		lg.Error("failed to download github release asset", "err", err.Error())
		err = errors.Join(ErrGithubAssetDownloadFailed, err)
		return
	}
	defer buff.Close()
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
func getReleaseAsset(ctx context.Context, log *slog.Logger, release *github.RepositoryRelease, tag string) (asset *github.ReleaseAsset, err error) {
	var lg *slog.Logger = log.With("func", "domain.govuk.govuk.getReleaseAsset")
	var assetName string = fmt.Sprintf("release-%s.zip", tag)

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
		lg.Error("no matching release asset found", "err", err.Error())
		return
	}
	lg.Debug("complete.")
	return
}

// getRelease finds the tagged release
func getRelease[T GitHubClient](ctx context.Context, log *slog.Logger, gh T, options *Options) (release *github.RepositoryRelease, err error) {
	var lg *slog.Logger = log.With("func", "domain.govuk.govuk.getRelease")
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
