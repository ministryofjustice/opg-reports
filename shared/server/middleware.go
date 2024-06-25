package server

import (
	"fmt"
	"log/slog"
	"net/http"
)

// MiddlewareHandler is alias for http.HandlerFunc to allow chaining
type MiddlewareHandlerFunc func(http.HandlerFunc) http.HandlerFunc

// MiddlewareFunc type matching for chaining middleware
type MiddlewareFunc func(w http.ResponseWriter, r *http.Request)

// Chain facilitates creation of middleware chains, with them being called in sequence
// Normally called by the MiddlewareServerMux Register func, but can be called directly
// if required
func Middleware(handlerFunc http.HandlerFunc, middleware ...MiddlewareFunc) http.HandlerFunc {

	for _, mw := range middleware {
		handlerFunc = wrap(mw)(handlerFunc)
	}
	return handlerFunc
}

// wrap is internal facing and used to create the correct wrapping aruond
// of function chains for middleware chaining so our handler methods
// dont need to deal with nested returns etc
func wrap(middewareFunc MiddlewareFunc) MiddlewareHandlerFunc {

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// call this function
			middewareFunc(w, r)
			// call next in middleware stack
			next(w, r)

		}
	}
}

var securityHeaders = map[string]string{
	"Content-Security-Policy":   "default-src 'self'",
	"Referrer-Policy":           "same-origin",
	"Strict-Transport-Security": "max-age=31536000; includeSubDomains; preload",
	"X-Content-Type-Options":    "nosniff",
	"X-Frame-Options":           "SAMEORIGIN",
	"X-Xss-Protection":          "1; mode=block",
	"Cache-Control":             "no-store",
	"Pragma":                    "no-cache",
}

// SecurityHeadersMW attaches standard headers to the response
func SecurityHeadersMW(w http.ResponseWriter, r *http.Request) {
	for header, value := range securityHeaders {
		w.Header().Add(header, value)
	}
}

func LoggingMW(w http.ResponseWriter, r *http.Request) {
	slog.Info(
		fmt.Sprintf(""),
		slog.String("request_method", r.Method),
		slog.String("request_uri", r.URL.String()),
	)
}