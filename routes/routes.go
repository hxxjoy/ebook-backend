package routes

import (
	v1 "book-fiber/routes/v1"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	
	// API Version
	router := app.Group("/api/v1")
	v1.RegisterRoutes(router)

}

