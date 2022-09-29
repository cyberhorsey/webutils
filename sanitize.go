package webutils

import (
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

// SanitizeHTML removes xss attack vectors, leaving
// you with the html and text by default.
// Opts can be passed in to configure how sanitization should occur.
func SanitizeHTML(s string) string {
	p := bluemonday.UGCPolicy()

	return strings.TrimSpace(p.Sanitize(s))
}

// SanitizeString removes html, css, JS, and all other vulnerabilities, leaving
// you with the text by default.
// Opts can be passed in to configure how sanitization should occur.
func SanitizeString(s string) string {
	p := bluemonday.NewPolicy()

	return strings.TrimSpace(p.Sanitize(s))
}
