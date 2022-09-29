package webutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SanitizeString(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected string
	}{
		{
			"html",
			"<html>string</html>",
			"string",
		},
		{
			"js",
			`Hello <STYLE>.XSS{background-image:url("javascript:alert('XSS')");}</STYLE><A CLASS=XSS></A>World`,
			"Hello World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, SanitizeString(tt.s))
		})
	}
}

func Test_SanitizeHTML(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected string
	}{
		{
			"html",
			"<section>string</section>",
			"<section>string</section>",
		},
		{
			"xssInjection",
			`Hello <STYLE>.XSS{background-image:url("javascript:alert('XSS')");}</STYLE><A CLASS=XSS></A>World`,
			"Hello World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, SanitizeHTML(tt.s))
		})
	}
}
