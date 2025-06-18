package config

import (
	"fmt"

	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

type Database struct {
	Driver string `json:"driver"`
	Path   string `json:"path"`
	Params string `json:"params"`
}

func (self *Database) Source() string {
	return fmt.Sprintf("%s%s", self.Path, self.Params)
}

type Config struct {
	Database *Database
}

func NewDatabaseConfig() (db *Database) {
	db = &Database{
		Driver: utils.GetEnvVar("DB_DRIVER", "sqlite3"),
		Path:   utils.GetEnvVar("DB_PATH", "test.db"),
		Params: utils.GetEnvVar("DB_PARAMS", "?_journal=WAL&_busy_timeout=5000&_vacuum=incremental&_synchronous=NORMAL&_cache_size=1000000000"),
	}
	return
}

func NewConfig() (cfg *Config) {
	cfg = &Config{
		Database: NewDatabaseConfig(),
	}
	return
}
