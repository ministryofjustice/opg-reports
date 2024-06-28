package server

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

func GetUrl(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	apiClient := http.Client{Timeout: time.Second * 3}
	resp, err = apiClient.Do(req)
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
