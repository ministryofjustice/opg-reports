package arguments

import (
	"database/sql"
	"fmt"
	"opg-reports/report/packages/types"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// Api is the main api argument struct that used in running
// the command and starting and handlers etc
//
// types.ApiArguments
type Api struct {
	DB      *DB
	Version *Versions
	Info    *ApiHost
}

func (self *Api) Database() types.DBer {
	return self.DB
}
func (self *Api) Versions() types.Versioner {
	return self.Version
}
func (self *Api) Host() types.Hoster {
	return self.Info
}

type ApiHost struct {
	host
}
type FrontHost struct {
	host
}

// Host is used to track the front & api host names
type host struct {
	Name     string `json:"label"`
	Hostname string `json:"host"`
}

func (self host) Label() string {
	return self.Name
}
func (self host) Host() string {
	return self.Hostname
}

// Versions contains semantic version information
//
// types.Versioner
type Versions struct {
	Version string `json:"semver"`
	SHA     string `json:"hash"`
}

func (self *Versions) Semver() string {
	return self.Version
}
func (self *Versions) Hash() string {
	return self.SHA
}

// DB struct for database connection details
//
// types.DBer
type DB struct {
	Driver   string `json:"driver"`
	Filepath string `json:"path"`
	Params   string `json:"params"`
}

// SetPath updates the database path to the location.
func (self *DB) SetPath(file string) {
	self.Filepath = file
}

// Connection create a connection - will also create the folder path
// to where the database is configured.
func (self *DB) Connection() (db *sql.DB, err error) {
	var connect = fmt.Sprintf("%s%s", self.Filepath, self.Params)

	os.MkdirAll(filepath.Dir(self.Filepath), os.ModePerm)

	db, err = sql.Open(self.Driver, connect)
	return
}
