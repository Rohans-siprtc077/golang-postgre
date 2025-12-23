package main

import (
	"log"

	"golang-postgre/config"
	"golang-postgre/models"
	"golang-postgre/routes"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load environment variables
	godotenv.Load()

	// Connect to DB
	config.ConnectDB()
	config.DB.AutoMigrate(&models.User{})

	// Create Echo instance
	e := echo.New()

	// ✅ Enable CORS (REQUIRED for jQuery / frontend)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.PUT,
			echo.DELETE,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
		},
	}))

	// Register routes
	routes.RegisterRoutes(e)

	// Start server
	log.Println("🚀 Server running at http://localhost:8080")
	log.Fatal(e.Start(":8080"))
}
