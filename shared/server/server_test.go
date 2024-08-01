package server

import (
	"opg-reports/internal/testhelpers"
	"opg-reports/shared/data"
	"opg-reports/shared/files"
	"opg-reports/shared/logger"
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
	logger.LogSetup()
	td := os.TempDir()
	tDir, _ := os.MkdirTemp(td, "server-test-*")
	defer os.RemoveAll(tDir)
	dfSys := os.DirFS(tDir).(files.IReadFS)

	f := files.NewFS(dfSys, tDir)
	store := data.NewStore[*tEntry]()
	resp := response.NewResponse[response.ICell, response.IRow[response.ICell]]()

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

func TestSharedServerApiGetParams(t *testing.T) {
	logger.LogSetup()
	td := os.TempDir()
	tDir, _ := os.MkdirTemp(td, "server-test-*")
	defer os.RemoveAll(tDir)
	dfSys := os.DirFS(tDir).(files.IReadFS)

	f := files.NewFS(dfSys, tDir)
	store := data.NewStore[*tEntry]()
	resp := response.NewResponse[response.ICell, response.IRow[response.ICell]]()

	tApi := NewApi(store, f, resp)
	_, r := testhelpers.WRGet("/test/?team=dev&team=ops&not-allowed=1")

	v := tApi.GetParameters([]string{"team"}, r)
	if len(v) != 1 {
		t.Errorf("get param logic failed: %v ", v)
	}
	if _, ok := v["team"]; !ok {
		t.Errorf("get param logic failed: %v ", v)
	}
}
