package server

import (
	"opg-reports/shared/data"
	"opg-reports/shared/files"
	"opg-reports/shared/server/response"
	"os"
	"testing"
	"time"
)

type tEntry struct {
	Id       string `json:"id"`
	Tag      string `json:"tag"`
	Category string `json:"category"`
}

func (i *tEntry) UID() string {
	return i.Id
}
func (i *tEntry) TS() time.Time {
	return time.Now().UTC()
}
func (i *tEntry) Valid() bool {
	return true
}
func TestSharedServerApi(t *testing.T) {
	td := os.TempDir()
	tDir, _ := os.MkdirTemp(td, "server-test-*")
	defer os.RemoveAll(tDir)
	dfSys := os.DirFS(tDir).(files.IReadFS)

	f := files.NewFS(dfSys, tDir)
	store := data.NewStore[*tEntry]()
	resp := response.NewResponse[*response.Cell, *response.Row[*response.Cell]]()

	tApi := NewApi(store, f, resp)

	if tApi.Store() != store {
		t.Errorf("store mismatch")
	}
	if tApi.FS() != f {
		t.Errorf("fs mismatch")
	}

	if tApi.GetResponse() != resp {
		t.Errorf("response mismatch")
	}
}
