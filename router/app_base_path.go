package router

import (
	"bufio"
	"net"
	"net/http"
	"strings"
)

type appBasePathResponseWriter struct {
	http.ResponseWriter
	basePath string
}

func (writer *appBasePathResponseWriter) rewriteLocation() {
	location := writer.Header().Get("Location")
	if location == "" || !strings.HasPrefix(location, "/") || strings.HasPrefix(location, "//") {
		return
	}
	if location == writer.basePath ||
		strings.HasPrefix(location, writer.basePath+"/") ||
		strings.HasPrefix(location, writer.basePath+"?") ||
		strings.HasPrefix(location, writer.basePath+"#") {
		return
	}
	writer.Header().Set("Location", writer.basePath+location)
}

func (writer *appBasePathResponseWriter) WriteHeader(statusCode int) {
	writer.rewriteLocation()
	writer.ResponseWriter.WriteHeader(statusCode)
}

func (writer *appBasePathResponseWriter) Write(data []byte) (int, error) {
	writer.rewriteLocation()
	return writer.ResponseWriter.Write(data)
}

func (writer *appBasePathResponseWriter) Flush() {
	writer.rewriteLocation()
	writer.ResponseWriter.(http.Flusher).Flush()
}

func (writer *appBasePathResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return writer.ResponseWriter.(http.Hijacker).Hijack()
}

func (writer *appBasePathResponseWriter) CloseNotify() <-chan bool {
	return writer.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func (writer *appBasePathResponseWriter) Unwrap() http.ResponseWriter {
	return writer.ResponseWriter
}

type appBasePathPushResponseWriter struct {
	*appBasePathResponseWriter
}

func (writer *appBasePathPushResponseWriter) Push(target string, options *http.PushOptions) error {
	return writer.ResponseWriter.(http.Pusher).Push(target, options)
}

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
		responseWriter := &appBasePathResponseWriter{ResponseWriter: writer, basePath: basePath}
		defer responseWriter.rewriteLocation()
		if _, ok := writer.(http.Pusher); ok {
			strippedHandler.ServeHTTP(&appBasePathPushResponseWriter{responseWriter}, forwardedRequest)
			return
		}
		strippedHandler.ServeHTTP(responseWriter, forwardedRequest)
	})
}
