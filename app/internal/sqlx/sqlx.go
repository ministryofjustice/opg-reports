// Package sqlx contains helpful additional functions to make
// working with the std lib sql easier.
//
// In particular, thie package provides a `Bind` function that
// allows combining of data from a struct into a SQL string
// by utilised field placeholders (`:fieldName`) and returns
// a correctly parameterised string and ordered arguments
// that are compatible with `db.ExecContext` or `db.QueryContext`
package sqlx
