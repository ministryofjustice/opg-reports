package types

import "opg-reports/report/packages/types/interfaces"

var (
	_ interfaces.Insertable = &ImportCost{}
	_ interfaces.Selectable = &CostByTeam{}
	_ interfaces.Statement  = &Select{}
)
