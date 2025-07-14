package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config contains all config the values required in the applications
// grouped by type
type Config struct {
	Database *Database // Database configuration values, such as filesystem location and connection flags
	Github   *Github   // Github configuration values used for accessing github data (such as releases and assets)
	Aws      *Aws      // AWS values for capturing environment authentication values (session token etc)
	Metadata *Metadata // Metadata contains details on where opg-metadata is stored for importing accounts & team records
	Existing *Existing // Existing contains information about where existing data is stored for the import command
	GovUK    *GovUK    // GovUK contains details about the gov uk front end assets and version to use
	Servers  *Servers  // Servers contains details about front & api server configuration (address, name etc)
	Versions *Versions // Versions contains semver and commit references and used for output
	Log      *Log      // Log contains settings that can be overridden by env values for LOG_LEVEL (warn, info, debug etc) and LOG_TYPE (text / json)
}

// Database stores all the config values relating to the database connection
type Database struct {
	Driver string // env: DATABASE_DRIVER
	Path   string // env: DATABASE_PATH
	Params string // env: DATABASE_PARAMS
}

// Source returns the full connection string to use with the database drivers
func (self *Database) Source() (src string) {
	src = fmt.Sprintf("%s%s", self.Path, self.Params)
	return
}

// Log handles the slog setup used for the application
type Log struct {
	Level string // env: LOG_LEVEL
	Type  string // env: LOG_TYPE
}

// Servers contains api & front end config
type Servers struct {
	Api   *Server
	Front *Server
}

// Server contains the address to use to contact the server
type Server struct {
	Name    string        // env: $X_NAME
	Addr    string        // env: $X_ADDR
	Timeout time.Duration // env: $X_TIMEOUT - used to set the http client timeout
}

// Versions contains version data about the build
type Versions struct {
	Semver string // env: VERSIONS_SEMVER
	Commit string // env: VERSIONS_COMMIT
}

func (self *Versions) Signature() string {
	return fmt.Sprintf("%s (%s)", self.Semver, self.Commit)
}

type GovUK struct {
	Front *frontEnd
}
type frontEnd struct {
	Owner      string // env: GOVUK_FRONT_OWNER
	Repository string // env: GOVUK_FRONT_REPOSITORY
	ReleaseTag string // env: GOVUK_FRONT_RELEASETAG
	AssetName  string // env: GOVUK_FRONT_ASSETNAME
	UseRegex   bool

	Directory string // env: GOVUK_FRONT_DIRECTORY - where downloaded gov uk assets are stored locally
}

// Metadata provides details on where metadata information used
// for generating team / account information is stored
type Metadata struct {
	Owner      string // env: METADATA_OWNER - defaults to ministryofjustice
	Repository string // env: METADATA_REPOSITORY
	ReleaseTag string // env: GOVUK_FRONT_RELEASETAG
	AssetName  string // env: GOVUK_FRONT_ASSETNAME
	UseRegex   bool
}

// Github provides connection details to access github org
type Github struct {
	Organisation string // env: GITHUB_ORGANISATION - defaults to ministryofjustice
	Token        string // env: GITHUB_TOKEN
}

type Existing struct {
	Costs *bucketInfo // env: EXISTING_COSTS_${X}
	DB    *bucketInfo // env: EXISTING_DB_${X}
}

// AWS tracks environment information and details on where
type Aws struct {
	BillingDate int    // env: AWS_BILLINGDATE - used to determin month data we can show
	Region      string // env: AWS_REGION
	Default     *def
	Session     *session
}

func (self *Aws) GetRegion() string {
	if self.Region != "" {
		return self.Region
	} else if self.Default.Region != "" {
		return self.Default.Region
	}
	return ""
}

func (self *Aws) GetToken() string {
	return self.Session.Token
}

// region
type def struct {
	Region string // env: AWS_DEFAULT_REGION
}

// session
type session struct {
	Token string // env: AWS_SESSION_TOKEN
}

// bucketInfo used to store details of a bucket, a path in a bucket or an item in a bucket
type bucketInfo struct {
	Bucket string // env: EXISTING_${X}_BUCKET
	Prefix string // env: EXISTING_${X}_PREFIX
	Key    string // env: EXISTING_${X}_KEY
}

func (self *bucketInfo) Path() string {
	return fmt.Sprintf("%s%s", self.Prefix, self.Key)
}

// setup a default config item that we use as a baseline
var defaultConfig = &Config{
	Database: &Database{
		Driver: "sqlite3",
		Path:   "./api.db",
		Params: "?_journal=WAL&_busy_timeout=5000&_vacuum=incremental&_synchronous=NORMAL&_cache_size=1000000000",
	},
	Github: &Github{
		Organisation: "ministryofjustice", // default organisations
		Token:        "",                  // needed for tests & data imports
	},
	Aws: &Aws{
		BillingDate: 15,
		Region:      "",
		Default:     &def{Region: "eu-west-1"},
		Session:     &session{Token: ""},
	},
	Metadata: &Metadata{
		Owner:      "ministryofjustice",
		Repository: "opg-metadata", // repository name for where meta data info is
		AssetName:  "metadata(.*)gz",
		ReleaseTag: "v(.*)",
		UseRegex:   true,
	},
	GovUK: &GovUK{
		Front: &frontEnd{
			Owner:      "alphagov",
			Repository: "govuk-frontend",
			AssetName:  "release-v5.11.0.zip",
			ReleaseTag: "v5.11.0",
			UseRegex:   false,
			Directory:  "govuk",
		},
	},
	Existing: &Existing{
		Costs: &bucketInfo{Bucket: "report-data-development", Prefix: "aws_costs/"},
		DB:    &bucketInfo{Bucket: "report-data-development", Prefix: "", Key: "database/api.db"},
	},
	Log: &Log{
		Level: "INFO",
		Type:  "JSON",
	},
	Versions: &Versions{
		Semver: "0.0.0",
		Commit: "000000",
	},
	Servers: &Servers{
		Api:   &Server{Name: "OPG Reports API", Addr: "localhost:8081"},
		Front: &Server{Name: "OPG Reports", Addr: "localhost:8080", Timeout: (30 * time.Second)},
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
	conf.ReadConfig(bytes.NewBuffer(mustMarshal(defaultConfig)))
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

// mustMarshal does a generic marshal on item to convert to json bytes
func mustMarshal[T any](item T) (bytes []byte) {
	bytes = []byte{}
	if b, err := json.MarshalIndent(item, "", "  "); err == nil {
		bytes = b
	}
	return
}
