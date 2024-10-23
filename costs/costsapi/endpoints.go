package costsapi

const (
	EndpointAwsTotal              string = "/{version}/costs/aws/total/{start_date}/{end_date}"
	EndpointAwsTaxOverview        string = "/{version}/costs/aws/tax-overview/{start_date}/{end_date}/{interval}?unit={unit}"
	EndpointAwsPerUnit            string = "/{version}/costs/aws/unit/{start_date}/{end_date}/{interval}?unit={unit}"
	EndpointAwsPerUnitEnvironment string = "/{version}/costs/aws/unit-environment/{start_date}/{end_date}/{interval}?unit={unit}"
	EndpointAwsDetailed           string = "/{version}/costs/aws/unit-environment/{start_date}/{end_date}/{interval}?unit={unit}"
)
