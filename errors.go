package webutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	qerrors "github.com/cyberhorsey/errors"
	echo "github.com/labstack/echo/v4"
	"github.com/neko-neko/echo-logrus/v2/log"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
)

var logger *log.MyLogger

func init() {
	logger = log.Logger()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg: "message",
		},
	})
}

// sentinel errors
var (
	ErrNoClaims                  = errors.New("claims is required")
	ErrNoKey                     = errors.New("key is required")
	ErrNoSecret                  = errors.New("secret is required")
	ErrNoToken                   = errors.New("token is required")
	ErrInvalidToken              = errors.New("jwt token not valid")
	ErrNoJWTClaimsInContext      = errors.New("jwt claim missing from context")
	ErrNoJWTInContext            = errors.New("jwt missing from context")
	ErrNoNotificationMessage     = qerrors.New("message is required")
	ErrNoPublicKeyFunction       = qerrors.New("public key func is required")
	ErrAuthorizationTokenInvalid = qerrors.Unauthorized.NewWithKeyAndDetail(
		"ERR_AUTHORIZATION_TOKEN_INVALID",
		"Authorization token is invalid",
	)
	ErrAuthorizationAccessTokenRequired = qerrors.Unauthorized.NewWithKeyAndDetail(
		"ERR_AUTHORIZATION_ACCESS_TOKEN_REQUIRED",
		"A valid Authorization access token is required",
	)
	ErrAuthorizationBearerRequired = qerrors.Unauthorized.NewWithKeyAndDetail(
		"ERR_AUTHORIZATION_BEARER_REQUIRED",
		"Authorization Bearer is required before token",
	)
)

// Logger returns a middleware that logs HTTP requests.
func Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Request().URL.Path
			if strings.Contains(path, "health") || strings.Contains(path, "metric") {
				return nil
			}

			req := c.Request()
			res := c.Response()
			start := time.Now()

			var err error
			if err = next(c); err != nil {
				c.Error(err)
			}

			stop := time.Now()

			pid := req.Header.Get(ProvenanceIDHeader)
			if pid == "" {
				pid = res.Header().Get(ProvenanceIDHeader)
			}

			rid, ok := RequestIDFromContext(c.Request().Context())
			if !ok {
				rid = ""
			}

			logger := log.Logger()
			logger.SetFormatter(&logrus.JSONFormatter{
				TimestampFormat: time.RFC3339,
				FieldMap: logrus.FieldMap{
					logrus.FieldKeyMsg: "message",
				},
			})

			l := logger.WithFields(logrus.Fields{
				"provenanceId": pid,
				"requestId":    rid,
				"ip":           c.RealIP(),
				"host":         req.Host,
				"method":       req.Method,
				"uri":          req.RequestURI,
				"status":       res.Status,
				"latency":      stop.Sub(start).Seconds(),
				"referer":      req.Referer(),
				"userAgent":    req.UserAgent(),
			})

			if res.Status == http.StatusOK || res.Status == http.StatusNoContent {
				l.Logger.Out = os.Stdout
				l.Info("successful http call")
			} else {
				l.Logger.Out = os.Stderr
				l.Error(fmt.Sprintf("%+v", err))
			}

			return err
		}
	}
}

// Error is a struct we return through RenderErrors to be able to return multiple errors at once
// from our API.
type Error struct {
	Cause  error  `json:"-"`
	Key    string `json:"key,omitempty"`
	Title  string `json:"title"`
	Detail string `json:"detail,omitempty"`
}

func (e Error) Error() string {
	if e.Cause != nil {
		return e.Cause.Error()
	}

	parts := make([]string, 0)
	if e.Key != "" {
		parts = []string{e.Key}
	}

	if e.Title != "" {
		parts = append(parts, e.Title)
	}

	if e.Detail != "" {
		parts = append(parts, e.Detail)
	}

	return strings.Join(parts, ": ")
}

// newUnexpectedError wrapps err in our defactor Internal Server Error
func newUnexpectedError(err error) Error {
	return Error{
		Cause:  err,
		Key:    "ERR_UNEXPECTED",
		Title:  http.StatusText(http.StatusInternalServerError),
		Detail: "An unexpected error occurred.",
	}
}

// ErrorResponse contains a collection of errors
type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

// MarshalJSON marshals errors
func (er ErrorResponse) MarshalJSON() (
	[]byte,
	error,
) {
	bs, err := json.Marshal(er.Errors)
	if err != nil {
		return nil, qerrors.Wrap(err, "json.Marshal(er.Errors)")
	}

	return append(
		[]byte("{\"errors\":"),
		append(
			bs,
			[]byte("}")...,
		)...,
	), nil
}

// UnmarshalJSON unmarshals errors
func (er *ErrorResponse) UnmarshalJSON(bs []byte) error {
	var errorsResult struct {
		Errors []Error `json:"errors"`
	}

	if err := json.Unmarshal(bs, &errorsResult); err != nil {
		return qerrors.Wrap(err, "json.Unmarshal(bs, &errorsResult)")
	}

	er.Errors = errorsResult.Errors

	for i := range er.Errors {
		detail := ""
		if er.Errors[i].Title != "" {
			detail = strings.Join(
				[]string{
					er.Errors[i].Title,
					er.Errors[i].Detail,
				},
				": ",
			)
		} else {
			detail = er.Errors[i].Detail
		}

		er.Errors[i].Cause = qerrors.NoType.NewWithKeyAndDetail(er.Errors[i].Key, detail)
	}

	return nil
}

// Error messages combined
func (er ErrorResponse) Error() string {
	var errMsgs []string
	for _, err := range er.Errors {
		errMsgs = append(errMsgs, err.Error())
	}

	return strings.TrimSuffix(strings.Join(errMsgs, "; "), ";")
}

// RenderErrors prepares collection of errors in an ErrorResponse
func RenderErrors(errs ...error) ErrorResponse {
	convertedErrs := make([]Error, 0)
	for _, err := range errs {
		convertedErrs = append(convertedErrs, convertError(err))
	}

	return ErrorResponse{Errors: convertedErrs}
}

func convertError(err error) Error {
	// If the error is already a webutils.Error; return it
	if werr, ok := err.(Error); ok {
		return werr
	}

	var title string

	switch qerrors.GetType(err) {
	case qerrors.BadRequest:
		title = http.StatusText(http.StatusBadRequest)
	case qerrors.Forbidden:
		title = http.StatusText(http.StatusForbidden)
	case qerrors.Unauthorized:
		title = http.StatusText(http.StatusUnauthorized)
	case qerrors.Validation:
		title = http.StatusText(http.StatusUnprocessableEntity)
	case qerrors.InvalidParameter:
		title = http.StatusText(http.StatusUnprocessableEntity)
	case qerrors.NotFound:
		title = http.StatusText(http.StatusNotFound)
	case qerrors.MissingParameter:
		title = "Missing Parameter"
	case qerrors.NoType:
		// Errors of unknown/default type should not be exposed
		return newUnexpectedError(err)
	}

	return Error{
		Cause:  err,
		Key:    qerrors.Key(err),
		Title:  title,
		Detail: qerrors.Detail(err),
	}
}

// RenderUnexpectedError renders an unexpected error
func RenderUnexpectedError(err error) ErrorResponse {
	e := make([]Error, 0)

	e = append(e, newUnexpectedError(err))

	errResp := ErrorResponse{Errors: e}

	return errResp
}

// ConvertErrorToStatusCode converts err to an HTTP status code using qerrors.GetType. If the
// type is qerrors.NoType or is unknown, http.StatusInternalServerError is used.
func ConvertErrorToStatusCode(err error) int {
	switch qerrors.GetType(err) {
	case qerrors.BadRequest, qerrors.MissingParameter:
		return http.StatusBadRequest
	case qerrors.Forbidden:
		return http.StatusForbidden
	case qerrors.Unauthorized:
		return http.StatusUnauthorized
	case qerrors.Validation:
		return http.StatusUnprocessableEntity
	case qerrors.InvalidParameter:
		return http.StatusUnprocessableEntity
	case qerrors.NotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

// ConvertErrorToGRPCCode converts err to a GRPC error code using qerrors.GetType. If the
// type is qerrors.NoType or is unknown, codes.Unknown is used.
func ConvertErrorToGRPCCode(err error) codes.Code {
	switch qerrors.GetType(err) {
	case qerrors.BadRequest, qerrors.MissingParameter,
		qerrors.Validation, qerrors.InvalidParameter:
		return codes.InvalidArgument
	case qerrors.Forbidden, qerrors.Unauthorized:
		return codes.PermissionDenied
	case qerrors.NotFound:
		return codes.NotFound
	default:
		return codes.Unknown
	}
}

// ConvertGRPCCodeToStatusCode converts err to an HTTP status code. If the
// type is codes.Unknown or is unknown, http.StatusInternalServerError is used.
func ConvertGRPCCodeToStatusCode(err codes.Code) int {
	switch err {
	case codes.OK:
		return http.StatusOK
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.NotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

// LogAndRenderErrors logs and renders errs to JSON with the provided statusCode.
func LogAndRenderErrors(c echo.Context, statusCode int, errs ...error) error {
	errResp := RenderErrors(errs...)

	// Log error stack trace
	pid, ok := ProvenanceIDFromContext(c.Request().Context())
	if !ok {
		pid = ""
	}

	rid, ok := RequestIDFromContext(c.Request().Context())
	if !ok {
		rid = ""
	}

	for _, err := range errs {
		logger.WithFields(logrus.Fields{"provenanceId": pid, "requestId": rid}).
			Error(err)
	}

	jsonErr := c.JSON(statusCode, errResp)
	if jsonErr != nil {
		logger.WithFields(logrus.Fields{"provenanceId": pid, "requestId": rid}).Error(jsonErr)
	}

	return errResp
}

// LogAndRenderUnexpectedError logs the stack trace for err and renders a generic internal server
// error message via `RenderUnexpectedAPIError()`.
func LogAndRenderUnexpectedError(c echo.Context, err error) error {
	// Log error stack trace
	pid, ok := ProvenanceIDFromContext(c.Request().Context())
	if !ok {
		pid = ""
	}

	rid, ok := RequestIDFromContext(c.Request().Context())
	if !ok {
		rid = ""
	}

	logger.WithFields(logrus.Fields{"provenanceId": pid, "requestId": rid}).Error(err)

	jsonErr := c.JSON(http.StatusInternalServerError, RenderUnexpectedError(err))
	if jsonErr != nil {
		logger.WithFields(logrus.Fields{"provenanceId": pid, "requestId": rid}).Error(jsonErr)
	}

	// return the original error which will be logged with Echo's access log
	return err
}
