package webutils

import (
	"os"
	"strings"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// ConfigureCORSMiddleware returns an echo.MiddlewareFunc with the whitelisted corsDomains origins
func ConfigureCORSMiddleware(corsDomains []string) echo.MiddlewareFunc {
	// CORS_DOMAINS env var can be one url or a comma-delinated list of urls, ie:
	// http://localhost:8004,https://website.com
	if corsDomains == nil {
		corsDomains = strings.Split(os.Getenv("CORS_DOMAINS"), ",")
	}

	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowOrigins:     corsDomains,
		// which we can explode
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
	})
}
