package types

import "opg-reports/report/packages/types/interfaces"

var (
	_ interfaces.Insertable = &ImportAccount{}
	_ interfaces.Selectable = &Account{}
	_ interfaces.Resultable = &Account{}
	_ interfaces.Statement  = &Select{}
)
