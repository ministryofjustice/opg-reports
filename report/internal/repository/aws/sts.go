package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// getSTSClient fetches a v2 version of the sts client loading from the
// env and setting the region
//
// Used to establish connection to the aws api for the sts / indentity calls
func (self *Repository) getSTSClient(region string) (client *sts.Client, err error) {
	var awscfg aws.Config

	awscfg, err = config.LoadDefaultConfig(self.ctx, config.WithRegion(region))
	if err != nil {
		return
	}
	client = sts.NewFromConfig(awscfg)
	return

}

// GetCallerIdentity returns the current identity, whihc should include the account id
// and arn being used.
//
// Used to find current aws session account details to verify data sources and add account
// id details to cost calls etc
func (self *Repository) GetCallerIdentity() (caller *sts.GetCallerIdentityOutput, err error) {
	var client *sts.Client
	if client, err = self.getSTSClient(self.conf.Aws.GetRegion()); err != nil {
		self.log.With("operation", "GetCallerIdentity").Error("failed to get client", "err", err.Error())
		return
	}
	return client.GetCallerIdentity(self.ctx, &sts.GetCallerIdentityInput{})
}
