package files

import (
	"io"
	"strings"
)

// Lines reads the content from the buffer and returns the content split by new lines
func Lines(content io.ReadCloser) (lines []string) {
	var err error

	if content == nil {
		return
	}
	b, err := io.ReadAll(content)
	if err != nil {
		return
	}
	err = content.Close()
	if err != nil {
		return
	}
	// trim the last new line from the file
	lines = strings.Split(strings.TrimRight(string(b), "\n"), "\n")
	return
}
