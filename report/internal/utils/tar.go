package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func addToArchive(tw *tar.Writer, filename string) error {
	// Open the file which will be written into the archive
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get FileInfo about our file providing file size, mode, etc.
	info, err := file.Stat()
	if err != nil {
		return err
	}

	// Create a tar Header from the FileInfo data
	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	// Use full path as name (FileInfoHeader only takes the basename)
	// If we don't do this the directory strucuture would
	// not be preserved
	// https://golang.org/src/archive/tar/common.go?#L626
	header.Name = filename

	// Write file header to the tar archive
	err = tw.WriteHeader(header)
	if err != nil {
		return err
	}

	// Copy file content to tar archive
	_, err = io.Copy(tw, file)
	if err != nil {
		return err
	}

	return nil
}

// TarGzCreate takes a slice of filenames and generates a
// tar.gz archive from those at the location specificed by
// writeTo.
//
// If the writeTo location exists, or any of the files cannot
// be found/opened then an error is returned.
func TarGzCreate(writeTo string, files []string) (err error) {
	var (
		fp    *os.File
		gzipw *gzip.Writer
		tarw  *tar.Writer

		parentDir string = filepath.Dir(writeTo)
	)
	// if destination exists, return error so we dont overwrite
	if FileExists(writeTo) {
		return fmt.Errorf("file [%s] already exists, won't overwrite", writeTo)
	}
	// create the parent directory path
	err = os.MkdirAll(parentDir, os.ModePerm)
	// create the file
	fp, err = os.Create(writeTo)
	if err != nil {
		return
	}
	// now create the archives & close them via defer
	gzipw = gzip.NewWriter(fp)
	tarw = tar.NewWriter(gzipw)
	defer gzipw.Close()
	defer tarw.Close()

	for _, file := range files {
		err = addToArchive(tarw, file)
		if err != nil {
			return
		}
	}

	return
}

// TarGzExtract extract the tar.gz file passed into the directory passed.
//
// The extractTo path is appended to the filepaths reported by the tar.gz
// structure to allow extraction to a fixed location
func TarGzExtract(extractTo string, source io.Reader) (err error) {
	var (
		uncompressed *gzip.Reader
		tarReader    *tar.Reader
	)
	uncompressed, err = gzip.NewReader(source)
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
			// create parent directory as well in case we get to a
			// file before a folder
			parent := filepath.Dir(header.Name)
			os.MkdirAll(parent, os.ModePerm)

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
