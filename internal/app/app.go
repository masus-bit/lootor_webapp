package app

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"lootor/internal/config"
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
	"lootor/internal/core/routes"
	"lootor/internal/core/services"
	"lootor/internal/pkg/auth"
	"lootor/internal/pkg/database"
	"os"
	"time"
)

type App struct {
	Echo *echo.Echo
}

func NewEchoApp(cfg *config.Config) (*App, error) {
	_ = godotenv.Load()

	e := echo.New()

	// Инициализация БД
	db, err := database.InitDB(&cfg.Database)
	if err != nil {
		return nil, err
	}

	jwtService := auth.NewJWTService(
		os.Getenv("JWT_SECRET"),
		15*time.Minute,
		7*24*time.Hour,
	)

	userRepo := repositories.NewUsersRepository(db)
	ciRepo := repositories.NewCiRepository(db)
	err = db.AutoMigrate(&models.Users{}, &models.Collections{}, &models.Tags{}, &models.CollectionItems{}, &models.Platforms{}, &models.Entities{}, models.Events{})
	if err != nil {
		return nil, err
	}
	userService := services.NewUserService(userRepo, jwtService, ciRepo)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 60 * time.Second,
	}))

	routes.RegisterRoutes(e, jwtService, *userService)

	return &App{Echo: e}, nil
}
