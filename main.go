package main

import (
	"category-service/config"
	"category-service/config/key"
	deliveryG "category-service/internal/delivery/http"
	"category-service/internal/repository"
	"category-service/internal/usecase"
	"category-service/pkg/database"
	"category-service/pkg/logger"
	"category-service/pkg/middleware"
	sharedDomain "category-service/pkg/shared/domain"
	"category-service/pkg/token"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"category-service/internal/grpcservice"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.LoadConfig()

	logger := logger.NewLogger("category-service", logrus.DebugLevel, os.Stdout)

	_, cancel := context.WithTimeout(context.Background(), 1*time.Minute) // timeout to shutdown server and init configuration
	defer func() {
		cancel()
		if r := recover(); r != nil {
			// notify if this service failed to start (to logstash or slack)
			fmt.Println("Failed to start book service:", r)
			fmt.Printf("Stack trace: \n%s\n", debug.Stack())
		}
	}()

	db := &database.GormDatabase{}
	if err := db.Connect(cfg); err != nil {
		logger.Panic(fmt.Sprintf("Database connection error: %v", err), "db-error", "connection")
	}

	if err := db.AutoMigrate(
		&sharedDomain.Category{},
	); err != nil {
		logger.Panic(fmt.Sprintf("Failed to perform migration: %v", err), "migration", "error")
	}

	privateKeyPath := "config/key/private_key.pem"
	publicKeyPath := "config/key/public_key.pem"

	privateKey, publicKey, err := key.LoadRSAKeys(privateKeyPath, publicKeyPath)
	if err != nil {
		logger.Panic(fmt.Sprintf("Failed to load RSA keys: %v", err), "load rsa", "error")
	}

	jwtService := token.NewJWT(publicKey, privateKey)

	bookClient := grpcservice.NewBookGRPCClient(cfg.GetBookGRPCHost() + ":" + cfg.GetBookGRPCPort())

	// Setup repository, usecase, dan handler
	categoryRepo := repository.NewAuthorRepository(db.GetDB())
	categoryUsecase := usecase.NewAuthorUsecase(categoryRepo, bookClient)
	categoryHandler := deliveryG.NewCategoryHandler(categoryUsecase)

	// Setup routes
	httpServer := gin.Default()

	categoryRoutes := httpServer.Group("/categories", middleware.JWTAuthMiddleware(jwtService))
	{
		categoryRoutes.POST("/", categoryHandler.CreateCategory)
		categoryRoutes.GET("/", categoryHandler.GetAllCategories)
		categoryRoutes.GET("/:id", categoryHandler.GetCategoryByID)
		categoryRoutes.PATCH("/:id", categoryHandler.UpdateCategory)
		categoryRoutes.DELETE("/:id", categoryHandler.DeleteCategory)
	}

	httpPort := cfg.GetHTTPPort()
	if httpPort == "" {
		httpPort = "8080"
	}

	httpSrv := &http.Server{
		Addr:    ":" + httpPort,
		Handler: httpServer,
	}

	go func() {
		logger.Info("HTTP Server listening on port "+httpPort, "server_startup", "port:"+httpPort)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Panic(fmt.Sprintf("HTTP server error: %v", err), "http_server", "error")
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	logger.Info("Shutdown signal received, shutting down gracefully...", "", "")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		logger.Error(fmt.Sprintf("HTTP server shutdown error: %v", err), "", "")
	}

	logger.Info("Closing database connection...", "", "")
	db.Close()

	logger.Info("Servers shut down successfully", "", "")
}
