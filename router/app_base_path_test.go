package router

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithAppBasePathPassesThroughRootDeployment(t *testing.T) {
	handler := WithAppBasePath(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusNoContent)
	}), "")

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/status", nil))

	assert.Equal(t, http.StatusNoContent, response.Code)
}

func TestWithAppBasePathRedirectsCanonicalEntryPoints(t *testing.T) {
	handler := WithAppBasePath(http.NotFoundHandler(), "/new-api")
	for _, requestPath := range []string{"/?aff=abc", "/new-api?aff=abc"} {
		t.Run(requestPath, func(t *testing.T) {
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, requestPath, nil))

			assert.Equal(t, http.StatusPermanentRedirect, response.Code)
			assert.Equal(t, "/new-api/?aff=abc", response.Header().Get("Location"))
		})
	}
}

func TestWithAppBasePathStripsPrefixBeforeRouting(t *testing.T) {
	handler := WithAppBasePath(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, err := fmt.Fprintf(writer, "%s|%s|%s", request.URL.Path, request.URL.RawQuery, request.RequestURI)
		require.NoError(t, err)
	}), "/tools/new-api")

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/tools/new-api/api/status?check=1", nil))

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "/api/status|check=1|/api/status?check=1", response.Body.String())
}

func TestWithAppBasePathRejectsUnprefixedAndNearMatchPaths(t *testing.T) {
	handler := WithAppBasePath(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusNoContent)
	}), "/new-api")

	for _, requestPath := range []string{"/api/status", "/v1/models", "/new-api-other/api/status"} {
		t.Run(requestPath, func(t *testing.T) {
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, requestPath, nil))

			assert.Equal(t, http.StatusNotFound, response.Code)
		})
	}
}
