// Package govukassets provides utils for downloading and working with govuk front end
package govukassets

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/ministryofjustice/opg-reports/pkg/consts"
	"github.com/ministryofjustice/opg-reports/pkg/fileutils"
)

// url (with version value to replace) for the release zip
const frontEndLocation string = "https://github.com/alphagov/govuk-frontend/releases/download/v{version}/release-v{version}.zip"

// FrontEndConfig contains setup for downloading and extracting
// govuk frontend built assets to directory of choice
type FrontEndConfig struct {
	Version           string   // Version number to use
	Source            string   // Source is the formatted string to fetch the release zip from
	DownloadLocation  string   // DownloadLocation is where the zip was downloaded to
	ExtractedLocation string   // ExtractedLocation is the new directory the zip was extracted in to
	Extracted         []string // List of files that have been extracted
}

// Url generates the url from the version and source strings
func (self *FrontEndConfig) Url() (url string) {
	url = strings.ReplaceAll(self.Source, "{version}", self.Version)
	return
}

// Download gets the zip file and downloads it into a temporary
// directoy and returns the full path to the zip
// This path can then be used with Extract
func (self *FrontEndConfig) Download(timeout time.Duration) (path string, err error) {
	var (
		url     string = self.Url()
		zipName string = "assets.zip"
		dir     string
	)
	dir, _ = os.MkdirTemp(os.TempDir(), "download-*")

	path, err = fileutils.DownloadFromUrl(url, dir, zipName, timeout)
	self.DownloadLocation = path
	return
}

// Extract with extract the zip file in the path into a new temp directory
// so it can be moved around to
func (self *FrontEndConfig) Extract(zipFilepath string) (extractedDir string, extracted []string, err error) {
	extractedDir, _ = os.MkdirTemp(os.TempDir(), "extract-*")

	extracted, err = fileutils.ZipExtract(zipFilepath, extractedDir)

	self.ExtractedLocation = extractedDir
	self.Extracted = extracted

	return
}

// Move takes all the files from the originDir and moves them over to the destinationDir
// Returns an updated list of files with their new paths and updates .Extracted.
//
// Note: the current govuk frontend has css & js at a top level so this function
// also moves those to live under the /assets/ path so there is one single
// path for all assets - `/assets/` - making includes / redirects easier
func (self *FrontEndConfig) Move(originDir string, files []string, destinationDir string) (moved []string, err error) {
	var v = consts.GovUKFrontendVersion
	moved = []string{}

	// move the css and js files into the assets folder
	// so all govuk related items are in one structure
	var moveToAssetsSubFolder []string = []string{
		fmt.Sprintf("govuk-frontend-%s.min.css", v),
		fmt.Sprintf("govuk-frontend-%s.min.css.map", v),
		fmt.Sprintf("govuk-frontend-%s.min.js", v),
		fmt.Sprintf("govuk-frontend-%s.min.js.map", v),
		"VERSION.txt",
	}

	// copy the remaining assets
	for _, current := range files {
		var origin = filepath.Join(originDir, current)
		var dest = filepath.Join(destinationDir, current)
		var moveSub = slices.Contains(moveToAssetsSubFolder, current)
		// see if it should be moved under assets
		if moveSub {
			dest = filepath.Join(destinationDir, "assets", current)
		}

		// copy the file over to new place
		if err = fileutils.CopyFromPath(origin, dest); err != nil {
			return
		}
		moved = append(moved, dest)
	}

	self.Extracted = moved
	err = os.RemoveAll(originDir)

	return
}

// Do handles the full logic of fetching and extracting the resources
func (self *FrontEndConfig) Do(destinationDir string) (resources []string, err error) {
	var (
		releaseZip   string
		extractedDir string
		extracted    []string
		timeout      = time.Second * 2
	)
	resources = []string{}
	// download the file
	if releaseZip, err = self.Download(timeout); err != nil {
		return
	}
	// extract the zip and capture dir and files
	if extractedDir, extracted, err = self.Extract(releaseZip); err != nil {
		return
	}
	// move the files from the extracted location into the destination
	if resources, err = self.Move(extractedDir, extracted, destinationDir); err != nil {
		return
	}
	if len(resources) != len(extracted) {
		err = fmt.Errorf("error moving extracted resources - resulting count does not match")
	}

	return
}

// Close removes tmp directories in use and generally cleans up after itself
func (self *FrontEndConfig) Close() {
	if self.DownloadLocation != "" {
		os.RemoveAll(self.DownloadLocation)
	}
	if self.ExtractedLocation != "" {
		os.RemoveAll(self.ExtractedLocation)
	}
	self.Extracted = []string{}
}

func FrontEnd() *FrontEndConfig {
	return &FrontEndConfig{
		Version: consts.GovUKFrontendVersion,
		Source:  frontEndLocation,
	}
}
