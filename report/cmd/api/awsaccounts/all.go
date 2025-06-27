package awsaccounts

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/awsaccount"
)

func RegisterGetAwsAccountsAll(log *slog.Logger, conf *config.Config, api huma.API) {
	var operation = huma.Operation{
		OperationID:   "get-awsaccounts-all",
		Method:        http.MethodGet,
		Path:          "/v1/awsaccounts/all",
		Summary:       "All AWS accounts",
		Description:   "Returns a list of all AWS accounts stored.",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"AWS Accounts"},
	}
	huma.Register(api, operation, func(ctx context.Context, input *struct{}) (*GetAwsAccountsAllResponse, error) {
		return handleGetAwsAccountsAll(ctx, log, conf, input)
	})
}

// handleGetAwsAccountsAll deals with each request and fetches
func handleGetAwsAccountsAll(ctx context.Context, log *slog.Logger, conf *config.Config, input *struct{}) (response *GetAwsAccountsAllResponse, err error) {
	var (
		service  *awsaccount.Service[*AwsAccount]
		accounts []*AwsAccount
	)
	response = &GetAwsAccountsAllResponse{}

	service, err = Service[*AwsAccount](ctx, log, conf)
	if err != nil {
		err = huma.Error500InternalServerError("failed to connect to service", err)
		return
	}
	defer service.Close()

	accounts, err = service.GetAllAccounts()
	if err != nil {
		err = huma.Error500InternalServerError("failed find all accounts", err)
		return
	}

	response.Body.Data = accounts

	return
}
