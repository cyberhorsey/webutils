package webutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalFields(t *testing.T) {
	var u64 uint64 = 0xFFFFFFFFFFFFFFFF

	input := map[string]interface{}{
		"a":      true,
		"b":      -2,
		"c":      "3",
		"d":      struct{}{},
		"float":  -1.23456789,
		"uint64": u64,
	}

	tests := []struct {
		name    string
		v       interface{}
		fields  []string
		want    string
		wantErr error
	}{
		{
			name:    "nil",
			v:       nil,
			fields:  nil,
			want:    "null",
			wantErr: nil,
		},
		{
			name:    "true",
			v:       true,
			fields:  nil,
			want:    "true",
			wantErr: nil,
		},
		{
			name:    "false",
			v:       false,
			fields:  nil,
			want:    "false",
			wantErr: nil,
		},
		{
			name:    "string",
			v:       "hello",
			fields:  nil,
			want:    `"hello"`,
			wantErr: nil,
		},
		{
			name:    "slice",
			v:       []int{1},
			fields:  nil,
			want:    "[1]",
			wantErr: nil,
		},
		{
			name:    "int",
			v:       -1,
			fields:  nil,
			want:    "-1",
			wantErr: nil,
		},
		{
			name:    "uint",
			v:       1,
			fields:  nil,
			want:    "1",
			wantErr: nil,
		},
		{
			name:    "uint64",
			v:       u64,
			fields:  nil,
			want:    fmt.Sprintf("%d", u64),
			wantErr: nil,
		},
		{
			name:    "float",
			v:       -1.23456789,
			fields:  nil,
			want:    "-1.23456789",
			wantErr: nil,
		},
		{
			name:    "struct",
			v:       struct{}{},
			fields:  nil,
			want:    "{}",
			wantErr: nil,
		},
		{
			name:    "noFields",
			v:       input,
			fields:  nil,
			want:    "{}",
			wantErr: nil,
		},
		{
			name:    "noFieldsCaseSensitive",
			v:       input,
			fields:  []string{"A", "B", "C", "D"},
			want:    "{}",
			wantErr: nil,
		},
		{
			name:    "allFields",
			v:       input,
			fields:  []string{"a", "b", "c", "d", "float", "uint64"},
			want:    `{"a":true,"b":-2,"c":"3","d":{},"float":-1.23456789,"uint64":18446744073709551615}`,
			wantErr: nil,
		},
		{
			name:    "oneField",
			v:       input,
			fields:  []string{"a"},
			want:    `{"a":true}`,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResponseFields(tt.v, tt.fields)
			assert.True(t, errors.Is(err, tt.wantErr))

			if tt.wantErr == nil {
				bs, err := json.Marshal(got)
				assert.Nil(t, err)
				assert.Equal(t, tt.want, string(bs))
			}
		})
	}
}
