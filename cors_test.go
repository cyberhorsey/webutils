package webutils

import (
	"os"
	"testing"
)

func Test_ConfigureCORSMiddleware(t *testing.T) {
	tests := []struct {
		name        string
		corsDomains []string
		envVarFunc  func()
	}{
		{
			"noCorsDomains",
			nil,
			nil,
		},
		{
			"corsDomain",
			[]string{"https://localhost:8004", "https://localhost:8003"},
			nil,
		},

		{
			"corsDomainsViaEnv",
			nil,
			func() {
				os.Setenv("CORS_DOMAINS", "https://localhost:8004,https://localhost:8003")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envVarFunc != nil {
				tt.envVarFunc()
			}
			_ = ConfigureCORSMiddleware(tt.corsDomains)
		})
	}
}
