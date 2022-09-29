package webutils

import (
	"context"
	"strings"

	"github.com/google/uuid"
	echo "github.com/labstack/echo/v4"
)

type ctxKey int

const ProvenanceIDHeader = "x-provenance-id"
const RequestIDHeader = "x-request-id"

const pidKey ctxKey = ctxKey(4200)
const ridKey ctxKey = ctxKey(4201)

// NewContext creates a context with provenance id
func NewContext(ctx context.Context, provenanceID string, requestID string) context.Context {
	ctx = context.WithValue(ctx, pidKey, provenanceID)
	ctx = context.WithValue(ctx, ridKey, requestID)

	return ctx
}

// ProvenanceIDFromContext returns the provenance id from context
func ProvenanceIDFromContext(ctx context.Context) (string, bool) {
	pid, ok := ctx.Value(pidKey).(string)
	return pid, ok
}

// RequestIDFromContext returns the provenance id from context
func RequestIDFromContext(ctx context.Context) (string, bool) {
	rid, ok := ctx.Value(ridKey).(string)
	return rid, ok
}

func ProvenanceIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		provenanceID := c.Request().Header.Get(ProvenanceIDHeader)
		if provenanceID == "" {
			provenanceID = strings.ReplaceAll(uuid.New().String(), "-", "")
			c.Request().Header.Set(ProvenanceIDHeader, provenanceID)
		}

		requestID := strings.ReplaceAll(uuid.New().String(), "-", "")

		ctx := NewContext(c.Request().Context(), provenanceID, requestID)

		c.Response().Header().Add(ProvenanceIDHeader, provenanceID)

		c.Response().Header().Add(RequestIDHeader, requestID)

		r := c.Request().Clone(ctx)

		c.SetRequest(r)

		return next(c)
	}
}
