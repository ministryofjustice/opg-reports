package govukassets_test

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ministryofjustice/opg-reports/internal/fileutils"
	"github.com/ministryofjustice/opg-reports/pkg/consts"
	"github.com/ministryofjustice/opg-reports/pkg/govukassets"
)

// TestGovUKAssetsFrontEndFull checks each method on the
// FrontEnder struct to make sure if downloads, extracts
// and moves files correctly
func TestGovUKAssetsFrontEndFull(t *testing.T) {
	var err error
	var tmpDir = t.TempDir() //"./tmp/"
	fe := govukassets.FrontEnd()

	// check the url
	url := fe.Url()
	if !strings.Contains(url, consts.GovUKFrontendVersion) {
		t.Errorf("version number was not set")
	}

	// now download the release
	release, err := fe.Download(time.Second)
	if err != nil {
		t.Errorf("error with download: [%s]", err.Error())
	}
	// check release zip exists
	if !fileutils.Exists(release) {
		t.Errorf("failed to download - local zip does not exist: [%s]", release)
	}

	// now extract the release
	dir, extracted, err := fe.Extract(release)
	if err != nil {
		t.Errorf("error with extraction: [%s]", err.Error())
	}
	// check dir exists
	if !fileutils.Exists(dir) {
		t.Errorf("failed to extract - local dir does not exist: [%s]", dir)
	}
	// check each asset exists
	for _, f := range extracted {
		path := filepath.Join(dir, f)
		if !fileutils.Exists(path) {
			t.Errorf("failed to extract - local asset does not exist: [%s]", path)
		}
	}
	// check the move to tmpDir works
	moved, err := fe.Move(dir, extracted, tmpDir)
	if err != nil {
		t.Errorf("error moving files: [%s]", err.Error())
	}

	if len(moved) != len(extracted) {
		t.Errorf("moved files dont match extracted ones")
	}

	for _, m := range moved {
		if !strings.Contains(m, "assets") {
			t.Errorf("error - assets should all be under folder: [%s]", m)
		}
	}

	// clean self up
	fe.Close()

	if fileutils.Exists(fe.DownloadLocation) {
		t.Errorf("clean up failed")
	}
	if fileutils.Exists(fe.ExtractedLocation) {
		t.Errorf("clean up failed")
	}

}

// TestGovUKAssetsFrontEndDo uses the short form helper call
// of .Do to run all the commands in sequence
func TestGovUKAssetsFrontEndDo(t *testing.T) {
	var err error
	var tmpDir = t.TempDir()
	var fe = govukassets.FrontEnd()
	var resources []string

	resources, err = fe.Do(tmpDir)

	if err != nil {
		t.Errorf("error within do: [%s]", err.Error())
	}

	for _, m := range resources {
		if !strings.Contains(m, "assets") {
			t.Errorf("error - assets should all be under folder: [%s]", m)
		}
	}

	fe.Close()

}

// testDo is used to test the defer close setup
func testDo(tmpDir string) (fe *govukassets.FrontEndAssets) {
	fe = govukassets.FrontEnd()
	defer fe.Close()
	fe.Do(tmpDir)
	return
}

// TestGovUKAssetsFrontEndDoDefer checks the defer fe.Close()
// works
func TestGovUKAssetsFrontEndDoDefer(t *testing.T) {
	var tmpDir = t.TempDir()

	fe := testDo(tmpDir)

	// so the file locations should all be removed
	if fileutils.Exists(fe.DownloadLocation) {
		t.Errorf("defer did not remove zip")
	}
	if fileutils.Exists(fe.ExtractedLocation) {
		t.Errorf("defer did not remove extraction")
	}
	if len(fe.Extracted) > 0 {
		t.Errorf("defer did not remove extraction data")
	}
}
