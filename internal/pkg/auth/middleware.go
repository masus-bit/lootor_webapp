package auth

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func (s *JWTService) EchoMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")

			userLogin, err := s.ParseToken(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Unauthorized",
				})
			}

			c.Set("userLogin", userLogin)
			return next(c)
		}
	}
}
