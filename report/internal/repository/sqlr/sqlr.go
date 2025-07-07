package sqlr

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/interfaces"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

const label string = "sql-repository"

type Repository struct {
	ctx  context.Context
	conf *config.Config
	log  *slog.Logger
}

type RepositoryWithSelect[T interfaces.Model] struct {
	Repository
	ctx  context.Context
	conf *config.Config
	log  *slog.Logger
}

// empty is used for selects that require a non-nil interface
// passed
type empty struct{}

// Init is called via New and it creates database using the
// full SCHEMA
func (self *Repository) init() (err error) {
	_, err = self.Exec(SCHEMA)
	return
}

// connection internal helper to handle connecting to the db
func (self *Repository) connection() (db *sqlx.DB, err error) {
	var dbSource string = self.conf.Database.Source()

	// create the file path
	dir := filepath.Dir(self.conf.Database.Path)
	os.MkdirAll(dir, os.ModePerm)

	self.log.With("dbSource", dbSource).Debug("connecting to database ...")
	db, err = sqlx.ConnectContext(self.ctx, self.conf.Database.Driver, dbSource)
	if err != nil {
		self.log.Error("connection failed", "error", err.Error(), "dbSource", dbSource)
	}

	return
}

// New creates a new repo that can write to the database
func New(ctx context.Context, log *slog.Logger, conf *config.Config) (rp *Repository, err error) {

	if log == nil {
		err = fmt.Errorf("no logger passed for %s", label)
		return
	}
	if conf == nil {
		err = fmt.Errorf("no config passed for %s", label)
		return
	}

	rp = &Repository{
		ctx:  ctx,
		log:  log.WithGroup(label),
		conf: conf,
	}

	if !utils.FileExists(rp.conf.Database.Path) {
		log.With("database.path", rp.conf.Database.Path).Warn("Database not found, so creating")
		err = rp.init()
	}

	return
}

// NewWithSelect creates a typed (T) version which allows selects to return a slice of T[]
func NewWithSelect[T Model](ctx context.Context, log *slog.Logger, conf *config.Config) (rps *RepositoryWithSelect[T], err error) {

	if log == nil {
		err = fmt.Errorf("no logger passed for %s", label)
		return
	}
	if conf == nil {
		err = fmt.Errorf("no config passed for %s", label)
		return
	}

	rps = &RepositoryWithSelect[T]{
		Repository: Repository{
			ctx:  ctx,
			log:  log,
			conf: conf,
		},
		ctx:  ctx,
		log:  log.WithGroup(label),
		conf: conf,
	}

	if !utils.FileExists(rps.conf.Database.Path) {
		log.With("database.path", rps.conf.Database.Path).Warn("Database not found, so creating")
		err = rps.init()
	}

	return
}

func Default(ctx context.Context, log *slog.Logger, conf *config.Config) (rp *Repository) {
	rp, err := New(ctx, log, conf)
	if err != nil {
		log.Error("error with default", "err", err.Error())
	}
	return
}

func DefaultWithSelect[T Model](ctx context.Context, log *slog.Logger, conf *config.Config) (rp *RepositoryWithSelect[T]) {
	rp, err := NewWithSelect[T](ctx, log, conf)
	if err != nil {
		log.Error("error with default", "err", err.Error())
	}
	return
}
