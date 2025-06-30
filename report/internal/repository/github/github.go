package github

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ministryofjustice/opg-reports/report/config"
)

const label string = "github-repository"

type Repository struct {
	ctx  context.Context
	conf *config.Config
	log  *slog.Logger
}

func (self *Repository) Close() {}

// GetReleaseOptions used in release queries
type GetReleaseOptions struct {
	ExcludePrereleases bool // Exclude releases marked as prereleases
	ExcludeDraft       bool // Exclude anything marked as a draft
	ExcludeNoAssets    bool // Exclude anything that does not have assets
}

func New(ctx context.Context, log *slog.Logger, conf *config.Config) (rp *Repository, err error) {
	rp = &Repository{}

	if log == nil {
		err = fmt.Errorf("no logger passed for %s", label)
		return
	}
	if conf == nil {
		err = fmt.Errorf("no config passed for %s", label)
		return
	}

	log = log.WithGroup(label)
	rp = &Repository{
		ctx:  ctx,
		log:  log,
		conf: conf,
	}

	return
}
