// Package `handler` contains api and fron end handlers and parsing.
package handler

import (
	"opg-reports/report/packages/args"
	"opg-reports/report/packages/types/interfaces"
)

// APIConfig
type ApiConfig struct {
	Name     string
	Database *args.DB
	Selector interfaces.Statement
}

func (self *ApiConfig) Label() string {
	return self.Name
}
func (self *ApiConfig) DB() *args.DB {
	return self.Database
}
func (self *ApiConfig) Statement() interfaces.Statement {
	return self.Selector
}
