package main

import (
	"log"

	"github.com/aliwert/go-wolf"
	"github.com/aliwert/go-wolf/pkg/context"
	"github.com/aliwert/go-wolf/pkg/middleware"
)

func main() {
	// Create a new Wolf application
	app := wolf.New()

	// Add middleware
	app.Use(middleware.Logger())
	app.Use(middleware.Recovery())

	// Hello World route
	app.GET("/", func(c *context.Context) error {
		return c.JSON(200, wolf.Map{
			"message": "Hello, Wolf! 🐺",
			"version": "1.0.0",
		})
	})

	// API routes
	app.GET("/api/users/:id", func(c *context.Context) error {
		id := c.Param("id")
		return c.JSON(200, wolf.Map{
			"user_id": id,
			"name":    "Wert",
		})
	})

	// Static route with wildcard
	app.GET("/static/*filepath", func(c *context.Context) error {
		filepath := c.Param("filepath")
		return c.JSON(200, wolf.Map{
			"filepath": filepath,
		})
	})

	// Start the server
	log.Println("Starting server on :8080")
	if err := app.Run(":8080"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
