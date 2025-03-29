package main

import (
	"category-service/config"
	"category-service/config/key"
	"category-service/internal/delivery/http"
	"category-service/internal/repository"
	"category-service/internal/usecase"
	"category-service/pkg/database"
	"category-service/pkg/middleware"
	sharedDomain "category-service/pkg/shared/domain"
	"category-service/pkg/token"
	"fmt"
	"log"

	"category-service/internal/grpcservice"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	db := &database.GormDatabase{}
	if err := db.Connect(cfg); err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	if err := db.AutoMigrate(
		&sharedDomain.Category{},
	); err != nil {
		log.Fatalf("Failed to perform migration: %v", err)
	}

	privateKeyPath := "config/key/private_key.pem"
	publicKeyPath := "config/key/public_key.pem"

	privateKey, publicKey, err := key.LoadRSAKeys(privateKeyPath, publicKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	jwtService := token.NewJWT(publicKey, privateKey)

	bookClient := grpcservice.NewBookGRPCClient(cfg.GetBookGRPCHost() + ":" + cfg.GetBookGRPCPort())
	fmt.Println("bookClient " + cfg.GetBookGRPCHost() + ":" + cfg.GetBookGRPCPort())

	// Setup repository, usecase, dan handler
	categoryRepo := repository.NewAuthorRepository(db.GetDB())
	categoryUsecase := usecase.NewAuthorUsecase(categoryRepo, bookClient)
	categoryHandler := http.NewCategoryHandler(categoryUsecase)

	// Setup routes
	r := gin.Default()
	r.Use()

	categoryRoutes := r.Group("/categories", middleware.JWTAuthMiddleware(jwtService))
	{
		categoryRoutes.POST("/", categoryHandler.CreateCategory)
		categoryRoutes.GET("/", categoryHandler.GetAllCategories)
		categoryRoutes.GET("/:id", categoryHandler.GetCategoryByID)
		categoryRoutes.PATCH("/:id", categoryHandler.UpdateCategory)
		categoryRoutes.DELETE("/:id", categoryHandler.DeleteCategory)
	}

	port := cfg.GetHTTPPort()
	if port == "" {
		port = "8080"
	}
	log.Println("HTTP Server listening on port ", port)
	r.Run(":" + port)
}
