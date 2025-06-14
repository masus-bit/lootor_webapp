package app

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	echoSwagger "github.com/swaggo/echo-swagger"
	"log"
	"lootor/internal/config"
	"lootor/internal/core/controllers"
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
	"lootor/internal/core/routes"
	"lootor/internal/core/services"
	"lootor/internal/pkg/auth"
	"lootor/internal/pkg/database"
	"lootor/internal/pkg/elasticsearch"
	"lootor/internal/pkg/feedback"
	"lootor/internal/pkg/mail"
	"lootor/internal/pkg/payment"
	"lootor/internal/pkg/s3"
	"lootor/internal/pkg/utils"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
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

	e.GET("/public/swagger/*", echoSwagger.WrapHandler)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"https://dev.lootor.me",
			"https://lootor.me",
			"https://www.lootor.me",
			"http://localhost:3000",
			"http://localhost:4173",
			"*",
		},
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.PUT,
			echo.DELETE,
			echo.OPTIONS,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			"X-Requested-With",
		},
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	// Инициализация БД
	db, err := database.InitDB(&cfg.Database)
	if err != nil {
		if err := database.Reconnect(&cfg.Database); err != nil {
			log.Printf("Reconnection failed: %v", err)
		}
	}

	jwtService := auth.NewJWTService(
		os.Getenv("JWT_SECRET"),
		100*time.Minute,
		7*24*time.Hour,
	)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: "",
		DB:       0,
	})
	searchService, err := elasticsearch.NewElasticService(getConfigPath())
	if err != nil {
		return nil, err
	}

	userRepo := repositories.NewUsersRepository(db, searchService)
	ciRepo := repositories.NewCiRepository(db, searchService)
	colRepo := repositories.NewCollectionsRepository(db, searchService)
	tagRepo := repositories.NewTagsRepository(db, searchService)
	platformRepo := repositories.NewPlatformsRepository(db)
	entityRepo := repositories.NewEntitiesRepository(db, searchService)
	eventsRepo := repositories.NewEventsRepository(db)
	itemTypesRepo := repositories.NewItemTypesRepository(db)
	wlRepo := repositories.NewWLRepository(db)
	subRepo := repositories.NewSubscriptionRepository(db)
	err = db.AutoMigrate(&models.Users{}, &models.Platforms{}, &models.Events{}, &models.Collections{}, &models.Tags{}, &models.CollectionItems{}, &models.Entities{}, &models.ItemTypes{}, &models.WishListItems{})
	if err != nil {
		return nil, err
	}

	host := os.Getenv("YANDEX_POST_HOST")
	port := os.Getenv("YANDEX_POST_PORT")
	portInt, _ := strconv.Atoi(port)
	user := os.Getenv("YANDEX_POST_USER")
	password := os.Getenv("YANDEX_POST_PASSWORD")

	s3Service := s3.NewS3Service(redisClient)
	mailService := mail.NewMailService(host, portInt, user, password, `"Lootor" <noreply@lootor.me>`)
	userService := services.NewUserService(userRepo, jwtService, ciRepo, mailService, eventsRepo)
	ciService := services.NewCiService(ciRepo, eventsRepo, colRepo, userRepo, platformRepo, entityRepo, s3Service, itemTypesRepo)
	colService := services.NewCollectionService(colRepo, tagRepo, userRepo, eventsRepo, ciRepo, s3Service)
	entitiesService := services.NewEntitiesService(entityRepo)
	platformsService := services.NewPlatformsService(platformRepo)
	itemTypesService := services.NewItemTypesService(itemTypesRepo)
	tagsService := services.NewTagsService(tagRepo)
	fbService := feedback.NewFeedbackService()
	wlService := services.NewWLService(wlRepo, userRepo, ciRepo, eventsRepo)
	enrichedCIService := utils.NewEnrichedCIService(ciService)
	subService := services.NewSubscriptionService(subRepo)
	paymentService := payment.NewPayService(userRepo, userService, subService)

	eventsService := services.NewEventsService(eventsRepo, userRepo)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 60 * time.Second,
	}))

	routes.RegisterRoutes(e, jwtService, *userService)
	routes.TagsRouter(e, jwtService, *tagsService)
	routes.RegisterCollectionsRoutes(e, jwtService, *colService)
	routes.RegisterCollectionItemsRoutes(e, jwtService, *ciService, *enrichedCIService)
	routes.EntitiesRouter(e, jwtService, *entitiesService)
	routes.PlatformsRouter(e, jwtService, *platformsService)
	routes.S3Router(e, jwtService, *s3Service)
	routes.SearchRouter(e, jwtService, *searchService)
	routes.ReindexRouter(e, jwtService, *searchService, *userRepo, *ciRepo, *colRepo, *tagRepo, *entityRepo)
	routes.EventsRouter(e, jwtService, *eventsService)
	routes.FeedbackRouter(e, jwtService, *fbService)
	routes.ItemTypesRouter(e, jwtService, *itemTypesService)
	routes.WLRouter(e, jwtService, *wlService)
	routes.PaymentRouter(e, jwtService, *paymentService)

	controllers.NewReindexController(searchService, userRepo, colRepo, ciRepo, tagRepo, entityRepo).ReindexInternal(context.Background())
	c := cron.New()
	_, err = c.AddFunc("@midnight", func() { userService.CheckExpiredSubscriptions() })
	if err != nil {
		return nil, err
	}
	c.Start()
	return &App{Echo: e}, nil
}
