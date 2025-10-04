package main

import (
	config "aprilpollo/internal/adapter/config"
	httpfiber "aprilpollo/internal/adapter/handler/fiber"
	"aprilpollo/internal/adapter/handler/fiber/middleware"
	"aprilpollo/internal/adapter/handler/fiber/routes"
	gormOrm "aprilpollo/internal/adapter/storage/gorm"
	"aprilpollo/internal/adapter/storage/gorm/repository"
	"aprilpollo/internal/core/service"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func main() {
	config, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✔ [INFO] Loading Configuration")

	db, err := gormOrm.NewGormDB(config.Database, &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	fmt.Println("✔ [INFO] Database Connection")

	// Initialize repositories
	authRepo := repository.NewAuthRepository(db.GetDB())

	// Initialize services
	authService := service.NewAuthService(authRepo)

	// Initialize handlers
	authHandler := routes.NewAuthHandler(authService)

	// Initialize Fiber app
	fiberApp := fiber.New()

	// Initialize middleware
	middlewareApp := middleware.NewMiddlewareHandler(fiberApp, config.App, config.JWT)
	middlewareApp.SetupGlobalMiddleware()

	// Initialize App routes
	app := httpfiber.NewApp(fiberApp, middlewareApp)
	app.MainRoutes()
	app.AuthRoutes(authHandler)
	fmt.Printf("✔ [INFO] Starting %s Version %s on %s Mode\n", config.App.AppName, config.App.AppVersion, func() string {
		if config.App.Development {
			return "Development"
		}
		return "Production"
	}())
	app.Serve(fmt.Sprintf(":%s", config.App.ApiPort))

}
