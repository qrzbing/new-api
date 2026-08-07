package common

import (
	"fmt"
	"strings"
)

// NormalizeAppBasePath validates and normalizes the URL prefix used to mount
// the application. An empty string represents root-path deployment.
func NormalizeAppBasePath(value string) (string, error) {
	if value == "" || value == "/" {
		return "", nil
	}
	if !strings.HasPrefix(value, "/") {
		return "", fmt.Errorf("APP_BASE_PATH must start with /")
	}

	normalized := value
	if strings.HasSuffix(normalized, "/") {
		normalized = strings.TrimSuffix(normalized, "/")
	}
	if normalized == "" {
		return "", nil
	}

	for _, segment := range strings.Split(normalized[1:], "/") {
		if segment == "" {
			return "", fmt.Errorf("APP_BASE_PATH must not contain empty path segments")
		}
		for i := 0; i < len(segment); i++ {
			char := segment[i]
			if (char >= 'a' && char <= 'z') ||
				(char >= 'A' && char <= 'Z') ||
				(char >= '0' && char <= '9') ||
				char == '-' || char == '_' {
				continue
			}
			return "", fmt.Errorf("APP_BASE_PATH contains invalid character %q", char)
		}
	}

	return normalized, nil
}

// AppPath adds the configured application base path to an app-owned absolute
// path. External URLs must not be passed to this function.
func AppPath(value string) string {
	if AppBasePath == "" {
		return value
	}
	if value == "" {
		return AppBasePath
	}
	if strings.HasPrefix(value, "/") {
		return AppBasePath + value
	}
	return AppBasePath + "/" + value
}
