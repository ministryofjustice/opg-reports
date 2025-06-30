package id

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/ministryofjustice/opg-reports/report/config"
)

// Repository
//
// interfaces:
//   - Repository
type Repository struct {
	ctx  context.Context
	conf *config.Config
	log  *slog.Logger
}

// connection is an internal helper to handle creating the client
func (self *Repository) connection() (client *sts.STS, err error) {
	var (
		sess           *session.Session
		sessionOptions *aws.Config = &aws.Config{Region: aws.String(self.conf.Aws.GetRegion())}
	)
	// create a new session
	sess, err = session.NewSession(sessionOptions)
	if err != nil {
		return
	}
	client = sts.New(sess)

	return
}

func (self *Repository) GetCallerIdentity() (id *sts.GetCallerIdentityOutput, err error) {
	var (
		client *sts.STS
		log    = self.log.With("operation", "GetCallerIdentity")
	)
	log.Debug("creating sts client ...")

	client, err = self.connection()
	if err != nil {
		return
	}

	log.Debug("getting caller identity details ...")
	id, err = client.GetCallerIdentityWithContext(self.ctx, &sts.GetCallerIdentityInput{})
	return
}

func (self *Repository) GetAccountID() (accountID string, err error) {
	var (
		id  *sts.GetCallerIdentityOutput
		log = self.log.With("operation", "GetAccountID")
	)
	log.Debug("getting call account id ...")
	id, err = self.GetCallerIdentity()
	if err != nil {
		return
	}
	// set the account
	accountID = *id.Account
	return
}

// New provides a configured repository instance
func New(ctx context.Context, log *slog.Logger, conf *config.Config) (rp *Repository, err error) {
	rp = &Repository{}

	if log == nil {
		err = fmt.Errorf("no logger passed for id repository")
		return
	}
	if conf == nil {
		err = fmt.Errorf("no config passed for id repository")
		return
	}

	log = log.WithGroup("sts")
	rp = &Repository{
		ctx:  ctx,
		log:  log,
		conf: conf,
	}

	return
}
