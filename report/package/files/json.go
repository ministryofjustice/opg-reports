package files

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
)

// ReadJSON
//
// If the file doesnt exist then nothing is returned (no error)
func ReadJSON[T any](ctx context.Context, file string, destination T) (err error) {
	var content []byte
	if !Exists(ctx, file) {
		return
	}
	if content, err = os.ReadFile(file); err != nil {
		return
	}
	err = json.Unmarshal(content, destination)
	return
}

func WriteAsJSON[T any](ctx context.Context, file string, src T) (err error) {
	var content []byte
	content, err = json.MarshalIndent(src, "", "  ")
	if err != nil {
		return
	}
	os.MkdirAll(filepath.Dir(file), os.ModePerm)
	err = os.WriteFile(file, content, os.ModePerm)
	return
}
