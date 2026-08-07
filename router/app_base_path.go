package router

import (
	"net/http"
	"strings"
)

// WithAppBasePath mounts an existing root-based handler under basePath.
func WithAppBasePath(handler http.Handler, basePath string) http.Handler {
	if basePath == "" {
		return handler
	}

	strippedHandler := http.StripPrefix(basePath, handler)
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/" || request.URL.Path == basePath {
			target := basePath + "/"
			if request.URL.RawQuery != "" {
				target += "?" + request.URL.RawQuery
			}
			http.Redirect(writer, request, target, http.StatusPermanentRedirect)
			return
		}
		if !strings.HasPrefix(request.URL.Path, basePath+"/") {
			http.NotFound(writer, request)
			return
		}
		forwardedRequest := request.Clone(request.Context())
		forwardedRequest.RequestURI = strings.TrimPrefix(request.RequestURI, basePath)
		strippedHandler.ServeHTTP(writer, forwardedRequest)
	})
}
