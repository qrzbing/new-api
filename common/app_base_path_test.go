package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeAppBasePath(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{name: "empty", value: "", want: ""},
		{name: "root", value: "/", want: ""},
		{name: "single segment", value: "/new-api", want: "/new-api"},
		{name: "trailing slash", value: "/new-api/", want: "/new-api"},
		{name: "multiple segments", value: "/tools/new_api2", want: "/tools/new_api2"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := NormalizeAppBasePath(test.value)
			require.NoError(t, err)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestNormalizeAppBasePathRejectsInvalidValues(t *testing.T) {
	values := []string{
		"new-api",
		"/new-api//nested",
		"/new-api//",
		"/new.api",
		"/new api",
		"/新接口",
		"/new%2Fapi",
		"/new-api?mode=test",
		"/new-api#section",
	}

	for _, value := range values {
		t.Run(value, func(t *testing.T) {
			_, err := NormalizeAppBasePath(value)
			require.Error(t, err)
		})
	}
}

func TestAppPathUsesConfiguredBasePath(t *testing.T) {
	previous := AppBasePath
	AppBasePath = "/tools/new-api"
	t.Cleanup(func() { AppBasePath = previous })

	assert.Equal(t, "/tools/new-api/api/status", AppPath("/api/status"))
	assert.Equal(t, "/tools/new-api/oauth/discord", AppPath("oauth/discord"))
}
