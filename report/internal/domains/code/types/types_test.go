package types

import "opg-reports/report/packages/types/interfaces"

var (
	_ interfaces.Insertable = &Codebase{}
	_ interfaces.Selectable = &Codebase{}
	_ interfaces.Statement  = &Select{}
)
