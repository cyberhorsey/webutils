package webutils

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cyberhorsey/errors"
	jwt "github.com/dgrijalva/jwt-go"
	echo "github.com/labstack/echo/v4"
)

// ContextKey is a key to add to request contexts
type ContextKey string // keys can not be untyped strings

// Context keys
const (
	ContextKeyJWT       = "jwt"
	ContextKeyJWTClaims = "jwt-claims"
)

// JWTType is a type of token
type JWTType string

const (
	// JWTAccess is a short-lived access token to grant access to the application
	JWTAccess JWTType = "access"

	// JWTRefresh is a long-lived refresh token used to gain access tokens
	JWTRefresh JWTType = "refresh"
)

// Claims contains jwt.Token.Claims data
type Claims struct {
	jwt.StandardClaims
	Type     string `json:"type"`
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
}

// AuthorizedUserID returns the authorized UserID from the claims
func (c *Claims) AuthorizedUserID() uint {
	return c.UserID
}

// Authorized indicates whether the claims are valid
func (c *Claims) Authorized() error {
	if err := c.StandardClaims.Valid(); err != nil {
		return err
	}

	return nil
}

// CreateJWT creates a JWT string for the provided Claims, signed with the RSA key.
func CreateJWT(claims Claims, key *rsa.PrivateKey) (string, error) {
	if err := claims.Valid(); err != nil {
		return "", errors.Wrap(err, "claims.Valid()")
	}

	if key == nil {
		return "", ErrNoKey
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)

	return token.SignedString(key)
}

// GetClaimsFromJWT parses and returns the Claims from the provided RSA encrypted token string
func GetClaimsFromJWT(token string, key *rsa.PublicKey) (*Claims, error) {
	if token == "" {
		return nil, ErrNoToken
	}

	if key == nil {
		return nil, ErrNoKey
	}

	claims := &Claims{}

	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// GetClaimsFromJWTToken creates Claims from a *jwt.Token
func GetClaimsFromJWTToken(token *jwt.Token) (*Claims, error) {
	if token == nil {
		return nil, ErrNoToken
	}

	output, err := json.Marshal(token.Claims)
	if err != nil {
		return nil, err
	}

	claims := &Claims{}

	err = json.Unmarshal(output, claims)

	if err != nil {
		return nil, err
	}

	return claims, nil
}

var defaultJWTMiddlewareSkipper = func(c echo.Context) bool {
	switch c.Request().URL.Path {
	case "/":
		return true
	case "/health":
		return true
	}
	return false
}

// JWTMiddlewareOpts contains the options for ConfigureJWTMiddleware
type JWTMiddlewareOpts struct {
	PublicKey func(c echo.Context) (*rsa.PublicKey, error)
	Skipper   func(c echo.Context) bool
}

// jwtMiddleware is a wrapper for echo jwt middleware
type jwtMiddleware struct {
	PublicKey func(c echo.Context) (*rsa.PublicKey, error)
	Skipper   func(c echo.Context) bool
}

// ConfigureJWTMiddleware configures JWT middleware
func ConfigureJWTMiddleware(opts JWTMiddlewareOpts) (echo.MiddlewareFunc, error) {
	mw := jwtMiddleware(opts)
	if mw.PublicKey == nil {
		return nil, ErrNoPublicKeyFunction
	}

	if mw.Skipper == nil {
		mw.Skipper = defaultJWTMiddlewareSkipper
	}

	return mw.Handler, nil
}

func (mw *jwtMiddleware) Handler(
	next echo.HandlerFunc,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		// skip jwt?
		if mw.Skipper(c) {
			return next(c)
		}

		pk, err := mw.PublicKey(c)
		if err != nil {
			if errors.GetType(err) != errors.NoType {
				return LogAndRenderErrors(c, ConvertErrorToStatusCode(err), err)
			}

			return LogAndRenderUnexpectedError(c, err)
		}

		claims, err := GetClaimsFromBearerJWT(c.Request().Header.Get("Authorization"), pk)
		if err != nil {
			return LogAndRenderErrors(c, http.StatusUnauthorized, errors.Wrap(err, "GetClaimsFromBearerJWT"))
		}

		if claims.Type != string(JWTAccess) {
			return LogAndRenderErrors(c, http.StatusUnauthorized, ErrAuthorizationAccessTokenRequired)
		}

		c.Set(ContextKeyJWTClaims, claims)
		ctx := context.WithValue(c.Request().Context(), ContextKey(ContextKeyJWTClaims), claims)

		header := c.Request().Header.Get(echo.HeaderAuthorization)
		jwt := header[len(bearerPrefix):]
		c.Set(ContextKeyJWT, jwt)
		ctx = context.WithValue(ctx, ContextKey(ContextKeyJWT), jwt)

		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

const (
	bearerPrefix = `Bearer `
)

// GetClaimsFromBearerJWT gets claims from a Bearer JWT
func GetClaimsFromBearerJWT(
	token string,
	jwtPublicKey *rsa.PublicKey,
) (*Claims, error) {
	if token == "" {
		return nil, ErrAuthorizationAccessTokenRequired
	}

	if strings.HasPrefix(token, bearerPrefix) {
		claims, err := GetClaimsFromJWT(token[len(bearerPrefix):], jwtPublicKey)
		if err != nil {
			return nil, errors.WithCause(ErrAuthorizationTokenInvalid, err)
		}

		return claims, nil
	}

	// enforce Bearer prefix
	return nil, ErrAuthorizationBearerRequired
}

// GetJWTClaimsFromEchoContext retrieves the *webutils.Claims from the provided echo.Context.
func GetJWTClaimsFromEchoContext(c echo.Context) (*Claims, error) {
	claims, ok := c.Get(ContextKeyJWTClaims).(*Claims)
	if !ok {
		return nil, ErrNoJWTClaimsInContext
	}

	return claims, nil
}

// GetJWTFromEchoContext retrieves the JWT from the provided echo.Context.
func GetJWTFromEchoContext(c echo.Context) (string, error) {
	jwt, ok := c.Get(ContextKeyJWT).(string)
	if !ok {
		return "", ErrNoJWTInContext
	}

	return jwt, nil
}

// GetJWTFromContext retrieves the JWT from the provided context.Context.
func GetJWTFromContext(ctx context.Context) (string, error) {
	jwt, ok := ctx.Value(ContextKey(ContextKeyJWT)).(string)
	if !ok {
		return "", ErrNoJWTInContext
	}

	return jwt, nil
}

// GetJWTClaimsFromContext retrieves the *webutils.Claims from the provided context.Context.
func GetJWTClaimsFromContext(ctx context.Context) (*Claims, error) {
	claims, ok := ctx.Value(ContextKey(ContextKeyJWTClaims)).(*Claims)
	if !ok || claims == nil {
		return nil, ErrNoJWTClaimsInContext
	}

	return claims, nil
}
