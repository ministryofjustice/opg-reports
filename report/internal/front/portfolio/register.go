package portfolio

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/global/frontmodels"
	"opg-reports/report/package/cntxt"
)

const ENDPOINT string = "/home/portfolio"

func Register(ctx context.Context, mux *http.ServeMux, args *frontmodels.FrontRegisterArgs) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "portfolio", "func", "Register")

	log.Info("registering handler [`" + ENDPOINT + "/{$}`] ...")
	mux.HandleFunc(fmt.Sprintf("%s/{$}", ENDPOINT), func(writer http.ResponseWriter, request *http.Request) {
		Handler(ctx, args, writer, request)
	})
}
