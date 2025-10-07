package githubr

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/utils"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/google/go-github/v75/github"
	"github.com/maruel/natural"
)

// DownloadRepositoryReleaseAssetOptions is used for configured finding releases from the sdk
type DownloadRepositoryReleaseAssetOptions struct {
	UseRegex  bool   // UseRegex enables using `AssetName` as a regex pattern rather than an exact match
	AssetName string // Comapred against `asset.Name` attribute
}

// Validate checks if the asset passed meets the requirements for this current setup
func (self *DownloadRepositoryReleaseAssetOptions) Validate(asset *github.ReleaseAsset) (valid bool) {
	var assetNameRegex *regexp.Regexp

	valid = false
	if self.AssetName != "" {
		assetNameRegex = regexp.MustCompile(self.AssetName)
		if self.UseRegex && len(assetNameRegex.FindStringIndex(*asset.Name)) > 0 {
			valid = true
		} else if !self.UseRegex && self.AssetName != "" && self.AssetName == *asset.Name {
			valid = true
		}
	}
	return
}

// GetRepositoryReleaseOptions used to configure what releases are returned from the github api
type GetRepositoryReleaseOptions struct {
	ExcludePrereleases bool   // Exclude releases marked as prereleases
	ExcludeDraft       bool   // Exclude anything marked as a draft
	ExcludeNoAssets    bool   // Exclude anything that does not have assets
	UseRegex           bool   // UseRegex enables using `ReleaseName` and `ReleaseTag` as a regex pattern rather than an exact match
	ReleaseName        string // ReleaseName is compared to the `.Name` attribute of the release
	ReleaseTag         string // ReleaseTag is compared to `.TagName` attribute
}

// Validate checks if the release passed meets the requirements for this current setup
func (self *GetRepositoryReleaseOptions) Validate(release *github.RepositoryRelease) (valid bool) {
	var (
		releaseNameRegex *regexp.Regexp
		releaseTagRegex  *regexp.Regexp
	)
	valid = true
	// no self to check, so true
	if self == nil {
		return
	}
	// skip draft when set
	if self.ExcludeDraft && *release.Draft == true {
		valid = false
	}
	// skip prereleases
	if self.ExcludePrereleases && *release.Prerelease == true {
		valid = false
	}
	// skipping releses without assets
	if self.ExcludeNoAssets && len(release.Assets) <= 0 {
		valid = false
	}
	// if the simple comparisions have determined false, then return to
	// avoid some of the more complex checks
	if !valid {
		return
	}
	// name filtering by either regex if enabled or exact match if not
	if self.ReleaseName != "" {
		releaseNameRegex = regexp.MustCompile(self.ReleaseName)
		releaseNameMatch := false
		if self.UseRegex && len(releaseNameRegex.FindStringIndex(*release.Name)) > 0 {
			releaseNameMatch = true
		} else if !self.UseRegex && self.ReleaseName == *release.Name {
			releaseNameMatch = true
		}
		valid = (valid && releaseNameMatch)
	}

	// filter by the release tag, similar to the name filtering
	if self.ReleaseTag != "" {
		releaseTagRegex = regexp.MustCompile(self.ReleaseTag)
		releaseTagMatch := false

		if self.UseRegex && len(releaseTagRegex.FindStringIndex(*release.TagName)) > 0 {
			releaseTagMatch = true
		} else if !self.UseRegex && self.ReleaseTag == *release.TagName {
			releaseTagMatch = true
		}
		valid = (valid && releaseTagMatch)
	}

	return
}

// GetRepositoryReleases returns fetches all releases from the client (by calling `ListReleases`) and in turn filters
// each release against the criteria set in `options` paramter, so only the desired type of releases are returned.
//
//   - if `options` is nil, then no filtering is applied
//   - Matching releases are order by TagName descending
//   - ClientRepositoryReleaseListReleases is github.Client.Repositories.ListReleases
func (self *Repository) GetRepositoryReleases(
	client ClientRepositoryReleaseListReleases,
	owner string, repository string,
	options *GetRepositoryReleaseOptions) (releases []*github.RepositoryRelease, err error) {
	var (
		ctx                                         = context.Background()
		page   int                                  = 1
		opts   *github.ListOptions                  = &github.ListOptions{PerPage: 200}
		log    *slog.Logger                         = self.log.With("repository", repository, "operation", "GetRepositoryReleases")
		toSort []string                             = []string{}
		rels   map[string]*github.RepositoryRelease = map[string]*github.RepositoryRelease{}
	)
	// loop around with pagination to fetch all releases
	// 	- filter the releases based on options set
	for page > 0 {
		var response *github.Response
		var list []*github.RepositoryRelease
		// set the page number
		opts.Page = page
		// get all releases for the repository
		log.With("page", page).Debug("getting repository releases ... ")
		list, response, err = client.ListReleases(ctx, owner, repository, opts)
		if err != nil {
			log.Error("error getting repository releases", "err", err.Error())
			return
		}
		log.With("page", page, "count", len(list)).Debug("found releases ...")
		// if there items in the list, them merge into all
		if len(list) > 0 {
			for _, item := range list {
				var include = options.Validate(item)
				log.With("include", include, "release", item.GetTagName()).Debug("include release?")
				if include && item.GetTagName() != "" {
					rels[item.GetTagName()] = item
				}
			}
		}
		// move to next page
		page = response.NextPage
	}
	// -- sorting
	// grab all the tagnames
	for _, r := range rels {
		toSort = append(toSort, *r.TagName)
	}
	// sort tag names in highest string value first
	sort.Sort(sort.Reverse(natural.StringSlice(toSort)))
	// now create the final set of releases from the sorted list of keys
	for _, key := range toSort {
		releases = append(releases, rels[key])
	}
	return
}

// GetRepositoryRelease calls GetRepositoryReleases on itself and simply returns the first release found
//
//   - if `options` is nil, then no filtering is applied
//   - ClientRepositoryReleaseListReleases is github.Client.Repositories.ListReleases
func (self *Repository) GetRepositoryRelease(
	client ClientRepositoryReleaseListReleases,
	owner string, repository string,
	options *GetRepositoryReleaseOptions) (release *github.RepositoryRelease, err error) {

	releases, err := self.GetRepositoryReleases(client, owner, repository, options)
	if err != nil {
		return
	}
	if len(releases) > 0 {
		release = releases[0]
	}
	return
}

// DownloadRepositoryReleaseAsset tries to download an asset that matches the oprions specified
// thats attached to the release passed along. File will be downloaded to destinationDirectory
// using the same name as is set on github.
//
//   - if `options` is nil, then no filtering is applied and first asset is used
//   - if no assets are found to match (or are present) then an error is returned
//   - ClientRepositoryReleaseDownloadReleaseAsset is github.Client.Repositories.DownloadReleaseAsset
func (self *Repository) DownloadRepositoryReleaseAsset(
	client ClientRepositoryReleaseDownloadReleaseAsset,
	owner string, repository string,
	release *github.RepositoryRelease,
	destinationDirectory string,
	options *DownloadRepositoryReleaseAssetOptions) (asset *github.ReleaseAsset, destination string, err error) {
	var (
		buff io.ReadCloser
		log  *slog.Logger = self.log.With("operation", "DownloadRepositoryReleaseAsset", "release", *release.ID)
	)
	// when find the first match, break the loop
	for _, a := range release.Assets {
		var found = options == nil || options.Validate(a)
		if found {
			asset = a
			break
		}
	}
	// return error for not matching asset
	if asset == nil {
		err = fmt.Errorf("no asset found matching criteria")
		log.Error("unable to find suitable release asset")
		return
	}

	// otherwise, lets download it
	buff, _, err = client.DownloadReleaseAsset(self.ctx, owner, repository, *asset.ID, http.DefaultClient)
	if err != nil {
		log.Error("failed to download release asset")
		return
	}
	if buff == nil {
		err = fmt.Errorf("asset download was empty")
		return
	}
	// close file
	defer buff.Close()
	// make the directory
	os.MkdirAll(destinationDirectory, os.ModePerm)
	// set the destination name
	destination = filepath.Join(destinationDirectory, *asset.Name)
	// make sure parent

	if err = utils.FileCopy(buff, destination); err != nil {
		log.Error("failed to copy asset to destination", "destination", destination, "err", err)
		return
	}

	return
}
