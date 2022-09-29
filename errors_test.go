package webutils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	errors "github.com/cyberhorsey/errors"
	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func TestError_MarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		er     *ErrorResponse
		expect string
	}{
		{
			"nil object",
			nil,
			`null`,
		},
		{
			"nil",
			&ErrorResponse{
				Errors: nil,
			},
			`{"errors":null}`,
		},
		{
			"empty",
			&ErrorResponse{
				Errors: []Error{},
			},
			`{"errors":[]}`,
		},
		{
			"empty error",
			&ErrorResponse{
				Errors: []Error{
					{
						// empty
					},
				},
			},
			`{"errors":[{"title":""}]}`,
		},
		{
			"nonempty error",
			&ErrorResponse{
				Errors: []Error{
					{
						Cause:  fmt.Errorf("cause"),
						Key:    "key",
						Title:  "title",
						Detail: "detail",
					},
				},
			},
			`{"errors":[{"key":"key","title":"title","detail":"detail"}]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs, err := json.Marshal(tt.er)
			if err != nil {
				assert.FailNowf(t, "failed to marshal", "error: %v", err)
			}

			assert.Equal(t, tt.expect, string(bs))
		})
	}
}

func TestError_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect ErrorResponse
	}{
		{
			"nil object",
			`null`,
			ErrorResponse{},
		},
		{
			"nil",
			`{"errors":null}`,
			ErrorResponse{
				Errors: nil,
			},
		},
		{
			"empty",
			`{"errors":[]}`,
			ErrorResponse{
				Errors: []Error{},
			},
		},
		{
			"empty error",
			`{"errors":[{"title":""}]}`,
			ErrorResponse{
				Errors: []Error{
					{
						// empty
						Cause: errors.NoType.NewWithKeyAndDetail("", ""),
					},
				},
			},
		},
		{
			"nonempty error",
			`{"errors":[{"key":"key","title":"title","detail":"detail"}]}`,
			ErrorResponse{
				Errors: []Error{
					{
						Cause:  errors.NoType.NewWithKeyAndDetail("key", "title: detail"),
						Key:    "key",
						Title:  "title",
						Detail: "detail",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ErrorResponse{}
			err := json.Unmarshal([]byte(tt.input), &result)
			if err != nil {
				assert.FailNowf(t, "failed to unmarshal", "error: %v", err)
			}

			// we can't compare our objects directly as is because they have different error pointers
			bs1, _ := json.Marshal(tt.expect)
			bs2, _ := json.Marshal(result)
			assert.Equal(t, bs1, bs2)
		})
	}
}

func TestConvertGRPCCodeToStatusCode(t *testing.T) {
	tests := []struct {
		err    codes.Code
		status int
	}{
		{
			codes.OK,
			http.StatusOK,
		},
		{
			codes.InvalidArgument,
			http.StatusBadRequest,
		},
		{
			codes.PermissionDenied,
			http.StatusForbidden,
		},
		{
			codes.NotFound,
			http.StatusNotFound,
		},
		{
			codes.Unknown,
			http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d -> %d", tt.err, tt.status), func(t *testing.T) {
			assert.Equal(t, ConvertGRPCCodeToStatusCode(tt.err), tt.status)
		})
	}
}

func TestError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  Error
		want string
	}{
		{
			"empty",
			Error{},
			"",
		},
		{
			"detailOnly",
			Error{Detail: "Error detail"},
			"Error detail",
		},
		{
			"keyAndDetail",
			Error{Key: "ERR_KEY", Detail: "Error detail"},
			"ERR_KEY: Error detail",
		},
		{
			"allFields",
			Error{Key: "ERR_KEY", Title: "Error title", Detail: "Error detail"},
			"ERR_KEY: Error title: Error detail",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestErrorResponse_Error(t *testing.T) {
	tests := []struct {
		name    string
		errResp ErrorResponse
		want    string
	}{
		{
			"withCause",
			ErrorResponse{
				Errors: []Error{
					{
						Cause: errors.Wrap(
							errors.Validation.NewWithKeyAndDetail("CAUSE_ERR_KEY", "Cause error detail"),
							"wrappedError",
						),
						Key:    "ERR_KEY",
						Title:  "An error",
						Detail: "Detail",
					},
				},
			},
			"CAUSE_ERR_KEY: wrappedError: Cause error detail",
		},
		{
			"withoutCause",
			ErrorResponse{
				Errors: []Error{
					{
						Key:    "ERR_KEY",
						Title:  "An error",
						Detail: "Detail",
					},
				},
			},
			"ERR_KEY: An error: Detail",
		},
	}

	for _, tt := range tests {
		got := tt.errResp.Error()
		assert.Equal(t, tt.want, got)
	}
}

func TestRenderErrors_Single(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			"badRequest",
			errors.BadRequest.NewWithKeyAndDetail("DB_ERR", "Detail here"),
			`{"errors":[{"key":"DB_ERR","title":"Bad Request","detail":"Detail here"}]}`,
		},
		{
			"missingParameter",
			errors.MissingParameter.NewWithDetail("Detail here"),
			`{"errors":[{"title":"Missing Parameter","detail":"Detail here"}]}`,
		},
		{
			"forbidden",
			errors.Forbidden.NewWithDetail("Detail here"),
			`{"errors":[{"title":"Forbidden","detail":"Detail here"}]}`,
		},
		{
			"unauthorized",
			errors.Unauthorized.NewWithDetail("Detail here"),
			`{"errors":[{"title":"Unauthorized","detail":"Detail here"}]}`,
		},
		{
			"validation",
			errors.Validation.NewWithDetail("Detail here"),
			`{"errors":[{"title":"Unprocessable Entity","detail":"Detail here"}]}`,
		},
		{
			"invalidParameter",
			errors.InvalidParameter.NewWithDetail("Detail here"),
			`{"errors":[{"title":"Unprocessable Entity","detail":"Detail here"}]}`,
		},
		{
			"notFound",
			errors.NotFound.NewWithDetail("Detail here"),
			`{"errors":[{"title":"Not Found","detail":"Detail here"}]}`,
		},
		{
			"noType",
			errors.NoType.NewWithDetail("Detail here"),
			formatJSONString(`
{
	"errors":[
		{
			"key":"ERR_UNEXPECTED",
			"title":"Internal Server Error",
			"detail":"An unexpected error occurred. If the problem persists, please contact support@gamestop.com."
		}
	]
}`),
		},
		{
			"regularError",
			errors.New("error"),
			formatJSONString(`
{
	"errors":[
		{
			"key":"ERR_UNEXPECTED",
			"title":"Internal Server Error",
			"detail":"An unexpected error occurred. If the problem persists, please contact support@gamestop.com."
		}
	]
}`),
		},
		{
			"noDetail",
			errors.Unauthorized.New("internal validation message without detail"),
			`{"errors":[{"title":"Unauthorized"}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := RenderErrors(tt.err)
			if len(errs.Errors) > 1 {
				t.Fatalf("expected errs length to be 1, got %v", len(errs.Errors))
			}
			marshaled, err := json.Marshal(errs)
			if err != nil {
				t.Fatalf("unable to marshal json response")
			}
			if tt.expected != string(marshaled) {
				t.Fatalf("expected json string to be %v, got %v", tt.expected, string(marshaled))
			}
		})
	}
}

func TestErrorResponse_Multiple(t *testing.T) {
	tests := []struct {
		name     string
		errs     []error
		expected string
	}{
		{
			"multipleErrors",
			[]error{
				errors.BadRequest.NewWithDetail("Detail here"),
				errors.Forbidden.NewWithKeyAndDetail("ERR_FORBIDDEN", "Detail here"),
			},
			formatJSONString(`
{
	"errors":[
		{
			"title":"Bad Request",
			"detail":"Detail here"
		},
		{
			"key":"ERR_FORBIDDEN",
			"title":"Forbidden",
			"detail":"Detail here"
		}
	]
}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := RenderErrors(tt.errs...)
			marshaled, err := json.Marshal(errs)
			if err != nil {
				t.Fatalf("unable to marshal json response")
			}
			if tt.expected != string(marshaled) {
				t.Fatalf("expected json string to be %v, got %v", tt.expected, string(marshaled))
			}
		})
	}
}

func TestRenderUnexpectedError(t *testing.T) {
	resp := RenderUnexpectedError(nil)
	assert.Equal(
		t,
		"ERR_UNEXPECTED: Internal Server Error: "+
			"An unexpected error occurred. If the problem persists, please contact support@gamestop.com.",
		resp.Error(),
	)
}

func TestConvertErrorToStatusCode(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{
			"badRequest",
			errors.BadRequest.NewWithDetail("error detail"),
			400,
		},
		{
			"missingParameter",
			errors.MissingParameter.NewWithDetail("error detail"),
			400,
		},
		{
			"forbidden",
			errors.Forbidden.NewWithDetail("error detail"),
			403,
		},
		{
			"unauthorized",
			errors.Unauthorized.NewWithDetail("error detail"),
			401,
		},
		{
			"validation",
			errors.Validation.NewWithDetail("error detail"),
			422,
		},
		{
			"invalidParameter",
			errors.InvalidParameter.NewWithDetail("error detail"),
			422,
		},
		{
			"notFound",
			errors.NotFound.NewWithDetail("error detail"),
			404,
		},
		{
			"noType",
			fmt.Errorf("standard error"),
			500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertErrorToStatusCode(tt.err)
			assert.Equal(t, tt.wantStatus, got)
		})
	}
}

func TestConvertErrorToGRPCCode(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus codes.Code
	}{
		{
			"badRequest",
			errors.BadRequest.NewWithDetail("error detail"),
			codes.InvalidArgument,
		},
		{
			"missingParameter",
			errors.MissingParameter.NewWithDetail("error detail"),
			codes.InvalidArgument,
		},
		{
			"forbidden",
			errors.Forbidden.NewWithDetail("error detail"),
			codes.PermissionDenied,
		},
		{
			"unauthorized",
			errors.Unauthorized.NewWithDetail("error detail"),
			codes.PermissionDenied,
		},
		{
			"validation",
			errors.Validation.NewWithDetail("error detail"),
			codes.InvalidArgument,
		},
		{
			"invalidParameter",
			errors.InvalidParameter.NewWithDetail("error detail"),
			codes.InvalidArgument,
		},
		{
			"notFound",
			errors.NotFound.NewWithDetail("error detail"),
			codes.NotFound,
		},
		{
			"noType",
			fmt.Errorf("standard error"),
			codes.Unknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertErrorToGRPCCode(tt.err)
			assert.Equal(t, tt.wantStatus, got)
		})
	}
}

func TestLogAndRenderErrors(t *testing.T) {
	e := echo.New()

	errs := []error{fmt.Errorf("internal error"), errors.Validation.NewWithDetail("error detail")}

	e.GET("/", func(c echo.Context) error {
		return LogAndRenderErrors(c, http.StatusUnprocessableEntity, errs...)
	})

	req, _ := http.NewRequest(echo.GET, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, 422, rec.Result().StatusCode)
	assert.Contains(t, rec.Body.String(), `"detail":"error detail"`)
	assert.Contains(t, rec.Body.String(), "An unexpected error occurred")
}

func TestLogAndRenderUnexpectedError(t *testing.T) {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return LogAndRenderUnexpectedError(c, fmt.Errorf("internal error"))
	})

	req, _ := http.NewRequest(echo.GET, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, 500, rec.Result().StatusCode)
	assert.NotContains(t, rec.Body.String(), "error detail")
	assert.Contains(t, rec.Body.String(), "An unexpected error occurred")
}

func Test_Logger(t *testing.T) {
	_ = Logger()
}
