package webutils

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FilterQuerystring(t *testing.T) {
	tests := []struct {
		name string
		qs   url.Values
		opts []string
		want map[string][]string
	}{
		{
			"empty",
			url.Values{},
			nil,
			make(map[string][]string),
		},
		{
			"fields",
			url.Values{
				"nil":        nil,
				"empty":      []string{},
				"name[cont]": []string{},
				"found":      []string{"a", "b", "c"},
			},
			[]string{
				"nil",
				"empty",
				"found",
				"name[cont]",
				"notfound",
			},
			map[string][]string{
				"nil":        nil,
				"empty":      {},
				"name[cont]": {},
				"found":      {"a", "b", "c"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterQuerystring(tt.qs, tt.opts)
			assert.Equal(t, tt.want, got)
		})
	}
}
