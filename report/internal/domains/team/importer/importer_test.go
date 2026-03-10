package importer

import (
	"opg-reports/report/internal/domains/team/types"
	"opg-reports/report/packages/types/interfaces"
)

var _ interfaces.ImportGetterF[*types.Team, *types.Team, *client] = Get
var _ interfaces.ImportFilterF[*types.Team] = Filter
var _ interfaces.ImportTransformF[*types.Team, *types.Team] = Transform
