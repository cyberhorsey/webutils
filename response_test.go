package webutils

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_JSON(t *testing.T) {
	w := httptest.NewRecorder()

	err := JSON(w, 200, nil)
	if err != nil {
		t.Fatalf("Test_JSON couldnt write to recorder, got err %v", nil)
	}

	if w.Result().StatusCode != 200 {
		t.Fatalf("Test_JSON did not write status")
	}
}

// MockNotificationService is responsible for publishing notifications
type MockNotificationService struct {
	NotifyFn func(
		ctx context.Context,
		opts NotificationOpts,
	) error
}

// Notify MockNotificationService.Notify
func (svc *MockNotificationService) Notify(
	ctx context.Context,
	opts NotificationOpts,
) error {
	return svc.NotifyFn(ctx, opts)
}

func Test_LogNotifyAndRenderUnexpectedError(t *testing.T) {
	resultErr := fmt.Errorf("failed to publish")

	var result NotificationOpts

	var wg sync.WaitGroup

	wg.Add(1)

	nSvc := &MockNotificationService{
		NotifyFn: func(
			ctx context.Context,
			opts NotificationOpts,
		) error {
			defer wg.Done()
			result = opts
			return nil
		},
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		// this will always return the same error
		return LogNotifyAndRenderUnexpectedError(
			c,
			nSvc,
			resultErr,
		)
	})

	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, 500, rec.Result().StatusCode)
	assert.Contains(t, rec.Body.String(), "An unexpected error occurred")

	// wait for notify
	wg.Wait()
	assert.Equal(
		t,
		result,
		NotificationOpts{
			Subject:  "Unexpected Error",
			Message:  "An unexpected error occurred that resulted in a 500 Internal Server Error",
			Priority: NotificationPriorityHigh,
			Error:    resultErr,
			Metadata: nil,
		},
	)
}

func Test_LogNotifyAndRenderUnexpectedError_NotifyError(t *testing.T) {
	resultErr := fmt.Errorf("failed to publish")

	var result NotificationOpts

	var wg sync.WaitGroup

	wg.Add(1)

	nSvc := &MockNotificationService{
		NotifyFn: func(
			ctx context.Context,
			opts NotificationOpts,
		) error {
			defer wg.Done()
			result = opts
			return fmt.Errorf("some other error")
		},
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		// this will always return the same error
		return LogNotifyAndRenderUnexpectedError(
			c,
			nSvc,
			resultErr,
		)
	})

	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, 500, rec.Result().StatusCode)
	assert.Contains(t, rec.Body.String(), "An unexpected error occurred")

	// wait for notify
	wg.Wait()
	assert.Equal(
		t,
		result,
		NotificationOpts{
			Subject:  "Unexpected Error",
			Message:  "An unexpected error occurred that resulted in a 500 Internal Server Error",
			Priority: NotificationPriorityHigh,
			Error:    resultErr,
			Metadata: nil,
		},
	)
}

func Test_CheckResponse(t *testing.T) {
	tests := []struct {
		name    string
		r       *http.Response
		wantErr bool
	}{
		{
			"noError",
			&http.Response{
				StatusCode: 201,
			},
			false,
		},
		{
			"errorBadBody",
			&http.Response{
				StatusCode: 401,
				Body:       ioutil.NopCloser(strings.NewReader("Hello")),
			},
			true,
		},
		{
			"error",
			&http.Response{
				StatusCode: 401,
				Body:       ioutil.NopCloser(strings.NewReader(`{"errors": "title": "title", "detail": "detail"}`)),
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckResponse(tt.r)
			if tt.wantErr {
				assert.NotEqual(t, nil, err)
			} else {
				assert.Equal(t, nil, err)
			}
		})
	}
}
