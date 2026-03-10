package api

import (
	"opg-reports/report/internal/domains/home/types"
	"opg-reports/report/packages/types/interfaces"
)

var _ interfaces.Row = &types.NilRow{}
