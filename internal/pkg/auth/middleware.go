package auth

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func (s *JWTService) RequireAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")

			if token == "" {
				return c.JSON(
					http.StatusUnauthorized, map[string]string{
						"error": "Authorization token required",
					},
				)
			}
			login, err := s.ParseToken(token)
			if err != nil {
				return c.JSON(
					http.StatusUnauthorized, map[string]string{
						"error": "Invalid token",
					},
				)
			}
			role, err := s.GetRole(token)
			if err != nil {
				return c.JSON(
					http.StatusUnauthorized, map[string]string{
						"error": "Invalid token",
					},
				)
			}
			c.Set("user_login", login)
			c.Set("role", role)
			return next(c)
		}
	}
}

func (s *JWTService) AuthInfoMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authInfo := struct {
				IsAuthenticated bool
				UserLogin       string
			}{
				IsAuthenticated: false,
			}

			token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
			if token != "" {
				if login, err := s.ParseToken(token); err == nil {
					authInfo.IsAuthenticated = true
					authInfo.UserLogin = login
				}
			}

			c.Set("auth_info", authInfo)
			return next(c)
		}
	}
}
