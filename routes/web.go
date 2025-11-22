package routes

import (
	"database/sql"
	"uas/app/services"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, postgreSQL *sql.DB) {

	api := app.Group("/api/v1")

	// Autentikasi & Otorisasi

	// Users (Admin)
	api.Get("/users", func(c *fiber.Ctx) error {
		return services.GetAllUsers(c, postgreSQL)
	})

	api.Get("/users/:id", func(c *fiber.Ctx) error {
		return services.GetUserByID(c, postgreSQL)
	})

	api.Post("/users", func(c *fiber.Ctx) error {
		return services.CreateUser(c, postgreSQL)
	})

	api.Put("/users/:id", func(c *fiber.Ctx) error {
		return services.UpdateUser(c, postgreSQL)
	})

	api.Delete("/users/:id", func(c *fiber.Ctx) error {
			return services.DeleteUser(c, postgreSQL)
	})

	// Achievements

	// Mahasiswa & Dosen

	 // Reports & Analytics 
}