package main

import (
	"log"

	"os"

	"github.com/TFX0019/api-go-gds/features/auth"
	"github.com/TFX0019/api-go-gds/features/customers"
	"github.com/TFX0019/api-go-gds/features/dashboard"
	"github.com/TFX0019/api-go-gds/features/materials"
	"github.com/TFX0019/api-go-gds/features/plans"
	"github.com/TFX0019/api-go-gds/features/products"
	"github.com/TFX0019/api-go-gds/features/subscriptions"
	"github.com/TFX0019/api-go-gds/features/tasks"
	"github.com/TFX0019/api-go-gds/features/user"
	"github.com/TFX0019/api-go-gds/features/wallets"
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
	// Migrate Auth models
	if err := database.DB.AutoMigrate(&auth.User{}, &auth.VerificationCode{}, &auth.Role{}, &auth.Session{}, &customers.Customer{}, &products.Product{}, &materials.Material{}, &tasks.Task{}, &wallets.Wallet{}, &wallets.CreditTransaction{}, &subscriptions.Subscription{}, &plans.Plan{}); err != nil {
		log.Fatal("Migration failed: ", err)
	}

	// Seed Roles
	var roles = []string{"admin", "member"}
	for _, roleName := range roles {
		var role auth.Role
		if err := database.DB.FirstOrCreate(&role, auth.Role{Name: roleName}).Error; err != nil {
			log.Printf("Failed to seed role %s: %v", roleName, err)
		}
	}

	// Seed Free Tier Plan
	var freePlan plans.Plan
	if err := database.DB.Where("product_id = ?", "free_tier").First(&freePlan).Error; err != nil {
		freePlan = plans.Plan{
			ProductID:    "free_tier",
			Title:        "Free Tier",
			Description:  "Starter plan for new users",
			Price:        0,
			Benefits:     []string{"20 Customer Limit", "20 Product Limit", "20 Material Limit", "20 Task Limit"},
			MaxCustomers: 20,
			MaxProducts:  20,
			MaxMaterials: 20,
			MaxTasks:     20,
		}
		if err := database.DB.Create(&freePlan).Error; err != nil {
			log.Printf("Failed to seed free plan: %v", err)
		} else {
			log.Println("Seeded 'free_tier' plan")
		}
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
	// Plans Feature
	plansRepo := plans.NewRepository(database.DB)

	// Auth Feature
	authRepo := auth.NewRepository(database.DB)
	authService := auth.NewService(authRepo, plansRepo)
	authController := auth.NewController(authService)
	auth.RegisterRoutes(app, authController)

	// Customers Feature
	customersRepo := customers.NewRepository(database.DB)
	customersService := customers.NewService(customersRepo, authRepo, plansRepo)
	customersController := customers.NewController(customersService)
	customers.RegisterRoutes(app, customersController)

	// Products Feature
	productsRepo := products.NewRepository(database.DB)
	productsService := products.NewService(productsRepo, authRepo, plansRepo)
	productsController := products.NewController(productsService)
	products.RegisterRoutes(app, productsController)

	// Materials Feature
	materialsRepo := materials.NewRepository(database.DB)
	materialsService := materials.NewService(materialsRepo)
	materialsController := materials.NewController(materialsService)
	materials.RegisterRoutes(app, materialsController)

	// Tasks Feature
	tasksRepo := tasks.NewRepository(database.DB)
	tasksService := tasks.NewService(tasksRepo)
	tasksController := tasks.NewController(tasksService)
	tasks.RegisterRoutes(app, tasksController)

	// Plans Feature Routes
	plans.RegisterRoutes(app, database.DB)

	// User Feature
	user.RegisterRoutes(app, database.DB)

	// Subscriptions Feature
	subscriptionsRepo := subscriptions.NewRepository(database.DB)
	subscriptionsService := subscriptions.NewService(subscriptionsRepo)
	subscriptionsController := subscriptions.NewController(subscriptionsService)
	subscriptions.RegisterRoutes(app, subscriptionsController)

	// Dashboard Feature
	dashboardRepo := dashboard.NewRepository(database.DB)
	dashboardService := dashboard.NewService(dashboardRepo)
	dashboardController := dashboard.NewController(dashboardService)
	dashboard.RegisterRoutes(app, dashboardController)

	// 6. Start Server
	port := config.GetEnv("PORT", "3000")
	log.Printf("Server running on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
