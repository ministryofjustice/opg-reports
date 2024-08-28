package govassets

import (
	"archive/zip"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

const govukVersion string = "v5.4.0"

// DonwloadGovUKAssets gets a zip file from govuk-frontend of all assets
// downloads it locally and populates file system in expected format
func DownloadAssets() {
	slog.Warn("Downloading govuk assets", slog.String("v", govukVersion))

	zipUrl := fmt.Sprintf("https://github.com/alphagov/govuk-frontend/releases/download/%s/release-%s.zip",
		govukVersion, govukVersion)

	resp, err := http.Get(zipUrl)
	if err != nil {
		slog.Error("error getting assets", slog.String("err", err.Error()))
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		slog.Error("error status with assets", slog.Int("status", resp.StatusCode))
		return
	}
	os.MkdirAll("govuk", os.ModePerm)
	// Create the file
	zFile := "govuk/assets.zip"
	out, err := os.Create(zFile)
	if err != nil {
		slog.Error("error creating assets zip", slog.String("err", err.Error()))
		return
	}

	defer out.Close()
	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		slog.Error("error writing to zip", slog.String("err", err.Error()))
		return
	}
	// extract
	archive, err := zip.OpenReader(zFile)
	defer archive.Close()
	if err != nil {
		slog.Error("error opening zip", slog.String("err", err.Error()))
		return
	}

	for _, f := range archive.File {
		path := filepath.Join("govuk", f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
		} else {
			// make parent dir
			os.MkdirAll(filepath.Dir(path), os.ModePerm)
			dstFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			defer dstFile.Close()
			if err != nil {
				slog.Error("error opening file", slog.String("err", err.Error()))
				return
			}
			srcFile, err := f.Open()
			defer srcFile.Close()
			if err != nil {
				slog.Error("error opening file", slog.String("err", err.Error()))
				return
			}
			_, err = io.Copy(dstFile, srcFile)
			if err != nil {
				slog.Error("error opening file", slog.String("err", err.Error()))
				return
			}
			slog.Debug("copied file", slog.String("to", path))

		}
	}
	// remove the zip
	os.Remove(zFile)
}
