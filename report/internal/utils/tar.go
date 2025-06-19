package utils

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

// TarGzExtract extract the tar.gz file passed into the directory passed.
//
// The extractTo path is appended to the filepaths reported by the tar.gz
// structure to allow extraction to a fixed location
func TarGzExtract(targz io.Reader, extractTo string) (err error) {
	var (
		uncompressed *gzip.Reader
		tarReader    *tar.Reader
	)
	uncompressed, err = gzip.NewReader(targz)
	if err != nil {
		return
	}

	// close file
	defer uncompressed.Close()
	tarReader = tar.NewReader(uncompressed)

	for true {
		var header *tar.Header

		// grab next & handle errors
		header, err = tarReader.Next()
		// end of tar ball - not actually an error
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return
		}
		// add prefix to extraction to push into another place
		header.Name = filepath.Join(extractTo, header.Name)

		switch header.Typeflag {
		// directories
		case tar.TypeDir:
			err = os.MkdirAll(header.Name, os.ModePerm)
			if err != nil {
				return
			}
		// files
		case tar.TypeReg:
			var outFile *os.File
			// create file
			outFile, err = os.Create(header.Name)
			if err != nil {
				return
			}
			// copy into file
			_, err = io.Copy(outFile, tarReader)
			if err != nil {
				return
			}
			// close the file
			outFile.Close()
		}
	}

	return
}
