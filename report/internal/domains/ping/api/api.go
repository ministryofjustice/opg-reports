package api

import (
	"context"
	"opg-reports/report/packages/httpx"
	"opg-reports/report/packages/slogx"
)

const label string = `get-ping`

// Ping returns a virtually empty endpoint for simple connection confirmation
func Ping(ctx context.Context, m httpx.Mux, r httpx.FitleredRequest, cfg httpx.MuxConfig, response *httpx.ResponseContent) {
	var log = slogx.FromContext(ctx)
	log.Info(ctx, "ping recieved.")
}
