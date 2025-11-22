package main

import (
	"log"
	"uas/config"
	"uas/database"
	"uas/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {

	// Menghubungkan ENV
	config.Config();

	// Database postgre	SQL
	postgreSQL := database.ConnectDB()

	// Inisialisasi fiber
	app := fiber.New(fiber.Config{
		ErrorHandler: func (c *fiber.Ctx, err error) error {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// routes
	routes.SetupRoutes(app, postgreSQL)

	// Server
	log.Fatal(app.Listen(":3000"))
}