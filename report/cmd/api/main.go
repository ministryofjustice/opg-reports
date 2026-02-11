package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"opg-reports/report/conf"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbmigrations"
	"opg-reports/report/internal/domain/accounts/accountapis/accountall"
	"opg-reports/report/internal/domain/codebases/codebaseapis/codebaseall"
	"opg-reports/report/internal/domain/codeowners/codeownerapis/codeownerall"
	"opg-reports/report/internal/domain/codeowners/codeownerapis/codeownerforteam"
	"opg-reports/report/internal/domain/infracosts/infracostapis/infracostsbymonthteam"
	"opg-reports/report/internal/domain/teams/teamapis/teamall"
	"opg-reports/report/internal/domain/uptime/uptimeapis/uptimebymonthteam"
	"opg-reports/report/internal/utils/logger"
	"os"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/jmoiron/sqlx"
)

const (
	cmdName   string = "api" // root command name
	shortDesc string = `api runs the main api command to start huma etc.`
	longDesc  string = `api runs the main api command to start huma api & docs endpoints.`
)

// config items
var (
	cfg *conf.Config    // default config
	ctx context.Context // default context
	log *slog.Logger    // default logger
)

// cli defaults
var (
	defaultDBPath   string = "./database/api.db"
	defaultDBDriver string = "sqlite3"
	defaultAddr     string = ":8081"
)

// apiOpts struct for running the main api cmd, fected from
// cli args
type apiOpts struct {
	Address string `doc:"host and pot to listen on" default:":8081"`
	DB      string `doc:"path to database file" default:"./database/api.db"` // database file path
	Driver  string `doc:"database driver type" default:"sqlite3"`            // database driver
}

func handlers(ctx context.Context, log *slog.Logger, opts *apiOpts, api huma.API) (err error) {
	var db *sqlx.DB
	// db connection to share with handlers
	db, err = dbconn(ctx, log, opts)
	if err != nil {
		return
	}

	// accounts
	accountall.Register(ctx, log, db, api)
	// codebases
	codebaseall.Register(ctx, log, db, api)
	// codeowners
	codeownerall.Register(ctx, log, db, api)
	codeownerforteam.Register(ctx, log, db, api)
	// infracosts
	infracostsbymonthteam.Register(ctx, log, db, api)
	// teams
	teamall.Register(ctx, log, db, api)
	// uptime
	uptimebymonthteam.Register(ctx, log, db, api)
	return
}

// runApiServer is the main function runner to start the api command
func runApiServer(ctx context.Context, log *slog.Logger) (err error) {

	var (
		humaapi       huma.API
		cli           humacli.CLI
		name          string         = "OPG Reports API"
		version       string         = "Test"
		mux           *http.ServeMux = http.NewServeMux()
		shutdownDelay time.Duration  = 5 * time.Second
		lg            *slog.Logger   = log.With("func", "api.runApiServer")
	)

	humaapi = humago.New(mux, huma.DefaultConfig(name, version))

	cli = humacli.New(func(hooks humacli.Hooks, opts *apiOpts) {
		var addr = opts.Address
		var server = http.Server{
			Addr:    addr,
			Handler: mux,
		}
		// add middleware
		addMiddleware(humaapi, log)
		// register handlers
		err = handlers(ctx, log, opts, humaapi)
		if err != nil {
			return
		}
		// startup
		hooks.OnStart(func() {
			lg.Info("Starting api server...")
			lg.Info(fmt.Sprintf("DB: %s", opts.DB))
			lg.Info(fmt.Sprintf("API: [http://%s/]", addr))
			lg.Info(fmt.Sprintf("Docs: [http://%s/docs]", addr))

			server.ListenAndServe()
		})
		// graceful shutdown
		hooks.OnStop(func() {
			lg.Info("Stopping api server...")
			ctx, cancel := context.WithTimeout(ctx, shutdownDelay)
			defer cancel()
			server.Shutdown(ctx)
		})
	})

	cli.Run()
	return
}

// addMiddleware add all middleware into the request; currently empty
func addMiddleware(hapi huma.API, log *slog.Logger) {
	hapi.UseMiddleware(func(ctx huma.Context, next func(huma.Context)) {
		next(ctx)
	})
}

// dbconn used to create a db connection or throw error, also runs migration
func dbconn(ctx context.Context, log *slog.Logger, opts *apiOpts) (db *sqlx.DB, err error) {
	var (
		driver string = defaultDBDriver
		path   string = defaultDBPath
	)
	//  replce defaults with input params
	if opts.Driver != "" {
		driver = opts.Driver
	}
	if opts.DB != "" {
		path = opts.DB
	}
	// db connection
	db, err = dbconnection.Connection(ctx, log, driver, conf.DBConnectionString(path, ""))
	if err == nil {
		err = dbmigrations.Migrate(ctx, log, db)
	}
	return
}

func main() {
	var err error

	cfg = conf.New()
	ctx = context.Background()
	log = logger.New(cfg.Log.Level, cfg.Log.Type)

	err = runApiServer(ctx, log)
	if err != nil {
		log.Error("error running command", "err", err.Error())
		os.Exit(1)
	}
}
