package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"lootor/internal/config"
	"lootor/internal/modules/user"
	"lootor/internal/pkg/auth"
	"lootor/internal/pkg/database"
	"os"
	"time"
)

type App struct {
	Router *chi.Mux
}

func NewChiApp(cfg *config.Config) (*App, error) {
	_ = godotenv.Load()

	db, err := database.InitDB(&cfg.Database)

	if err != nil {
		return nil, err
	}

	jwtService := auth.NewJWTService(
		os.Getenv("JWT_SECRET"),
		15*time.Minute, // Access token expiration
		7*24*time.Hour, // Refresh token expiration
	)

	userRepo := user.NewRepository(db) // Предполагается, что репозиторий реализован

	db.AutoMigrate(&user.User{})

	userService := user.NewService(userRepo, jwtService)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(auth.JSONMiddleware)

	user.RegisterRoutes(r, jwtService, *userService)

	//user.RegisterRoutes(r, db)
	//auth.RegisterRoutes(r, db)

	return &App{Router: r}, nil
}
