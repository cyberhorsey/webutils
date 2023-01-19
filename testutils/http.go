package testutils

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/cyberhorsey/webutils"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

// NewAuthenticatedRequest creates an authenticated request with the specified account and user IDs
func NewAuthenticatedRequest(
	accountID uint,
	userID uint,
	privateKey *rsa.PrivateKey,
	method string,
	path string,
	body interface{},
) *http.Request {
	r := httptest.NewRequest(method, path, strings.NewReader(requestBodyToString(body)))

	if accountID != 0 && userID != 0 {
		setTokenAuthorizationHeader(r, generateSignedJWT(accountID, userID, privateKey))
	}

	setContentTypeHeader(r, "application/json")

	return r
}

// NewAuthenticatedRequestWithJWT creates an authenticated request with the specified JWT
func NewAuthenticatedRequestWithJWT(token, method, path string, body interface{}) *http.Request {
	r := httptest.NewRequest(method, path, strings.NewReader(requestBodyToString(body)))
	setTokenAuthorizationHeader(r, token)
	setContentTypeHeader(r, "application/json")

	return r
}

// NewAuthenticatedRequestWithCookie creates an authenticated request with the specified JWT
func NewAuthenticatedRequestWithCookie(cookieName, cookieValue, method, path string, body interface{}) *http.Request {
	r := httptest.NewRequest(method, path, strings.NewReader(requestBodyToString(body)))
	r.AddCookie(
		&http.Cookie{
			Name:  cookieName,
			Value: cookieValue,
		},
	)
	setContentTypeHeader(r, "application/json")

	return r
}

// NewUnauthenticatedRequest creates an unauthenticated request
func NewUnauthenticatedRequest(method, path string, body interface{}) *http.Request {
	r := httptest.NewRequest(method, path, strings.NewReader(requestBodyToString(body)))
	setContentTypeHeader(r, "application/json")

	return r
}

func requestBodyToString(body interface{}) string {
	if bs, ok := body.(string); ok {
		return bs
	}

	if bb, ok := body.([]byte); ok {
		return string(bb)
	}

	bb, err := json.Marshal(body)
	if err != nil {
		panic(fmt.Sprintf("unexpected json.Marshal error: %v", err))
	}

	return string(bb)
}

func generateSignedJWT(accountID, userID uint, privateKey *rsa.PrivateKey) string {
	if accountID == 0 || userID == 0 {
		return "badjwt"
	}

	token, err := webutils.CreateJWT(
		webutils.Claims{
			Username: "tester",
			UserID:   userID,
			Type:     "access",
			StandardClaims: jwt.StandardClaims{
				Issuer:    "gamestop",
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
				NotBefore: time.Now().Unix(),
			},
		},
		privateKey,
	)
	if err != nil {
		panic(err)
	}

	return token
}

func setTokenAuthorizationHeader(r *http.Request, token string) {
	r.Header.Set("Authorization", "Bearer "+token)
}

func setContentTypeHeader(r *http.Request, value string) {
	r.Header.Set("Content-Type", value)
}

// AssertStatusAndBody asserts the status and body of an HTTP response
func AssertStatusAndBody(
	t *testing.T,
	rec *httptest.ResponseRecorder,
	wantStatus int,
	wantBodyRegexpMatches []string,
) {
	assert.Equal(t, wantStatus, rec.Code)
	gotBody := rec.Body.String()

	for _, wantMatch := range wantBodyRegexpMatches {
		assert.Regexp(t, regexp.MustCompile(wantMatch), gotBody)
	}
}

// AssertNotRegexpy asserts the matches do not exist in string s
func AssertNotRegexp(t *testing.T, matches []string, s string) {
	for _, match := range matches {
		assert.NotRegexp(t, regexp.MustCompile(match), s)
	}
}
