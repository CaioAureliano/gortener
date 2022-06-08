package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/CaioAureliano/gortener/internal/auth/service"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ConfigJwt() echo.MiddlewareFunc {
	config := middleware.JWTConfig{
		ParseTokenFunc: customParseToken(),
	}

	return middleware.JWTWithConfig(config)
}

func customParseToken() func(string, echo.Context) (interface{}, error) {
	return func(auth string, c echo.Context) (interface{}, error) {

		token, err := jwt.Parse(auth, keyFunc())
		if err != nil {
			return nil, err
		}

		if !token.Valid {
			return nil, errors.New("invalid token")
		}

		return token, nil
	}
}

func keyFunc() jwt.Keyfunc {
	return func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != "HS256" {
			return nil, fmt.Errorf("unexpected jwt signing method: %s", t.Header["alg"])
		}

		if isExpiredToken(t.Claims) {
			return nil, echo.NewHTTPError(http.StatusUnauthorized, "expired token")
		}

		return []byte(service.TOKEN_SECRET_KEY), nil
	}
}

func isExpiredToken(claims jwt.Claims) bool {
	expiresAt := int64(getExpiresAtFromClaims(claims))
	now := time.Now().UnixMilli()
	return now > expiresAt
}

func getExpiresAtFromClaims(claims jwt.Claims) float64 {
	mapClaims := claims.(jwt.MapClaims)
	exp, exists := mapClaims["expires_at"]
	if !exists {
		return 0
	}
	return exp.(float64)
}
