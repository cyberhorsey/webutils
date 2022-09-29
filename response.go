package webutils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	echo "github.com/labstack/echo/v4"
)

// JSON encodes data into a JSON format
func JSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(data)
}

// LogNotifyAndRenderUnexpectedError logs the stack trace for err and renders a generic internal server
// error message via `RenderUnexpectedAPIError()`.
func LogNotifyAndRenderUnexpectedError(
	c echo.Context,
	nSvc Notifier,
	err error,
) error {
	// Log error stack trace
	c.Logger().Errorf("unexpected error: %+v", err)
	// render a response before we use our notification service
	jsonErr := c.JSON(http.StatusInternalServerError, RenderUnexpectedError(err))
	if jsonErr != nil {
		c.Logger().Errorf("webutils.LogAndRenderError encountered unexpected c.JSON error: %v", jsonErr)
	}
	// notify?
	if nSvc != nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			if err := nSvc.Notify(
				ctx,
				NotificationOpts{
					Subject:  "Unexpected Error",
					Message:  "An unexpected error occurred that resulted in a 500 Internal Server Error",
					Priority: NotificationPriorityHigh,
					Error:    err,
				},
			); err != nil {
				c.Logger().Errorf("failed to publish notification error: \"%+v\"", err)
			}
		}()
	}
	// return the original error which will be logged with Echo's access log
	return err
}

// CheckResponse checks the API response for error, and returns them if present. A response is
// considered successful if it has a status code in the 200 range.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	// Attempt to decode the ErrorResponse
	errResp := new(ErrorResponse)
	if err := json.NewDecoder(r.Body).Decode(&errResp); err != nil {
		return fmt.Errorf("%v: %v", r.StatusCode, http.StatusText(r.StatusCode))
	}

	return errResp
}
