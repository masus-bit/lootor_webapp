package app

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"lootor/internal/config"
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
	"lootor/internal/core/routes"
	"lootor/internal/core/services"
	"lootor/internal/pkg/auth"
	"lootor/internal/pkg/database"
	"lootor/internal/pkg/elasticsearch"
	"lootor/internal/pkg/mail"
	"lootor/internal/pkg/s3"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type App struct {
	Echo *echo.Echo
}

func getConfigPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "..", "elasticConfig.json")
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

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

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
	mailService := mail.NewMailService("smtp.yandex.ru", 465, "noreply@lootor.me", "cytuhekbl13", `"Lootor" <noreply@lootor.me>`)
	userService := services.NewUserService(userRepo, jwtService, ciRepo, mailService)
	ciService := services.NewCiService(ciRepo, eventsRepo, colRepo, userRepo, platformRepo, entityRepo)
	colService := services.NewCollectionService(colRepo, tagRepo, userRepo, eventsRepo, ciRepo)
	entitiesService := services.NewEntitiesService(entityRepo)
	platformsService := services.NewPlatformsService(platformRepo)
	tagsService := services.NewTagsService(tagRepo)
	s3Service := s3.NewS3Service(redisClient)
	searchService, err := elasticsearch.NewElasticService(getConfigPath())
	if err != nil {
		return nil, err
	}

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
	routes.S3Router(e, jwtService, *s3Service)
	routes.SearchRouter(e, jwtService, *searchService)
	routes.ReindexRouter(e, jwtService, *searchService, *userRepo, *ciRepo, *colRepo, *tagRepo, *entityRepo)

	return &App{Echo: e}, nil
}
