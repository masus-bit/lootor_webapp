package main

import (
	"log"
	"lootor/internal/app"
	"lootor/internal/config"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Something went wrong", err)
	}

	application, e := app.NewChiApp(cfg)
	if e != nil {
		log.Fatalf("роктер не инициализирован")
	}

	server := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      application.Router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("Server is running on %s", cfg.Server.Address)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}

}
