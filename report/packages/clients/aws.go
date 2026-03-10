package clients

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// AccountID returns just the account id detials from the session information
func AWSAccountID(ctx context.Context, region string) (account string) {
	var err error
	var client *sts.Client
	var id *sts.GetCallerIdentityOutput

	client, err = New[*sts.Client](ctx, region)
	if err != nil {
		return
	}

	id, err = client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return
	}

	if id != nil {
		account = *id.Account
	}

	return
}
