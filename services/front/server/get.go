package server

import (
	"fmt"
	"log/slog"
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

func Url(scheme string, host string, path string) (u *url.URL) {
	slog.Debug("generating url",
		slog.String("scheme", scheme),
		slog.String("host", host),
		slog.String("path", path))

	// set defaults
	if scheme == "" {
		scheme = "http"
	}
	scheme = strings.ReplaceAll(scheme, "://", "")

	parsedHost := host
	if host == "" {
		parsedHost = "localhost"
	} else if host[0:1] == ":" {
		parsedHost = "localhost" + parsedHost
	}

	parsedHost = strings.ReplaceAll(parsedHost, "https://", "")
	parsedHost = strings.ReplaceAll(parsedHost, "http://", "")

	path = strings.ReplaceAll(path, "https://", "")
	path = strings.ReplaceAll(path, "http://", "")
	path = strings.ReplaceAll(path, host, "")
	path = strings.ReplaceAll(path, parsedHost, "")

	raw := fmt.Sprintf("%s://%s%s", scheme, parsedHost, path)
	u, err := url.Parse(raw)

	// add trialing slash to the end of the path
	p := u.Path
	last := p[len(p)-1:]
	if err == nil && last != "/" {
		u.Path = p + "/"
	}

	slog.Debug("generated url",
		slog.String("scheme", scheme),
		slog.String("host", host),
		slog.String("parsedHost", parsedHost),
		slog.String("path", path),
		slog.String("raw", raw),
		slog.String("u", u.String()))

	if err != nil {
		return nil
	}
	return

}
