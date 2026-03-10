package files

import (
	"context"
	"os"
)

// FileExists checks if the file exists
func Exists(ctx context.Context, path string) bool {
	info, err := os.Stat(path)

	// if there is an error, or the filepath doesnt exist, return false
	if err != nil || os.IsNotExist(err) {
		return false
	}
	// return false for directories - its not a file
	if info.IsDir() {
		return false
	}

	return true
}
