package config

import (
	"fmt"

	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

// Database stores all the config values relating to the database connection
type Database struct {
	Driver string `json:"driver"`
	Path   string `json:"path"`
	Params string `json:"params"`
}

// Source returns the full connection string to use with the database drivers
func (self *Database) Source() (src string) {
	src = fmt.Sprintf("%s%s", self.Path, self.Params)

	return
}

// Github provides connection details to access github org
type Github struct {
	Organisation string `json:"organisation"`
	Token        string `json:"-"`
}

// Config is the overacrhing config item
type Config struct {
	Database *Database
	Github   *Github // this should always by empty for the running servers and and only used by the import commands
}

func NewDatabaseConfig() (db *Database) {
	db = &Database{
		Driver: utils.GetEnvVar("DB_DRIVER", "sqlite3"),
		Path:   utils.GetEnvVar("DB_PATH", "test.db"),
		Params: utils.GetEnvVar("DB_PARAMS", "?_journal=WAL&_busy_timeout=5000&_vacuum=incremental&_synchronous=NORMAL&_cache_size=1000000000"),
	}
	return
}

func NewGithubConfig() (gh *Github) {
	gh = &Github{
		Organisation: utils.GetEnvVar("GH_ORG", ""),
		Token:        utils.GetEnvVar("GH_TOKEN", ""),
	}
	return
}

func NewConfig() (cfg *Config) {
	cfg = &Config{
		Database: NewDatabaseConfig(),
		Github:   NewGithubConfig(),
	}
	return
}
