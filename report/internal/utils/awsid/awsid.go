package awsid

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/utils/awsclients"

	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// Identity returns the general sts caller identity details for the session & region
func Identity(ctx context.Context, log *slog.Logger, region string) (id *sts.GetCallerIdentityOutput) {
	var err error
	var client *sts.Client
	var lg *slog.Logger = log.With("func", "awsid.Identity")

	lg.Debug("starting ...")

	client, err = awsclients.New[*sts.Client](ctx, log, region)
	if err != nil {
		lg.Error("failed to create sts client.")
		return
	}

	id, err = client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		lg.Error("failed to get identity details.")
		return
	}
	lg.Debug("complete.")
	return

}

// AccountID returns just the account id detials from the session information
func AccountID(ctx context.Context, log *slog.Logger, region string) (account string) {
	var lg *slog.Logger = log.With("func", "awsid.AccountID")
	var id *sts.GetCallerIdentityOutput
	lg.Debug("starting ...")
	id = Identity(ctx, log, region)
	if id != nil {
		account = *id.Account
	}
	lg.Debug("complete.")
	return
}
