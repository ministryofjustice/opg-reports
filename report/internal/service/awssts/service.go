package awssts

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/interfaces"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/caller"
)

// Service is used to download, covnert and return data files from within s3 buckets.
//
// interfaces:
//   - Service
type Service[T interfaces.Model] struct {
	ctx       context.Context
	log       *slog.Logger
	conf      *config.Config
	store     interfaces.STSRepository
	directory string
}

// Close cleans up
func (self *Service[T]) Close() (err error) {
	return
}

func (self *Service[T]) GetAccountID() (account string) {
	var log *slog.Logger = self.log.With("operation", "GetAccountID")
	log.Debug("getting account details from current current caller identity ...")
	account = self.store.GetAccountID()
	return
}

// NewService returns a configured s3 service object
func NewService[T interfaces.Model](ctx context.Context, log *slog.Logger, conf *config.Config, store interfaces.STSRepository) (srv *Service[T], err error) {
	if log == nil {
		return nil, fmt.Errorf("no logger passed for sts service")
	}
	if conf == nil {
		return nil, fmt.Errorf("no config passed for sts service")
	}
	if conf.Aws == nil ||
		conf.Aws.Session == nil {
		return nil, fmt.Errorf("missing aws config details for sts service")
	}
	if conf.Aws.Region == "" ||
		conf.Aws.Session.Token == "" {
		return nil, fmt.Errorf("missing aws config details for sts service")
	}

	if store == nil {
		return nil, fmt.Errorf("no repository passed for sts service")
	}

	srv = &Service[T]{
		ctx:   ctx,
		log:   log.With("service", "sts"),
		conf:  conf,
		store: store,
	}
	return
}

// Default generates the default gh repository and then the service
func Default[T interfaces.Model](ctx context.Context, log *slog.Logger, conf *config.Config) (srv *Service[T]) {

	store, err := caller.New(ctx, log, conf)
	if err != nil {
		log.Error("error creating awssts repository for sts service", "error", err.Error())
		return nil
	}
	srv, err = NewService[T](ctx, log, conf, store)
	if err != nil {
		log.Error("error creating sts service", "error", err.Error())
		return nil
	}

	return
}
