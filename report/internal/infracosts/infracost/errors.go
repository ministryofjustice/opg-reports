package infracost

import "errors"

// ErrGettingCostData is used when the an error is returned from the sdk GetCostAndUsage
var ErrGettingCostData error = errors.New("call to GetCostAndUsage failed with an error.")
