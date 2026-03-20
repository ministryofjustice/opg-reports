package api

import (
	"context"
	"opg-reports/report/packages/httpx"
	"opg-reports/report/packages/slogx"
)

const label string = `get-ping`

// Ping returns a virtually empty endpoint for simple connection confirmation
func Ping[T httpx.ResponseData](ctx context.Context, cfg httpx.MuxConfigurer, r httpx.FitleredRequest, response *httpx.ResponseData) {
	var log = slogx.FromContext(ctx)
	log.Info(ctx, "ping recieved.")
}
