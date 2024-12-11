/*
convertor
*/
package main

import (
	"log/slog"

	"github.com/ministryofjustice/opg-reports/convertor/lib"
)

var args = &lib.Arguments{}

func main() {
	var err error
	slog.Info("[convertor] starting ")
	lib.SetupArgs(args)
	err = lib.Run(args)
	if err != nil {
		panic(err)
	}
	slog.Info("[convertor] done     ✅")
}
