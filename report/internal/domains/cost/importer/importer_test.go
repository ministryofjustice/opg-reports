package importer

import (
	ct "opg-reports/report/internal/domains/cost/types"
	"opg-reports/report/packages/types/interfaces"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
)

var _ interfaces.ImportGetterF[types.ResultByTime, types.ResultByTime, *costexplorer.Client] = Get
var _ interfaces.ImportFilterF[types.ResultByTime] = Filter
var _ interfaces.ImportTransformF[*ct.ImportCost, types.ResultByTime] = Transform
