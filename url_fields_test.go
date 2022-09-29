package webutils

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldsQuerystring(t *testing.T) {
	tests := []struct {
		name          string
		querystring   url.Values
		requestFields map[string]struct{}
		want          []string
	}{
		{
			name: "empty",
			want: []string{},
		},
		{
			name: "emptyRequestField",
			querystring: url.Values{
				"fields": []string{" "},
			},
			want: []string{},
		},
		{
			name: "emptyRequestFields",
			querystring: url.Values{
				"fields": []string{",,"},
			},
			want: []string{},
		},
		{
			name: "singleValue",
			querystring: url.Values{
				"fields": []string{"id"},
			},
			requestFields: map[string]struct{}{
				"id": {},
			},
			want: []string{"id"},
		},
		{
			name: "notFound",
			querystring: url.Values{
				"fields": []string{"asdf"},
			},
			requestFields: map[string]struct{}{
				"id": {},
			},
			want: []string{},
		},
		{
			name: "multipleValues",
			querystring: url.Values{
				"fields": []string{"name,id"},
			},
			requestFields: map[string]struct{}{
				"id":   {},
				"name": {},
			},
			want: []string{"name", "id"},
		},
		{
			name: "multipleFieldsSingleValue",
			querystring: url.Values{
				"fields": []string{"id", "name"},
			},
			requestFields: map[string]struct{}{
				"id":   {},
				"name": {},
			},
			want: []string{"id", "name"},
		},
		{
			name: "multipleFieldsMultipleValues",
			querystring: url.Values{
				"fields": []string{"id,modified", "name,,created"},
			},
			requestFields: map[string]struct{}{
				"id":       {},
				"name":     {},
				"modified": {},
				"created":  {},
			},
			want: []string{"id", "modified", "name", "created"},
		},
		{
			name: "singleValueDupes",
			querystring: url.Values{
				"fields": []string{"id", "id"},
			},
			requestFields: map[string]struct{}{
				"id": {},
			},
			want: []string{"id"},
		},
		{
			name: "multipleValueDupes",
			querystring: url.Values{
				"fields": []string{"id,id"},
			},
			requestFields: map[string]struct{}{
				"id": {},
			},
			want: []string{"id"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FieldsQuerystring(tt.querystring, tt.requestFields)
			assert.Equal(t, tt.want, got)
		})
	}
}
