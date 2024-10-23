// Package endpoiints is a list of all endpoints running on the api
package endpoints

import "github.com/ministryofjustice/opg-reports/costs/costsapi"

type ApiEndpoint string

// Endpoints for the costsapi
// - Used in navigation configuration
const (
	CostsAwsTotal              ApiEndpoint = ApiEndpoint(costsapi.EndpointAwsTotal)
	CostsAwsTaxOverview        ApiEndpoint = ApiEndpoint(costsapi.EndpointAwsTaxOverview)
	CostsAwsPerUnit            ApiEndpoint = ApiEndpoint(costsapi.EndpointAwsPerUnit)
	CostsAwsPerUnitEnvironment ApiEndpoint = ApiEndpoint(costsapi.EndpointAwsPerUnitEnvironment)
	CostsAwsDetailed           ApiEndpoint = ApiEndpoint(costsapi.EndpointAwsDetailed)
)
