package config

import (
	"database/sql"
	"fmt"
	"opg-reports/report/packages/env"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const govUKVersion string = `5.14.0`

type Config struct {
	// Sematic version number
	Semver string `json:"semver"`
	// Git commit sha
	SHA string `json:"sha"`
	// Location to use for the api server
	ApiHost string `json:"api_host"`
	// FrontHost is the location for the front host
	FrontHost string `json:"front_host"`
	// DB location
	DBPath string `json:"db"`
	// DB driver
	DBDriver string `json:"db_driver"`
	// DB connection parameters
	DBParams string `json:"db_params"`
	// govUKVersion number
	GovUKVersion string `json:"govuk_version"`
	// Root is the base directory for template and assets to be relative from
	Root string `json:"root_dir"`
	// localAssetsDir directory
	localAssetsDir string
	// govuk assets
	govukAssetsDir string
	// template directory
	templateAssetsDir string
}

// Connection returns a connection to the database
func (self *Config) Connection() (db *sql.DB) {
	os.MkdirAll(filepath.Dir(self.DBPath), os.ModePerm)
	conn := fmt.Sprintf("%s%s", self.DBPath, self.DBParams)
	driver := self.DBDriver
	db, _ = sql.Open(driver, conn)
	return
}

func (self *Config) Version() string {
	return fmt.Sprintf("%s (%s)", self.Semver, self.SHA)
}

func (self *Config) Conf() *Config {
	return self
}

// ApiHostname returns the hostname for the api
func (self *Config) ApiHostname() string {
	return self.ApiHost
}

// ConfigHostname
func (self *Config) FrontHostname() string {
	return self.FrontHost
}

// RootDir
func (self *Config) RootDir() string {
	return self.Root
}

// LocalAssetsDir provides updated path to the local asset directory
func (self *Config) LocalAssetsDir() string {
	return filepath.Join(self.Root, self.localAssetsDir)
}

// GovUKAssetsDir provides updated path to the govuk asset directory
func (self *Config) GovUKAssetsDir() string {
	return filepath.Join(self.Root, self.govukAssetsDir)
}

// TemplateAssetsDir provides updated path to the template asset directory
func (self *Config) TemplateAssetsDir() string {
	return filepath.Join(self.Root, self.templateAssetsDir)
}

// // DirAsURL converts a directory path into a url path to use in static handling
// func (self *Config) DirAsURL(dir string) (url string) {
// 	var root = self.RootDir()
// 	// strip root directory
// 	url = strings.ReplaceAll(dir, root, "")
// 	// trim slash from start & end
// 	url = strings.Trim(url, "/")
// 	// re-add start and end
// 	url = fmt.Sprintf(`/%s/`, url)
// 	// remove and doubles (for / etc)
// 	url = strings.ReplaceAll(url, `//`, `/`)

// 	return
// }

// Bind attachs the struct fields as cobra flags to the command passed
func (self *Config) BindFront(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&self.Semver, `semver`, self.Semver, `semver number.`)
	cmd.PersistentFlags().StringVar(&self.SHA, `sha`, self.SHA, `git commit hash.`)

	cmd.PersistentFlags().StringVar(&self.FrontHost, `front-host`, self.FrontHost, `location of front end server.`)
	cmd.PersistentFlags().StringVar(&self.ApiHost, `api-host`, self.ApiHost, `location of api server.`)

	cmd.PersistentFlags().StringVar(&self.Root, `root-dir`, self.Root, `root directory to use for file system.`)

	cmd.PersistentFlags().StringVar(&self.GovUKVersion, `govuk-version`, self.GovUKVersion, `GovUK version number.`)
}

// Bind attachs the struct fields as cobra flags to the command passed
func (self *Config) BindApi(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&self.Semver, `semver`, self.Semver, `semver number.`)
	cmd.PersistentFlags().StringVar(&self.SHA, `sha`, self.SHA, `git commit hash.`)
	cmd.PersistentFlags().StringVar(&self.ApiHost, `api-host`, self.ApiHost, `location of api server.`)

	cmd.PersistentFlags().StringVar(&self.DBPath, `db`, self.DBPath, `database path.`)
	cmd.PersistentFlags().StringVar(&self.DBDriver, `db-driver`, self.DBDriver, `database driver.`)
	cmd.PersistentFlags().StringVar(&self.DBParams, `db-params`, self.DBParams, `database connection parameters.`)
}

// NewFront create a default front instance, checks and replaces
// any values that are within the env and returns the value
func NewFront() (cfg *Config) {
	cfg = &Config{
		Semver:            `0.0.1`,
		SHA:               `abcdef`,
		FrontHost:         `:8080`,
		ApiHost:           `:8081`,
		GovUKVersion:      govUKVersion,
		Root:              `./`,
		localAssetsDir:    `web`,
		templateAssetsDir: `templates`,
		govukAssetsDir:    `govuk`,
	}
	// overwrite values from the os env
	env.OverwriteStruct(cfg)
	return
}

// NewApi create a default api instance, checks and replaces
// any values that are within the env and returns the value
func NewApi() (cfg *Config) {
	cfg = &Config{
		Semver:   `0.0.1`,
		SHA:      `abcdef`,
		ApiHost:  `:8081`,
		DBPath:   `./database/api.db`,
		DBDriver: `sqlite3`,
		DBParams: `?_journal=WAL&_busy_timeout=5000&_vacuum=incremental&_synchronous=NORMAL&_cache_size=1000000000`,
	}
	// overwrite values from the os env
	env.OverwriteStruct(cfg)
	return
}
