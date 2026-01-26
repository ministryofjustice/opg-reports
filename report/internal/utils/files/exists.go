package files

import "os"

// FileExists checks if the file exists
func Exists(path string) bool {
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

// DirExists checks if the path exists and is a directory
func DirExists(path string) (exists bool) {
	exists = false
	info, err := os.Stat(path)

	if err == nil && info.IsDir() {
		exists = true
	}
	return
}
