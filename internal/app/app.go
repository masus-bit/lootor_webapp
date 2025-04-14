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
		100*time.Minute,
		7*24*time.Hour,
	)

	userRepo := repositories.NewUsersRepository(db)
	ciRepo := repositories.NewCiRepository(db)
	colRepo := repositories.NewCollectionsRepository(db)
	tagRepo := repositories.NewTagsRepository(db)
	platformRepo := repositories.NewPlatformsRepository(db)
	entityRepo := repositories.NewEntitiesRepository(db)
	eventsRepo := repositories.NewEventsRepository(db)
	err = db.AutoMigrate(&models.Users{}, &models.Platforms{}, models.Events{}, &models.Collections{}, &models.Tags{}, &models.CollectionItems{}, &models.Entities{})
	if err != nil {
		return nil, err
	}
	userService := services.NewUserService(userRepo, jwtService, ciRepo)
	ciService := services.NewCiService(ciRepo, eventsRepo, colRepo, userRepo, platformRepo, entityRepo)
	colService := services.NewCollectionService(colRepo, tagRepo, userRepo, eventsRepo, ciRepo)
	entitiesService := services.NewEntitiesService(entityRepo)
	platformsService := services.NewPlatformsService(platformRepo)
	tagsService := services.NewTagsService(tagRepo)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 60 * time.Second,
	}))

	routes.RegisterRoutes(e, jwtService, *userService)
	routes.TagsRouter(e, jwtService, *tagsService)
	routes.RegisterCollectionsRoutes(e, jwtService, *colService)
	routes.RegisterCollectionItemsRoutes(e, jwtService, *ciService)
	routes.EntitiesRouter(e, jwtService, *entitiesService)
	routes.PlatformsRouter(e, jwtService, *platformsService)

	return &App{Echo: e}, nil
}
