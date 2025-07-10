package awsaccounts

import (
	"context"
	"log/slog"
	"net/http"

	"opg-reports/report/config"
	"opg-reports/report/endpoints"
	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/service/api"

	"github.com/danielgtaylor/huma/v2"
)

// GetAwsAccountsAllResponse
type GetAwsAccountsAllResponse[T api.Model] struct {
	Body struct {
		Count int `json:"count,omitempty"`
		Data  []T `json:"data"`
	}
}

// RegisterGetAwsAccountsAll handles registering the endpoint to return all aws accounts stored
// in the data.
//
// Attaches both the handler (`handleGetAwsAccountsAll`) and the spec details to the huma api
func RegisterGetAwsAccountsAll[T api.Model](log *slog.Logger, conf *config.Config, humaapi huma.API, service api.AwsAccountsGetter[T], store sqlr.Reader) {
	var operation = huma.Operation{
		OperationID:   "get-awsaccounts-all",
		Method:        http.MethodGet,
		Path:          endpoints.AWSACCOUNTS_GET_ALL,
		Summary:       "All AWS accounts",
		Description:   "Returns a list of all AWS accounts stored.",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"AWS Accounts"},
	}
	huma.Register(humaapi, operation, func(ctx context.Context, input *struct{}) (*GetAwsAccountsAllResponse[T], error) {
		return handleGetAwsAccountsAll(ctx, log, conf, service, store, input)
	})
}

// handleGetAwsAccountsAll deals with each request and fetches
func handleGetAwsAccountsAll[T api.Model](ctx context.Context, log *slog.Logger, conf *config.Config,
	service api.AwsAccountsGetter[T], store sqlr.Reader, input *struct{}) (response *GetAwsAccountsAllResponse[T], err error) {
	var accounts []T
	response = &GetAwsAccountsAllResponse[T]{}

	if service == nil {
		err = huma.Error500InternalServerError("failed to connect to service", err)
		return
	}
	defer service.Close()

	accounts, err = service.GetAllAwsAccounts(store)
	if err != nil {
		err = huma.Error500InternalServerError("failed find all accounts", err)
		return
	}

	response.Body.Data = accounts
	response.Body.Count = len(accounts)

	return
}
