package webutils

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PaginateQuerystring(t *testing.T) {
	tests := []struct {
		name string
		qs   url.Values
		want PaginatedResults
	}{
		{
			"empty",
			url.Values{},
			PaginatedResults{
				Page:  1,
				Limit: 100,
			},
		},
		{
			"page",
			url.Values{
				"page": {"2", "3", "4"},
			},
			PaginatedResults{
				Page:  2,
				Limit: 100,
			},
		},
		{
			"pageField",
			url.Values{
				"page2": {"2", "3", "4"},
			},
			PaginatedResults{
				Page:  1,
				Limit: 100,
			},
		},
		{
			"pageZero",
			url.Values{
				"page": {"0"},
			},
			PaginatedResults{
				Page:  1,
				Limit: 100,
			},
		},
		{
			"pageInvalid",
			url.Values{
				"page": {"a"},
			},
			PaginatedResults{
				Page:  1,
				Limit: 100,
			},
		},
		{
			"pageZero",
			url.Values{
				"page": {"0"},
			},
			PaginatedResults{
				Page:  1,
				Limit: 100,
			},
		},
		{
			"limit",
			url.Values{
				"limit": {"3", "2", "1"},
			},
			PaginatedResults{
				Page:  1,
				Limit: 3,
			},
		},
		{
			"limitField",
			url.Values{
				"limit2": {"3", "2", "1"},
			},
			PaginatedResults{
				Page:  1,
				Limit: 100,
			},
		},
		{
			"limitMax",
			url.Values{
				"limit": {"1000"},
			},
			PaginatedResults{
				Page:  1,
				Limit: 100,
			},
		},
		{
			"limitZero",
			url.Values{
				"limit": {"0"},
			},
			PaginatedResults{
				Page:  1,
				Limit: 100,
			},
		},
		{
			"limitInvalid",
			url.Values{
				"limit": {"a"},
			},
			PaginatedResults{
				Page:  1,
				Limit: 100,
			},
		},
		{
			"both",
			url.Values{
				"page":  {"4"},
				"limit": {"8"},
			},
			PaginatedResults{
				Page:  4,
				Limit: 8,
			},
		},
		{
			"bothFields",
			url.Values{
				"page2":  {"4"},
				"limit2": {"8"},
			},
			PaginatedResults{
				Page:  1,
				Limit: 100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PaginateQuerystring(tt.qs)
			assert.Equal(t, tt.want, got)
		})
	}
}
