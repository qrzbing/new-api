package router

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
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

func TestWithAppBasePathPreservesPrefixInGinTrailingSlashRedirect(t *testing.T) {
	engine := gin.New()
	engine.POST("/api/channel/", func(context *gin.Context) {
		context.Status(http.StatusNoContent)
	})
	handler := WithAppBasePath(engine, "/new-api")

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/new-api/api/channel", nil))

	assert.Equal(t, http.StatusTemporaryRedirect, response.Code)
	assert.Equal(t, "/new-api/api/channel/", response.Header().Get("Location"))
}

func TestWithAppBasePathRewritesOnlyRootRelativeLocations(t *testing.T) {
	testCases := []struct {
		name     string
		location string
		expected string
	}{
		{name: "root relative", location: "/api/channel/?page=1#result", expected: "/new-api/api/channel/?page=1#result"},
		{name: "absolute", location: "https://example.com/callback", expected: "https://example.com/callback"},
		{name: "protocol relative", location: "//example.com/callback", expected: "//example.com/callback"},
		{name: "relative", location: "sign-in", expected: "sign-in"},
		{name: "base path", location: "/new-api", expected: "/new-api"},
		{name: "already prefixed", location: "/new-api/dashboard", expected: "/new-api/dashboard"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			handler := WithAppBasePath(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				writer.Header().Set("Location", testCase.location)
				writer.WriteHeader(http.StatusCreated)
			}), "/new-api")

			response := httptest.NewRecorder()
			handler.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/new-api/api/resource", nil))

			assert.Equal(t, http.StatusCreated, response.Code)
			assert.Equal(t, testCase.expected, response.Header().Get("Location"))
		})
	}
}

func TestWithAppBasePathRewritesLocationWithoutExplicitStatus(t *testing.T) {
	handler := WithAppBasePath(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Location", "/sign-in")
	}), "/new-api")

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/new-api/dashboard", nil))

	assert.Equal(t, "/new-api/sign-in", response.Header().Get("Location"))
}
