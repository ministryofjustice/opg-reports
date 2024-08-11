package urls

import (
	"net/http"
	"net/url"
	"time"
)

func Get(url *url.URL) (resp *http.Response, err error) {
	u := url.String()
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return
	}
	apiClient := http.Client{Timeout: time.Second * 3}
	resp, err = apiClient.Do(req)
	return
}
