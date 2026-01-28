package uptime

import "errors"

var (
	ErrIncorrectRegion          = errors.New("metrics must be fetched via us-east-1.")
	ErrFailedGettingMetricsList = errors.New("failed to get metrics list with error.")
	ErrFailedGettingMetricStats = errors.New("failed to get metric statistics with error.")
)
