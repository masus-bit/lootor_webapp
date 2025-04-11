package main

import (
	"context"
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
		log.Fatalf("Config error: %v", err)
	}

	application, err := app.NewEchoApp(cfg)
	if err != nil {
		log.Fatalf("Application init error: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := application.Echo.Start(cfg.Server.Address); err != nil && err != http.ErrServerClosed {
			application.Echo.Logger.Fatal("Shutting down the server")
		}
	}()

	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := application.Echo.Shutdown(ctx); err != nil {
		application.Echo.Logger.Fatal(err)
	}
}
