package app

import (
	"context"
	"encoding/json"
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
	"lootor/internal/infrastructure/commentsclient"
	"lootor/internal/infrastructure/eventsclient"
	"lootor/internal/infrastructure/newsclient"
	"lootor/internal/infrastructure/notificationsclient"
	"lootor/internal/infrastructure/postsclient"
	"lootor/internal/infrastructure/tagsclient"
	"lootor/internal/pkg/auth"
	"lootor/internal/pkg/database"
	"lootor/internal/pkg/elasticsearch"
	"lootor/internal/pkg/feedback"
	"lootor/internal/pkg/mail"
	"lootor/internal/pkg/migrator"
	"lootor/internal/pkg/payment"
	"lootor/internal/pkg/s3"
	"lootor/internal/pkg/utils"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

type App struct {
	Echo *echo.Echo
}

func setupMonitoring(s3Service *s3.S3Service) {
	http.HandleFunc("/debug/stats", func(w http.ResponseWriter, r *http.Request) {
		stats := s3Service.GetStats()
		json.NewEncoder(w).Encode(stats)
	})

	go func() {
		log.Println("Monitoring server started on :8081")
		log.Fatal(http.ListenAndServe(":8081", nil))
	}()
}

func getConfigPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "..", "elasticConfig.json")
}
func NewEchoApp(cfg *config.Config) (*App, error) {
	_ = godotenv.Load()

	e := echo.New()

	e.GET("/public/swagger/*", echoSwagger.WrapHandler)
	e.Server.MaxHeaderBytes = 1 << 20
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"https://dev.lootor.me",
			"https://lootor.me",
			"https://www.lootor.me",
			"https://www.dev.lootor.me",
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

	// db init
	db, err := database.InitDB(&cfg.Database)
	if err != nil {
		if err := database.Reconnect(&cfg.Database); err != nil {
			log.Printf("Reconnection failed: %v", err)
		}
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:         os.Getenv("REDIS_HOST"),
		Password:     "",
		DB:           0,
		PoolSize:     200,
		MinIdleConns: 50,
		ReadTimeout:  500 * time.Millisecond,
		WriteTimeout: 500 * time.Millisecond,
	})
	searchService, err := elasticsearch.NewElasticService(getConfigPath())
	if err != nil {
		return nil, err
	}

	newsAddr := os.Getenv("NEWS_SERVICE_ADDR")
	commentsAddr := os.Getenv("COMMENTS_SERVICE_ADDR")

	newsClient, err := newsclient.NewGRPCClient(newsAddr)
	if err != nil {
		log.Fatal("failed to create news client:", err)
	}
	commentsClient, err := commentsclient.NewGRPCCommentsClient(commentsAddr)
	if err != nil {
		log.Fatal("failed to create comments client:", err)
	}
	likesClient, err := commentsclient.NewGRPCLikesClient(commentsAddr)
	if err != nil {
		log.Fatal("failed to create likes client:", err)
	}
	notificationsClient, err := notificationsclient.NewGRPCNotificationsClient(os.Getenv("NOTIFICATIONS_SERVICE_ADDR"))
	if err != nil {
		log.Fatal("failed to create notifications client:", err)
	}
	postsClient, err := postsclient.NewGRPCPostsClient(os.Getenv("POSTS_SERVICE_ADDR"))
	if err != nil {
		log.Fatal("failed to create posts client:", err)
	}
	eventsClient, err := eventsclient.NewGRPCClient(os.Getenv("EVENTS_SERVICE_ADDR"))
	if err != nil {
		log.Fatal("failed to create events client:", err)
	}
	tagsClient, err := tagsclient.NewGRPCTagsClient(os.Getenv("TAGS_SERVICE_ADDR"))
	if err != nil {
		log.Fatal("failed to create tags client:", err)
	}

	ciRepo := repositories.NewCiRepository(db, searchService)
	colRepo := repositories.NewCollectionsRepository(db, searchService)
	platformRepo := repositories.NewPlatformsRepository(db)
	itemTypesRepo := repositories.NewItemTypesRepository(db)
	wlRepo := repositories.NewWLRepository(db)
	subRepo := repositories.NewSubscriptionRepository(db)
	paymentsRepo := repositories.NewPaymentsRepository(db)
	userRepo := repositories.NewUsersRepository(db, searchService, colRepo)
	err = db.AutoMigrate(&models.Users{}, &models.Platforms{}, &models.Collections{}, &models.CollectionItems{}, &models.ItemTypes{}, &models.WishListItems{}, &models.Subscription{}, &models.Payments{}, &models.Migrations{})
	if err != nil {
		return nil, err
	}

	host := os.Getenv("YANDEX_POST_HOST")
	port := os.Getenv("YANDEX_POST_PORT")
	portInt, _ := strconv.Atoi(port)
	user := os.Getenv("YANDEX_POST_USER")
	password := os.Getenv("YANDEX_POST_PASSWORD")

	jwtTtl, _ := strconv.Atoi(os.Getenv("JWT_EXPIRESS"))
	refreshTtl, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_EXPIRES"))

	jwtService := auth.NewJWTService(
		os.Getenv("JWT_SECRET"),
		time.Duration(jwtTtl)*time.Minute,
		time.Duration(refreshTtl)*time.Minute,
		*userRepo,
	)
	notificationsService := services.NewNotificationsService(notificationsClient, userRepo, colRepo, ciRepo, postsClient)
	eventsService := services.NewEventsService(userRepo, eventsClient, colRepo, ciRepo, postsClient, wlRepo, tagsClient)
	postsService := services.NewPostsService(postsClient, userRepo, eventsService, notificationsService, tagsClient)

	s3Service := s3.NewS3Service(redisClient)
	mailService := mail.NewMailService(host, portInt, user, password, `"Lootor" <noreply@lootor.me>`)
	userService := services.NewUserService(userRepo, jwtService, ciRepo, mailService, eventsService, notificationsService, colRepo, postsService, tagsClient)
	ciService := services.NewCiService(ciRepo, eventsService, colRepo, userRepo, platformRepo, s3Service, itemTypesRepo, notificationsService, tagsClient)
	colService := services.NewCollectionService(colRepo, userRepo, eventsService, ciRepo, s3Service, notificationsService, tagsClient)
	platformsService := services.NewPlatformsService(platformRepo)
	itemTypesService := services.NewItemTypesService(itemTypesRepo)
	tagsService := services.NewTagsService(tagsClient, userRepo, postsService, ciRepo, colRepo, eventsService, itemTypesRepo, searchService)
	fbService := feedback.NewFeedbackService()
	wlService := services.NewWLService(wlRepo, userRepo, ciRepo, eventsService)
	enrichedCIService := utils.NewEnrichedCIService(ciService, tagsClient)
	subService := services.NewSubscriptionService(subRepo)
	paymentService := payment.NewPayService(userRepo, userService, subService, paymentsRepo)
	captchaService := auth.NewRecaptchaService()
	reportsService := feedback.NewReportsService(fbService, ciRepo, userRepo, colRepo, wlRepo)
	feedService := services.NewFeedService(newsClient)
	commentsService := services.NewCommentsService(commentsClient, likesClient, userRepo, notificationsService, colRepo, ciRepo, postsService)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 60 * time.Second,
	}))

	routes.RegisterRoutes(e, jwtService, *userService)
	routes.TagsRouter(e, jwtService, *tagsService)
	routes.RegisterCollectionsRoutes(e, jwtService, *colService)
	routes.RegisterCollectionItemsRoutes(e, jwtService, *ciService, *enrichedCIService)
	routes.PlatformsRouter(e, jwtService, *platformsService)
	routes.S3Router(e, jwtService, *s3Service)
	routes.SearchRouter(e, jwtService, *searchService)
	routes.ReindexRouter(e, jwtService, *searchService, *userRepo, *ciRepo, *colRepo, *tagsService)
	routes.EventsRouter(e, jwtService, *eventsService)
	routes.FeedbackRouter(e, jwtService, *fbService)
	routes.ItemTypesRouter(e, jwtService, *itemTypesService)
	routes.WLRouter(e, jwtService, *wlService)
	routes.PaymentRouter(e, jwtService, *paymentService)
	routes.RecaptchaRouter(e, jwtService, *captchaService)
	routes.ReportsRouter(e, jwtService, *reportsService)
	routes.FeedRouter(e, jwtService, *feedService)
	routes.CommentsRouter(e, jwtService, *commentsService)
	routes.NotificationsRouter(e, jwtService, *notificationsService)
	routes.PostsRouter(e, jwtService, *postsService)

	controllers.NewReindexController(searchService, userRepo, colRepo, ciRepo, tagsService).ReindexInternal(context.Background())
	setupMonitoring(s3Service)

	if err = migrator.RunMigrations(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	c := cron.New()
	_, err = c.AddFunc("@midnight", func() { userService.CheckExpiredSubscriptions() })
	//_, err = c.AddFunc("@every 2m", func() { userService.CheckExpiredSubscriptions() })
	if err != nil {
		return nil, err
	}

	_, err = c.AddFunc("@midnight", func() { userService.CheckDeletedUsers() })
	if err != nil {
		return nil, err
	}

	c.Start()
	return &App{Echo: e}, nil
}
