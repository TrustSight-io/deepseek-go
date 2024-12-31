// Package util provides internal utility functions for the DeepSeek client.
package util

import (
	"bytes"
	"strings"
)

// IsHTML checks if the given bytes appear to be HTML content
func IsHTML(body []byte) bool {
	// Convert to lowercase for case-insensitive matching
	lower := bytes.ToLower(body)

	// Check for common HTML indicators
	return bytes.Contains(lower, []byte("<!doctype html")) ||
		bytes.Contains(lower, []byte("<html")) ||
		bytes.Contains(lower, []byte("</html>")) ||
		bytes.Contains(lower, []byte("<body")) ||
		bytes.Contains(lower, []byte("</body>")) ||
		bytes.Contains(lower, []byte("<head")) ||
		bytes.Contains(lower, []byte("</head>"))
}

// JoinURL joins base URL with path, ensuring proper formatting
func JoinURL(baseURL, path string) string {
	baseURL = strings.TrimRight(baseURL, "/")
	path = strings.TrimLeft(path, "/")
	return baseURL + "/" + path
}
