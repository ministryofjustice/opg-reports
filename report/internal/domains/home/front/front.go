package front

import (
	"context"
	"opg-reports/report/packages/httpx"
	"opg-reports/report/packages/slogx"
)

const label string = `front-home`

// SetPageBaseline configures the baseline values for this page;
// generally this is the page name & title
func SetPageBaseline[T *Page](ctx context.Context, cfg httpx.MuxConfigurer, r httpx.FitleredRequest, response *Page) {
	var filter = r.Filter()
	var log = slogx.FromContext(ctx)
	log.Info(ctx, "setting page baseline", "label", label)

	// rebuild the response with all the info we have
	response = &Page{
		HTMLResponseData: httpx.HTMLResponseData{
			ResponseData: httpx.ResponseData{
				Version:      response.Version,
				Request:      response.Request,
				GovUKVersion: response.GovUKVersion,
				Teams:        []string{},
			},
			Name:  `OPG Reports`,
			Title: `OPG Reports`,
		},
	}

	if filter.Team != "" {
		response.Title = filter.Team + ` - ` + response.Title
	}
	log.Info(ctx, "page baseline complete.", "label", label)
}
