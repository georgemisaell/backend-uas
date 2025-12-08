package routes

import (
	"database/sql"
	"uas/app/services"
	"uas/middleware"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(app *fiber.App, postgreSQL *sql.DB, mongoDB *mongo.Database) {

	api := app.Group("/api/v1")

	// Autentikasi & Otorisasi 
	// (tidak perlu login)
	auth := api.Group("/auth")
	auth.Post("/login", services.Login)
	auth.Post("/refresh", services.Refresh)
	auth.Get("/profile", middleware.AuthRequired(), services.GetProfile)

	// Protected routes (perlu login) 
	protected := api.Group("", middleware.AuthRequired()) 

	// Users (Admin)
	protected.Get("/users", middleware.RequirePermission("user:manage"), func(c *fiber.Ctx) error {
		return services.GetAllUsers(c, postgreSQL)
	})

	protected.Get("/users/:id", middleware.RequirePermission("user:manage"), func(c *fiber.Ctx) error {
		return services.GetUserByID(c, postgreSQL)
	})

	protected.Post("/users", middleware.RequirePermission("user:manage"), func(c *fiber.Ctx) error {
		return services.CreateUser(c, postgreSQL)
	})

	protected.Put("/users/:id", middleware.RequirePermission("user:manage"), func(c *fiber.Ctx) error {
		return services.UpdateUser(c, postgreSQL)
	})

	protected.Delete("/users/:id", middleware.RequirePermission("user:manage"), func(c *fiber.Ctx) error {
			return services.DeleteUser(c, postgreSQL)
	})

	protected.Put("/users/:id/role", middleware.RequirePermission("user:manage"), func(c *fiber.Ctx) error {
			return services.UpdateUserRole(c, postgreSQL)
	})
}