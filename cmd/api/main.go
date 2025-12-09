package main

import (
	"log"

	"os"

	"github.com/TFX0019/api-go-gds/features/auth"
	"github.com/TFX0019/api-go-gds/features/customers"
	"github.com/TFX0019/api-go-gds/features/products"
	"github.com/TFX0019/api-go-gds/pkg/config"
	"github.com/TFX0019/api-go-gds/pkg/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// 1. Config
	config.LoadConfig()

	// 2. Database
	database.ConnectDB()

	// 3. Migrations
	// Migrate Auth models
	if err := database.DB.AutoMigrate(&auth.User{}, &customers.Customer{}, &products.Product{}); err != nil {
		log.Fatal("Migration failed: ", err)
	}

	// Ensure uploads directory exists
	if err := os.MkdirAll("./uploads", 0755); err != nil {
		log.Fatal("Failed to create uploads directory: ", err)
	}

	// 4. Fiber App
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())
	app.Static("/uploads", "./uploads")

	// 5. Setup Features
	// Auth Feature
	authRepo := auth.NewRepository(database.DB)
	authService := auth.NewService(authRepo)
	authController := auth.NewController(authService)
	auth.RegisterRoutes(app, authController)

	// Customers Feature
	customersRepo := customers.NewRepository(database.DB)
	customersService := customers.NewService(customersRepo)
	customersController := customers.NewController(customersService)
	customers.RegisterRoutes(app, customersController)

	// Products Feature
	productsRepo := products.NewRepository(database.DB)
	productsService := products.NewService(productsRepo)
	productsController := products.NewController(productsService)
	products.RegisterRoutes(app, productsController)

	// 6. Start Server
	port := config.GetEnv("PORT", "3000")
	log.Printf("Server running on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
