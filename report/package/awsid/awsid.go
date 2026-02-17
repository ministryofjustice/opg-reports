package awsid

import (
	"context"
	"opg-reports/report/package/awsclients"

	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// Identity returns the general sts caller identity details for the session & region
func Identity(ctx context.Context, region string) (id *sts.GetCallerIdentityOutput) {
	var err error
	var client *sts.Client

	client, err = awsclients.New[*sts.Client](ctx, region)
	if err != nil {
		return
	}

	id, err = client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return
	}
	return

}

// AccountID returns just the account id detials from the session information
func AccountID(ctx context.Context, region string) (account string) {

	var id *sts.GetCallerIdentityOutput

	id = Identity(ctx, region)
	if id != nil {
		account = *id.Account
	}

	return
}
