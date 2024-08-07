package files

import (
	"embed"
	"fmt"
	"opg-reports/shared/logger"
	"os"
	"testing"
)

func TestSharedFilesWriteFileError(t *testing.T) {
	logger.LogSetup()
	var err error
	td := os.TempDir()
	tDir, _ := os.MkdirTemp(td, "files-write-err-*")
	defer os.RemoveAll(tDir)
	// error making a directory as its reserved
	os.Mkdir(tDir+"/readonly", os.ModePerm)
	os.Chmod(tDir+"/readonly", 0444)
	text := `{"hello":"world"}`
	err = WriteFile(tDir+"/readonly/sub", "test.json", []byte(text))
	if err == nil {
		t.Errorf("should not be nil")
	}
	err = os.Chmod(tDir+"/readonly", 0777)

}

//go:embed testdata/*
var testFS embed.FS

func TestSharedFilesAllFromEmbedded(t *testing.T) {
	logger.LogSetup()
	tDir := "./testdata"
	fSys := NewFS(testFS, tDir)

	// these will be ignored as embed is done at build
	for i := 0; i < 8; i++ {
		f, _ := os.CreateTemp(tDir, "dummy-*.json")
		defer os.Remove(f.Name())
	}

	all := All(fSys, false)
	if len(all) != 2 {
		t.Errorf("incorrect number of files found: %d", len(all))
	}
	all = All(fSys, true)
	if len(all) != 1 {
		t.Errorf("incorrect number of json files found")
	}
}

func TestSharedFilesReadFromEmbedded(t *testing.T) {
	logger.LogSetup()
	tDir := "testdata"
	fSys := NewFS(testFS, tDir)

	text := `{"id": "001"}` + "\n"
	all := All(fSys, true)
	first := all[0]
	content, err := ReadFile(fSys, first)
	if err != nil {
		t.Errorf("error reading file: %v", err.Error())
	}
	if string(content) != text {
		t.Errorf("content mismtach [%v] [%v]", content, text)
	}

	err = SaveFile(fSys, first, []byte(text))
	if err != nil {
		t.Errorf("error saving file: %v", err.Error())
	}

}

func TestSharedFilesAllFromDir(t *testing.T) {
	logger.LogSetup()
	td := os.TempDir()
	tDir, _ := os.MkdirTemp(td, "files-all-*")
	dfSys := os.DirFS(tDir).(IReadFS)
	defer os.RemoveAll(tDir)
	for i := 0; i < 8; i++ {
		os.CreateTemp(tDir, "dummy-*.json")
	}
	for i := 0; i < 2; i++ {
		os.CreateTemp(tDir, "dummy-*.txt")
	}
	fSys := NewFS(dfSys, tDir)
	all := All(fSys, false)

	if len(all) != 10 {
		t.Errorf("incorrect number of files found")
	}

	filtered := Filter(all, `txt$`)
	if len(filtered) != 2 {
		t.Errorf("failed to filter")
		fmt.Printf("%+v\n", filtered)
	}

	all = All(fSys, true)
	if len(all) != 8 {
		t.Errorf("incorrect number of json files found")
	}

}

func TestSharedFilesWriteReadFromDir(t *testing.T) {
	logger.LogSetup()
	td := os.TempDir()
	tDir, _ := os.MkdirTemp(td, "files-write-read-*")
	dfSys := os.DirFS(tDir).(IReadFS)
	defer os.RemoveAll(tDir)

	os.CreateTemp(tDir, "dummy-*.json")
	fSys := NewFS(dfSys, tDir)
	all := All(fSys, true)

	if len(all) != 1 {
		t.Errorf("length error")
	}

	first := all[0]
	text := `{"hello":"world"}`
	err := WriteFile(tDir, first.Name(), []byte(text))
	if err != nil {
		t.Errorf("error writing to file")
	}

	content, err := ReadFile(fSys, first)
	if err != nil {
		t.Errorf("error reading file")
	}
	if string(content) != text {
		t.Errorf("content mismtach")
	}

	newTxt := `{"foo":"bar"}`
	err = SaveFile(fSys, first, []byte(newTxt))
	if err != nil {
		t.Errorf("error saving file")
	}
}
