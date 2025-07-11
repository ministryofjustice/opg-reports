package front

import (
	"fmt"
	"opg-reports/report/internal/repository/githubr"
	"opg-reports/report/internal/utils"
	"os"
	"strings"

	"github.com/google/go-github/v62/github"
)

// DownloadGovUKFrontEnd fetches the govuk front end generated release files
// directly from github based on the configuration values (Config.GovUK) to
// the local file system and extracts the zip into the directory stated
func (self *Service[T, R]) DownloadGovUKFrontEnd(
	client githubr.ClientRepositoryReleases,
	store githubr.RepositoryReleasesGetOneDownloader,
	directory string,
) (files []string, path string, err error) {
	var (
		file       string
		release    *github.RepositoryRelease
		asset      *github.ReleaseAsset
		owner      = self.conf.GovUK.Front.Owner      //"alphagov"
		repository = self.conf.GovUK.Front.Repository //"govuk-frontend"
		log        = self.log.With("operation", "DownloadGovUKFrontEnd", "repository", repository)
		relOpts    = &githubr.GetRepositoryReleaseOptions{
			ExcludePrereleases: true,
			ExcludeDraft:       true,
			ExcludeNoAssets:    true,
			ReleaseTag:         self.conf.GovUK.Front.ReleaseTag,
			UseRegex:           self.conf.GovUK.Front.UseRegex,
		}
		asOpts = &githubr.DownloadRepositoryReleaseAssetOptions{
			AssetName: self.conf.GovUK.Front.AssetName,
			UseRegex:  self.conf.GovUK.Front.UseRegex,
		}
	)
	// get the release
	release, err = store.GetRepositoryRelease(client, owner, repository, relOpts)
	if err != nil {
		log.Error("failed to find a matching repository release")
		return
	}
	log.With("release", *release.Name, "tag", *release.TagName).Debug("found release")
	// download the release
	asset, file, err = store.DownloadRepositoryReleaseAsset(client, owner, repository, release, directory, asOpts)
	if err != nil {
		log.Error("failed to download release asset")
		return
	}
	log.With("asset", *asset.Name, "content-type", *asset.ContentType).Debug("downloaded asset")
	// remove the source file
	defer os.Remove(file)

	// handle it being a zip
	if strings.HasSuffix(*asset.Name, "zip") || *asset.ContentType == "application/zip" {
		files, err = utils.ZipExtract(file, directory)
	} else {
		err = fmt.Errorf("unsupported file type: [name: %s type:%s]", *asset.Name, *asset.ContentType)
	}

	path = directory

	return
}
