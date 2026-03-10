package importer

import (
	"opg-reports/report/internal/domains/account/types"
	"opg-reports/report/packages/types/interfaces"
)

var _ interfaces.ImportGetterF[*types.ImportAccount, *types.ImportAccount, *client] = Get
var _ interfaces.ImportFilterF[*types.ImportAccount] = Filter
var _ interfaces.ImportTransformF[*types.ImportAccount, *types.ImportAccount] = Transform
