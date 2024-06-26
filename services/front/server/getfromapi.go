package server

import (
	"net/http"
	"net/url"
	"opg-reports/shared/server/response"
	"strings"
	"time"
)

// *response.TableData[*response.Cell, *response.Row[*response.Cell]]
func GetFromApi(url string) (resp *response.Result[*response.Cell, *response.Row[*response.Cell], *response.TableData[*response.Cell, *response.Row[*response.Cell]]], err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	apiClient := http.Client{Timeout: time.Second * 3}
	apiResp, err := apiClient.Do(req)
	if err != nil {
		return
	}
	resp = response.NewResponse()
	response.NewResponseFromHttp(apiResp, resp)
	return
}

func Url(scheme string, host string, path string) *url.URL {
	if scheme == "" {
		scheme = "http"
	}
	if host != "" && host[0:1] == ":" {
		host = "localhost" + host
	}
	host = strings.ReplaceAll(host, "http://", "")
	host = strings.ReplaceAll(host, "https://", "")

	if path[len(path)-1:] != "/" {
		path = path + "/"
	}
	path = strings.ReplaceAll(path, "http://", "")
	path = strings.ReplaceAll(path, "https://", "")

	return &url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}
}
