package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"

	"github.com/takechiyo-19940627/medicalquest/config"
	handlers "github.com/takechiyo-19940627/medicalquest/handler"
	"github.com/takechiyo-19940627/medicalquest/handler/validation"
	"github.com/takechiyo-19940627/medicalquest/infrastructure/database"
)

func main() {
	// Initialize configuration
	cfg := config.New()

	// Setup database connection
	client, err := database.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer client.Close()

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.Validator = validation.NewCustomValidator(validator.New())

	// Register routes
	handlers.RegisterRoutes(e, client)

	// Start server
	go func() {
		if err := e.Start(fmt.Sprintf(":%s", cfg.ServerPort)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
