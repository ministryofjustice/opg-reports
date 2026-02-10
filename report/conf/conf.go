package conf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config is the main config struct and is generated
// from yaml file combined with environment variables
type Config struct {
	DB          *db    `json:"db"`
	Log         *log   `json:"log"`
	GithubToken string `json:"github_token" mapstructure:"github_token"`
}

// DB handles database related env & config values
type db struct {
	Driver string `json:"driver"`
	Path   string `json:"path"`
	Params string `json:"params"`
}

func (self *db) ConnectionString() (conn string) {
	return DBConnectionString(self.Path, self.Params)
}

func DBConnectionString(path string, params string) (conn string) {
	if params == "" {
		params = defaultParams
	}
	conn = fmt.Sprintf("%s%s", path, params)
	return
}

// log tracks logging configuration values used within
// all apps
type log struct {
	Level string `json:"level"`
	Type  string `json:"type"`
}

var (
	vp            *viper.Viper // viper instance
	defaultParams string       = "?_journal=WAL&_busy_timeout=5000&_vacuum=incremental&_synchronous=NORMAL&_cache_size=1000000000"
)

func defaults() (cfg *Config) {
	cfg = &Config{
		Log: &log{
			Level: "info",
			Type:  "json",
		},
		DB: &db{
			Driver: "sqlite3",
			Path:   "./api.db",
			Params: defaultParams,
		},
		GithubToken: "",
	}
	return
}

// setup loads default config reads that to viper and
// triggers env overwrites, then returns the unmarshaled
// value
func setup() (v *viper.Viper, cfg *Config, err error) {
	cfg = defaults()
	v = viper.New()
	// load from struct
	v.SetConfigType("json")
	v.ReadConfig(bytes.NewBuffer(mustMarshal(cfg)))
	// use underscores from env vars
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// load from env
	v.AutomaticEnv()
	err = v.Unmarshal(&cfg)
	return
}

// mustMarshal does a generic marshal on item to convert to json bytes
func mustMarshal[T any](item T) (bytes []byte) {
	bytes = []byte{}
	if b, err := json.MarshalIndent(item, "", "  "); err == nil {
		bytes = b
	}
	return
}

func New() (c *Config) {
	vp, c, _ = setup()
	return
}
