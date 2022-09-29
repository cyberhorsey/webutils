package webutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMappedFields(t *testing.T) {
	tests := []struct {
		name          string
		requestFields []string
		columnMap     map[string]string
		want          []string
	}{
		{
			name: "emptyBoth",
			want: []string{},
		},
		{
			name: "emptyFields",
			columnMap: map[string]string{
				"abc": "xyz",
			},
			want: []string{},
		},
		{
			name:          "notFound",
			requestFields: []string{"123"},
			columnMap: map[string]string{
				"abc": "xyz",
			},
			want: []string{},
		},
		{
			name:          "found",
			requestFields: []string{"abc"},
			columnMap: map[string]string{
				"abc": "xyz",
			},
			want: []string{"xyz"},
		},
		{
			name:          "foundDupes",
			requestFields: []string{"abc", "abc"},
			columnMap: map[string]string{
				"abc": "xyz",
			},
			want: []string{"xyz"},
		},
		{
			name:          "foundDupeColumns",
			requestFields: []string{"abc", "xyz"},
			columnMap: map[string]string{
				"abc": "123",
				"xyz": "123",
			},
			want: []string{"123"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MappedFields(tt.requestFields, tt.columnMap)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestContainsReturnField(t *testing.T) {
	tests := []struct {
		name    string
		fields  []string
		matches []string
		want    bool
	}{
		{
			name:    "noFields",
			fields:  []string{},
			matches: []string{"category"},
			want:    true,
		},
		{
			name:    "exists",
			fields:  []string{"category"},
			matches: []string{"category"},
			want:    true,
		},
		{
			name:    "doesntExist",
			fields:  []string{"supplier"},
			matches: []string{"category"},
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IncludesAnyReturnField(tt.fields, tt.matches...)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestReturnFieldsWith(t *testing.T) {
	tests := []struct {
		name          string
		fields        []string
		fieldToAppend string
		wantFields    []string
	}{
		{
			name:          "fieldsWithfieldNotExistent",
			fields:        []string{"item"},
			fieldToAppend: "id",
			wantFields:    []string{"item", "id"},
		},
		{
			name:          "fieldsWithFieldExistent",
			fields:        []string{"id"},
			fieldToAppend: "id",
			wantFields:    []string{"id"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReturnFieldsWith(tt.fields, tt.fieldToAppend)
			assert.Equal(t, tt.wantFields, got)
		})
	}
}
