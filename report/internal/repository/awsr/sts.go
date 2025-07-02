package awsr

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// ClientSTS fetches a v2 version of the sts client loading from the
// env and setting the region
//
// Used to establish connection to the aws api for the sts / indentity calls
func ClientSTS(ctx context.Context, region string) (client *sts.Client, err error) {
	var awscfg aws.Config

	awscfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(region))
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
func (self *Repository) GetCallerIdentity(client ClientSTSCaller) (caller *sts.GetCallerIdentityOutput, err error) {
	return client.GetCallerIdentity(self.ctx, &sts.GetCallerIdentityInput{})
}
