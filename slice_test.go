package webutils

import "testing"

func TestContainsString(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		contains string
		want     bool
	}{
		{
			name: "empty",
		},
		{
			name:     "notFound",
			slice:    []string{"abc"},
			contains: "xyz",
		},
		{
			name:     "found",
			slice:    []string{"abc"},
			contains: "abc",
			want:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsString(tt.slice, tt.contains); got != tt.want {
				t.Errorf("ContainsString() = %v, want %v", got, tt.want)
			}
		})
	}
}
