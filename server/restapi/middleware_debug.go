//go:build debug

package restapi

import (
	"net/http"
	"strings"
)

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return cors(handler)
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}

func cors(h http.Handler) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		origin := req.Header.Get("Origin")
		if strings.HasPrefix(origin, "http://localhost") {
			res.Header().Set("Access-Control-Allow-Origin", origin)
		}
		h.ServeHTTP(res, req)
	}
}
