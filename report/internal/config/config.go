package config

import (
	"database/sql"
	"fmt"
	"html/template"
	"opg-reports/report/packages/env"
	"opg-reports/report/packages/tmpl"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const govUKVersion string = `5.14.0`

type Config struct {
	IsJSON bool `json:"is_json"`
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
	LocalAssetPath string
	// govuk assets
	GovUKAssetPath string
	// template directory
	TemplateAssetPath string
}

// Connection returns a connection to the database
func (self *Config) Connection() (db *sql.DB) {
	os.MkdirAll(filepath.Dir(self.DBPath), os.ModePerm)
	conn := fmt.Sprintf("%s%s", self.DBPath, self.DBParams)
	driver := self.DBDriver
	db, _ = sql.Open(driver, conn)
	return
}

// Version returns the version signature of the semver & sha
func (self *Config) Version() string {
	return self.Semver
}

func (self *Config) GovukVersion() string {
	return self.GovUKVersion
}

// Directories returns the updated paths for:
//   - root
//   - local
//   - govuk
//   - templates
func (self *Config) Directories() map[string]string {
	var root = self.Root
	return map[string]string{
		"root":      root,
		"local":     filepath.Join(root, self.LocalAssetPath),
		"govuk":     filepath.Join(root, self.GovUKAssetPath),
		"templates": filepath.Join(root, self.TemplateAssetPath),
	}
}

// Template creates a complied template with functions and files ready to be
// used by a response
//
// Only returns values when the template directory is set
func (self *Config) Template(name string) (t *template.Template, err error) {
	var files []string

	// if this json only, then return nil so triggers json rendering
	if self.IsJSON {
		fmt.Println("json, so returning nothing")
		return nil, nil
	}
	t = template.New(name).Funcs(tmpl.Functions())
	if self.TemplateAssetPath != "" {
		files = tmpl.Files(filepath.Join(self.Root, self.TemplateAssetPath))
		if len(files) > 0 {
			t, err = t.ParseFiles(files...)
		}
	}
	return
}

// ApiHostname returns the hostname for the api
func (self *Config) ApiHostname() string {
	return self.ApiHost
}

// ConfigHostname
func (self *Config) FrontHostname() string {
	return self.FrontHost
}

// BindFront attaches the struct fields as cobra flags to the command passed
func (self *Config) BindFront(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&self.Semver, `semver`, self.Semver, `semver number.`)
	cmd.PersistentFlags().StringVar(&self.SHA, `sha`, self.SHA, `git commit hash.`)
	cmd.PersistentFlags().StringVar(&self.FrontHost, `front-host`, self.FrontHost, `location of front end server.`)
	cmd.PersistentFlags().StringVar(&self.ApiHost, `api-host`, self.ApiHost, `location of api server.`)
	cmd.PersistentFlags().StringVar(&self.Root, `root-dir`, self.Root, `root directory to use for file system.`)
	cmd.PersistentFlags().StringVar(&self.GovUKVersion, `govuk-version`, self.GovUKVersion, `GovUK version number.`)
}

// BindApi attaches the struct fields as cobra flags to the command passed
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
		IsJSON:            false,
		Semver:            `0.0.1`,
		SHA:               `abcdef`,
		FrontHost:         `:8080`,
		ApiHost:           `:8081`,
		GovUKVersion:      govUKVersion,
		Root:              `./`,
		LocalAssetPath:    `web`,
		TemplateAssetPath: `templates`,
		GovUKAssetPath:    `govuk`,
	}
	// overwrite values from the os env
	env.OverwriteStruct(cfg)
	return
}

// NewApi create a default api instance, checks and replaces
// any values that are within the env and returns the value
func NewApi() (cfg *Config) {
	cfg = &Config{
		IsJSON:   true,
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
