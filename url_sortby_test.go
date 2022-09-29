package webutils

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SortByQuerystring(t *testing.T) {
	tests := []struct {
		name string
		qs   url.Values
		opts []string
		want []string
	}{
		{
			"empty",
			url.Values{},
			nil,
			[]string{},
		},
		{
			"emptySortBy",
			url.Values{
				"sort_by": nil,
			},
			[]string{},
			[]string{},
		},
		{
			"found",
			url.Values{
				"sort_by": []string{
					"created",
				},
			},
			[]string{
				"created",
			},
			[]string{
				"created",
			},
		},
		{
			"foundTwice",
			url.Values{
				"sort_by": []string{
					"created",
					"created",
				},
			},
			[]string{
				"created",
			},
			[]string{
				"created",
			},
		},
		{
			"foundDupeAscDesc",
			url.Values{
				"sort_by": []string{
					"created",
					"-created",
				},
			},
			[]string{
				"created",
				"-created",
			},
			[]string{
				"created",
			},
		},
		{
			"foundDupeDescAsc",
			url.Values{
				"sort_by": []string{
					"-created",
					"created",
				},
			},
			[]string{
				"created",
				"-created",
			},
			[]string{
				"-created",
			},
		},
		{
			"foundBoth",
			url.Values{
				"sort_by": []string{
					"created",
					"-updated",
				},
			},
			[]string{
				"created",
				"-updated",
			},
			[]string{
				"created",
				"-updated",
			},
		},
		{
			"foundBoth2",
			url.Values{
				"sort_by": []string{
					"created",
					"updated",
				},
			},
			[]string{
				"updated",
				"created",
			},
			[]string{
				"created",
				"updated",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SortByQuerystring(tt.qs, tt.opts)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_parseSortValues(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{
			name:  "empty",
			input: nil,
			want:  []string{},
		},
		{
			name:  "single",
			input: []string{"a"},
			want:  []string{"a"},
		},
		{
			name:  "singleDash",
			input: []string{"-"},
			want:  []string{"-"},
		},
		{
			name:  "dupeSingleAscAsc",
			input: []string{"a,a"},
			want:  []string{"a"},
		},
		{
			name:  "dupeSingleAscDesc",
			input: []string{"a,-a"},
			want:  []string{"a"},
		},
		{
			name:  "dupeSingleDescAsc",
			input: []string{"-a,a"},
			want:  []string{"-a"},
		},
		{
			name:  "multiple",
			input: []string{"a,,c", "d"},
			want:  []string{"a", "c", "d"},
		},
		{
			name:  "dupeMultipleAscAsc",
			input: []string{"a", "a"},
			want:  []string{"a"},
		},
		{
			name:  "dupeMultipleAscDesc",
			input: []string{"a", "-a"},
			want:  []string{"a"},
		},
		{
			name:  "dupeMultipleDescAsc",
			input: []string{"-a", "a"},
			want:  []string{"-a"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseSortValues(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_parseSortBy(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "empty",
			input: "",
			want:  []string{},
		},
		{
			name:  "single",
			input: "a",
			want:  []string{"a"},
		},
		{
			name:  "multiple",
			input: "a,,c",
			want:  []string{"a", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseSortBy(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
