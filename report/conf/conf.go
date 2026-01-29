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
	DB     *db     `json:"db"`
	Log    *log    `json:"log"`
	Github *github `json:"github"`
	AWS    *aws    `json:"aws"`

	Accounts *accounts `json:"accounts"`
	Teams    *teams    `json:"teams"`
}

// DB handles database related env & config values
type db struct {
	Driver string `json:"driver"`
	Path   string `json:"path"`
	Params string `json:"params"`
}

func (self *db) ConnectionString() (conn string) {
	conn = fmt.Sprintf("%s%s", self.Path, self.Params)
	return
}

// accounts contains data relating to the accounts domain
// and how to find / fetch that
type accounts struct {
	Release string `json:"release"`
}

// teams contains data relating to the accounts domain
// and how to find / fetch that
type teams struct {
	Release string `json:"release"`
}

// github handles env token for access to github
// during data imports
type github struct {
	Token  string `json:"token"`
	Org    string `json:"org"`
	Parent string `json:"parent"`
}

// aws handles env token / session access to aws
// during data imports
type aws struct {
	Region  string      `json:"region"`
	Default *awsDefault `json:"default"`
	Session *awsSession `json:"session"`
}
type awsDefault struct {
	Region string `json:"region"`
}
type awsSession struct {
	Token string `json:"token"`
}

// log tracks logging configuration values used within
// all apps
type log struct {
	Level string `json:"level"`
	Type  string `json:"type"`
}

var vp *viper.Viper // viper instance

func defaults() (cfg *Config) {
	cfg = &Config{
		Log: &log{
			Level: "info",
			Type:  "json",
		},
		DB: &db{
			Driver: "sqlite3",
			Path:   "./api.db",
			Params: "?_journal=WAL&_busy_timeout=5000&_vacuum=incremental&_synchronous=NORMAL&_cache_size=1000000000",
		},
		Github: &github{
			Token:  "",
			Parent: "opg",
			Org:    "ministryofjustice",
		},
		AWS: &aws{
			Region:  "eu-west-1",
			Default: &awsDefault{Region: "eu-west-1"},
			Session: &awsSession{Token: ""},
		},
		Accounts: &accounts{
			Release: "v0.1.26",
		},
		Teams: &teams{
			Release: "v0.1.26",
		},
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
