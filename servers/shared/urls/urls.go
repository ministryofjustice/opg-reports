package urls

import (
	"fmt"
	"log/slog"
	"net/url"
	"strings"
)

// clean tries to clean up and set correct values for each
// scheme, host and path
// replaces a lot of common errors (like host in the path etc)
func clean(scheme string, host string, path string) string {
	slog.Debug("cleaning url",
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
	path = strings.ReplaceAll(path, "localhost", "")
	path = strings.ReplaceAll(path, host, "")
	path = strings.ReplaceAll(path, parsedHost, "")

	slog.Debug("cleaned url",
		slog.String("scheme", scheme),
		slog.String("host", host),
		slog.String("parsedHost", parsedHost),
		slog.String("path", path))

	return fmt.Sprintf("%s://%s%s", scheme, parsedHost, path)
}

// Parse takes a mix of scheme, host and path strings and generates
// a url object from those
// Generally used to generate the url to call for the api from the config data
func Parse(scheme string, host string, path string) (u *url.URL) {

	var raw string
	// if path contains the scheme, the use the path directly
	if strings.HasPrefix(path, "http") && strings.Contains(path, "://") {
		raw = path
	} else {
		raw = clean(scheme, host, path)
	}
	u, err := url.Parse(raw)

	// add trialing slash to the end of the path
	p := u.Path
	if len(p) > 0 {
		last := p[len(p)-1:]
		if err == nil && last != "/" {
			u.Path = p + "/"
		}
	}

	slog.Debug("generated url",
		slog.String("url", u.String()),
		slog.String("err", fmt.Sprintf("%+v", err)))

	if err != nil {
		return nil
	}
	return

}
