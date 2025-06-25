package awsaccounts

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/awsaccount"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

func RegisterGetAwsAccountsAll(log *slog.Logger, conf *config.Config, api huma.API) {
	var operation = huma.Operation{
		OperationID:   "get-awsaccounts-all",
		Method:        http.MethodGet,
		Path:          "/v1/awsaccounts/all",
		Summary:       "Return all AWS accounts",
		Description:   "Returns a list of all AWS accounts known about.",
		DefaultStatus: http.StatusOK,
	}
	huma.Register(api, operation, func(ctx context.Context, input *struct{}) (*GetAwsAccountsAllResponse, error) {
		return handleGetAwsAccountsAll(ctx, log, conf, input)
	})
}

// handleGetAwsAccountsAll deals with each request and fetches
func handleGetAwsAccountsAll(ctx context.Context, log *slog.Logger, conf *config.Config, input *struct{}) (response *GetAwsAccountsAllResponse, err error) {
	var (
		service      *awsaccount.Service[*awsaccount.AwsAccount]
		accounts     []*awsaccount.AwsAccount
		responseData []*AwsAccount = []*AwsAccount{}
	)
	response = &GetAwsAccountsAllResponse{}

	service, err = Service(ctx, log, conf)
	if err != nil {
		err = huma.Error500InternalServerError("failed to connect to service", err)
		return
	}

	accounts, err = service.GetAllAccounts()
	if err != nil {
		err = huma.Error500InternalServerError("failed find all accounts", err)
		return
	}
	// convert database data to response version - stripping any extra fields etc
	err = utils.Convert(accounts, &responseData)
	if err != nil {
		err = huma.Error500InternalServerError("failed to convert data to accounts", err)
		return
	}

	response.Body.Data = responseData

	return
}
