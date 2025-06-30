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

type Repository[T interfaces.Model] struct {
	ctx  context.Context
	conf *config.Config
	log  *slog.Logger
}

// empty is used for selects that require a non-nil interface
// passed
type empty struct{}

// Init is called via New and it creates database using the
// full SCHEMA
func (self *Repository[T]) init() (err error) {
	_, err = self.Exec(SCHEMA)
	return
}

// connection internal helper to handle connecting to the db
func (self *Repository[T]) connection() (db *sqlx.DB, err error) {
	var dbSource string = self.conf.Database.Source()

	// create the file path
	dir := filepath.Dir(self.conf.Database.Path)
	os.MkdirAll(dir, os.ModePerm)

	self.log.With("dbSource", dbSource).Debug("connecting to database...")
	db, err = sqlx.ConnectContext(self.ctx, self.conf.Database.Driver, dbSource)
	if err != nil {
		self.log.Error("connection failed", "error", err.Error(), "dbSource", dbSource)
	}

	return
}

// New creates a new repo
func New[T interfaces.Model](ctx context.Context, log *slog.Logger, conf *config.Config) (rp *Repository[T], err error) {
	rp = &Repository[T]{}

	if log == nil {
		err = fmt.Errorf("no logger passed for %s", label)
		return
	}
	if conf == nil {
		err = fmt.Errorf("no config passed for %s", label)
		return
	}

	log = log.WithGroup(label)
	rp = &Repository[T]{
		ctx:  ctx,
		log:  log,
		conf: conf,
	}

	if !utils.FileExists(rp.conf.Database.Path) {
		log.With("database.path", rp.conf.Database.Path).Warn("Database not found, so creating")
		err = rp.init()
	}

	return
}
