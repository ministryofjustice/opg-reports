package awsr

import (
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// GetCallerIdentity returns the current identity, whihc should include the account id
// and arn being used.
//
// Used to find current aws session account details to verify data sources and add account
// id details to cost calls etc
func (self *Repository) GetCallerIdentity(client ClientSTSCaller) (caller *sts.GetCallerIdentityOutput, err error) {
	return client.GetCallerIdentity(self.ctx, &sts.GetCallerIdentityInput{})
}
