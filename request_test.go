package webutils

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_NewContext(t *testing.T) {
	ctx := context.Background()
	ctx = NewContext(ctx, "pid", "rid")
	val, ok := ProvenanceIDFromContext(ctx)
	assert.NotEqual(t, false, ok)
	assert.Equal(t, val, "pid")
	val, ok = RequestIDFromContext(ctx)
	assert.NotEqual(t, false, ok)
	assert.Equal(t, val, "rid")
}

func Test_ProvenanceIDMiddleware(t *testing.T) {
	e := echo.New()
	handler := func(c echo.Context) error {
		pid, _ := ProvenanceIDFromContext(c.Request().Context())
		return c.String(http.StatusOK, pid)
	}

	h := ProvenanceIDMiddleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.NoError(t, h(c))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotEqual(t, "", rec.Header().Get(ProvenanceIDHeader))
	assert.NotEqual(t, "", rec.Body.String())
}
