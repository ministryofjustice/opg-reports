package config

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ministryofjustice/opg-reports/report/internal/utils"
	"github.com/spf13/viper"
)

// Config contains all config the values required in the applications
// grouped by type
type Config struct {
	Database *Database
	Github   *Github
	Versions *Versions
	Servers  *Servers
	Log      *Log
}

// Database stores all the config values relating to the database connection
type Database struct {
	Driver string
	Path   string
	Params string
}

// Source returns the full connection string to use with the database drivers
func (self *Database) Source() (src string) {
	src = fmt.Sprintf("%s%s", self.Path, self.Params)
	return
}

// Github provides connection details to access github org
type Github struct {
	Organisation string
	Token        string
}

// Log handles the slog setup used for the application
type Log struct {
	Level string
	Type  string
}

// Servers contains api & front end config
type Servers struct {
	Api   *Server
	Front *Server
}

// Server contains the address to use to contact the server
type Server struct {
	Name string
	Addr string
}

// Versions contains version data about the build
type Versions struct {
	Semver string
	Commit string
}

// setup a default config item that we use as a baseline
var defaultConfig = &Config{
	Database: &Database{
		Driver: "sqlite3",
		Path:   "./__database/api.db",
		Params: "?_journal=WAL&_busy_timeout=5000&_vacuum=incremental&_synchronous=NORMAL&_cache_size=1000000000",
	},
	Github: &Github{
		Organisation: "ministryofjustice",
		Token:        "",
	},
	Log: &Log{
		Level: "INFO",
		Type:  "TEXT",
	},
	Versions: &Versions{
		Semver: "0.0.0",
		Commit: "000000",
	},
	Servers: &Servers{
		Api:   &Server{Name: "OPG", Addr: "localhost:8081"},
		Front: &Server{Name: "OPG", Addr: "localhost:8080"},
	},
}

// NewViper configures an new viper instance with default values set from the defaultConfig struct
//
// Automatic environment variable mapping is enabled, so config values can be replace from the
// env directly
//
// Nest values need to use `_` notation in environment variables, so `DATABASE_PATH` - this is for
// easier mapping to AWS and similar env values. This also applies to `viper.Get` calls
func NewViper() (conf *viper.Viper) {
	conf = viper.New()

	conf.SetConfigType("json")
	conf.ReadConfig(bytes.NewBuffer(utils.MustMarshal(defaultConfig)))
	conf.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	conf.AutomaticEnv()

	return
}

// NewConfig uses NewViper to return a precomplied Config struct
// which will have values from the environment and the standard
// defaults.
func NewConfig() (cfg *Config) {
	cfg, _ = New()
	return
}

// New returns both a Config struct and the viper instance
// used in creation with it to allow cli flags etc to be
// bound afterwards
func New() (cfg *Config, vCfg *viper.Viper) {
	vCfg = NewViper()
	cfg = &Config{}

	vCfg.Unmarshal(cfg)
	return
}
