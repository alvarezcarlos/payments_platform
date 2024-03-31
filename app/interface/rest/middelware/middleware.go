package middelware

import (
	"github.com/alvarezcarlos/payment/app/config"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type Middleware interface {
	JwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc
}

type middleware struct {
}

func NewMiddleware() Middleware {
	return &middleware{}
}

func (m *middleware) JwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return echo.ErrUnauthorized
		}

		claims := &jwt.StandardClaims{}
		jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Config().SecretKey), nil
		})

		if err != nil || !jwtToken.Valid {
			return echo.ErrUnauthorized
		}

		return next(c)
	}
}
