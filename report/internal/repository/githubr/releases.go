package githubr

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/google/go-github/v62/github"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

// getAllReleases returns all releases for a repository without any filtering.
//
// The repositoryName should not include the organsiation name
func (self *Repository) getAllReleases(organisation string, repositoryName string) (releases []*github.RepositoryRelease, err error) {
	var (
		client *github.Client
		page   int                 = 1
		opts   *github.ListOptions = &github.ListOptions{PerPage: 200}
		log                        = self.log.With("repositoryName", repositoryName, "operation", "getAllReleases")
	)
	// get api client
	client, err = self.client()
	if err != nil {
		return
	}
	// loop around the pagination
	for page > 0 {
		var response *github.Response
		var list []*github.RepositoryRelease
		// set the page number
		opts.Page = page
		// get all releases for the repository
		log.With("page", page).Debug("getting next page of releases")
		list, response, err = client.Repositories.ListReleases(self.ctx, organisation, repositoryName, opts)
		if err != nil {
			return
		}
		// if there items in the list, them merge into all
		if len(list) > 0 {
			releases = append(releases, list...)
		}
		// move to next page
		page = response.NextPage
	}

	return
}

// GetReleases returns all releases for a repository with some basic filtering options available.
//
// If options is nil (or all values are false) then all releases are returned.
// The repositoryName should not include the organsiation name
func (self *Repository) GetReleases(organisation string, repositoryName string, options *GetReleaseOptions) (releases []*github.RepositoryRelease, err error) {
	// setup log to be for this operation
	var log = self.log.With("repositoryName", repositoryName, "operation", "GetReleases")
	releases = []*github.RepositoryRelease{}
	// first, get all releases
	all, err := self.getAllReleases(organisation, repositoryName)

	// if there are no filter, then return everything
	if options == nil || (!options.ExcludeDraft && !options.ExcludeNoAssets && !options.ExcludePrereleases) {
		log.Debug("no filtering set, returning all releases directly")
		releases = all
		return
	}
	// there are filters, so look what to return
	log.Debug("filtering required, checking release values")
	// find only the required releases
	for _, release := range all {
		var keep bool = true
		// add release info to the log output
		var lg = log.With("draft", *release.Draft, "prerelease", *release.Prerelease, "assets", len(release.Assets), "id", *release.ID)
		// if this is a draft, but we are excluding them, swap keep to false
		if *release.Draft == true && options.ExcludeDraft == true {
			keep = false
		}
		// if this is a prerelease and we are excluding them, keep to false
		if *release.Prerelease == true && options.ExcludePrereleases == true {
			keep = false
		}
		// if there are no assets and we are excluding those without, the swap
		if len(release.Assets) == 0 && options.ExcludeNoAssets == true {
			keep = false
		}
		if keep {
			releases = append(releases, release)
		}
		lg.With("keep", keep).Debug("release checked")
	}

	return
}

// GetLatestRelease returns the latest published release for a repository.
//
// If you are looking for a prerelease / draft then use GetReleases to return all and then limit to what you need.
// The repositoryName should not include the organsiation name
func (self *Repository) GetLatestRelease(organisation string, repositoryName string) (release *github.RepositoryRelease, err error) {
	var client *github.Client
	var log = self.log.With("repositoryName", repositoryName, "operation", "GetLatestRelease")
	// get api client
	client, err = self.client()
	if err != nil {
		return
	}

	// get just 1 release - this should be the latest
	log.Debug("getting last release")
	release, _, err = client.Repositories.GetLatestRelease(self.ctx, organisation, repositoryName)
	return
}

// GetLatestReleaseAsset gets the latest release and then checks for an asset with a matching name.
//
// If the regex is true a regex is used to match against the asset name
func (self *Repository) GetLatestReleaseAsset(organisation string, repositoryName string, assetName string, useRegex bool) (asset *github.ReleaseAsset, err error) {
	// setup
	var (
		release *github.RepositoryRelease
		re      *regexp.Regexp = regexp.MustCompile(assetName)
		log                    = self.log.With("repositoryName", repositoryName, "operation", "GetLatestReleaseAsset", "assetName", assetName)
	)
	// get the release
	release, err = self.GetLatestRelease(organisation, repositoryName)
	// if error or not found, return
	if err != nil || release == nil {
		return
	}

	// if there are no assets, return but without an error
	if len(release.Assets) == 0 {
		log.With("assets", len(release.Assets), "id", *release.ID).Warn("no assets found")
		return
	}
	// check each asset and return when asset match is found
	for _, a := range release.Assets {
		var nm = *a.Name
		// if regex is enabled, then run the re check against the name
		if useRegex && len(re.FindStringIndex(nm)) > 0 {
			asset = a
			return
		} else if nm == assetName {
			asset = a
			return
		}
	}
	log.Warn("no matching asset found")

	return
}

// DownloadReleaseAsset tries to download the file associated with the assetID from via the github api
// and copy the content to the destination path given.
//
// If the body from the api call is empty the an error is returned. The repositoryName should not
// include the organsiation name
//
// The returned `*os.File` is not closed and will need to be handled
func (self *Repository) DownloadReleaseAsset(organisation string, repositoryName string, assetID int64, destinationFilePath string) (destination *os.File, err error) {
	var (
		rc io.ReadCloser
		//redirect string
		client *github.Client
		log    = self.log.With("operation", "DownloadAsset", "assetID", assetID)
	)

	// get api client
	client, err = self.client()
	if err != nil {
		return
	}
	log.Debug("downloading asset")
	rc, _, err = client.Repositories.DownloadReleaseAsset(self.ctx, organisation, repositoryName, assetID, http.DefaultClient)
	if err != nil {
		return
	}
	// not an error, but warn
	if rc == nil {
		err = fmt.Errorf("asset download was empty")
		return
	}
	// close at the end
	defer rc.Close()
	// copy the file
	err = utils.FileCopy(rc, destinationFilePath)
	// if there was no error with the copy, return a pointer to the file
	// and clean up the download directory path
	if err == nil {
		destination, err = os.Open(destinationFilePath)
	}

	return
}

func (self *Repository) DownloadReleaseAssetByName(organisation string, repositoryName string, assetName string, regex bool, directory string) (downloadedTo string, err error) {
	var asset *github.ReleaseAsset

	downloadedTo = directory
	asset, err = self.GetLatestReleaseAsset(organisation, repositoryName, assetName, regex)
	if err != nil || asset == nil {
		self.log.Error("error getting latest release asset", "err", err.Error())
		return
	}

	downloadedTo = filepath.Join(downloadedTo, asset.GetName())
	f, err := self.DownloadReleaseAsset(organisation, repositoryName, *asset.ID, downloadedTo)
	if err != nil {
		self.log.Error("error downloading the release asset", "err", err.Error())
		return
	}
	f.Close()

	return
}
