package dbstatements

import "opg-reports/report/internal/db/dbmodels"

type DataStatement[T dbmodels.Model, R dbmodels.Result] struct {
	Statement string // Statement is the SQL string with placeholders etc to execute
	Data      T      // Data is the model to use with Statement during execution
	Returned  R      // Returned is the result from the database query
}

type Statement string
