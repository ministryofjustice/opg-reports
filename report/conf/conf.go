package conf

import (
	"strings"

	"github.com/spf13/viper"
)

type config struct {
	DB     *DB     `mapstructure:"db"`
	Github *GitHub `mapstructure:"github"`
}

type DB struct {
	Driver string `mapstructure:"driver"`
	Path   string `mapstructure:"path"`
	Params string `mapstructure:"params"`
}

type GitHub struct {
	Token string `mapstructure:"token"`
}

var (
	vp   *viper.Viper             // viper instance
	name string       = "default" // name of config file to use
)

func setup(filename string) (v *viper.Viper, cfg *config, err error) {
	cfg = &config{}
	v = viper.New()
	// load from the default file
	v.SetConfigName(filename)
	v.AddConfigPath(".")
	// use underscores from env vars
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// load from evn
	v.AutomaticEnv()
	// read
	err = v.ReadInConfig()
	if err != nil {
		return
	}
	err = v.Unmarshal(&cfg)
	return
}

func New() (c *config) {
	vp, c, _ = setup(name)
	return
}
