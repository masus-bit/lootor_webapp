package user

import (
	"github.com/go-chi/chi/v5"
	"lootor/internal/pkg/auth"
)

func RegisterRoutes(r chi.Router, jwtService *auth.JWTService, userService Service) {
	controller := NewUserController(userService)

	r.Group(func(r chi.Router) {
		r.Post("/public/auth/signup", controller.SignUp)
		r.Post("/public/auth/signin", controller.SignIn)
	})

	r.Group(func(r chi.Router) {
		r.Use(jwtService.RequireAuth)
		r.Get("/secured/user", controller.GetByLogin)
	})
}
