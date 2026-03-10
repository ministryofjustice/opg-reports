package interfaces

import "opg-reports/report/packages/args"

type ApiLabelGetter interface {
	Label() string
}

type ApiDBGetter interface {
	DB() *args.DB
}

type ApiStatementGetter interface {
	Statement() Statement
}

type ApiConfiguration interface {
	ApiLabelGetter
	ApiDBGetter
	ApiStatementGetter
}
