package routes

import (
	"database/sql"
	"uas/app/repository"
	"uas/app/services"
	"uas/middleware"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(app *fiber.App, postgreSQL *sql.DB, mongoDB *mongo.Database) {

	api := app.Group("/api/v1") // (tidak perlu login)

	// Autentikasi & Otorisasi 
	auth := api.Group("/auth")
	auth.Post("/login", services.Login)
	auth.Post("/refresh", services.Refresh)
	auth.Get("/profile", middleware.AuthRequired(), services.GetProfile)

	// Protected routes (perlu login) 
	protected := api.Group("", middleware.AuthRequired()) 
	
	// Users (Admin)
	userService := services.NewUserService(postgreSQL)
	protected.Post("/users", middleware.RequirePermission("users:create"), userService.CreateUser)
	protected.Get("/users", middleware.RequirePermission("users:read"), userService.GetAllUsers)
	protected.Get("/users/:id", middleware.RequirePermission("users:read"), userService.GetUserByID)
	protected.Put("/users/:id", middleware.RequirePermission("users:update"), userService.UpdateUser)
	protected.Delete("/users/:id", middleware.RequirePermission("users:delete"), userService.DeleteUser)
	protected.Put("/users/:id/role", userService.UpdateUserRole)

	// Students (Admin)
	studentRepo := repository.NewStudentRepository(postgreSQL)
	studentService := services.NewStudentService(studentRepo)
	protected.Get("/students", middleware.RequirePermission("students:read"), studentService.GetStudents)
	protected.Get("/students/:id", middleware.RequirePermission("students:read"), studentService.GetStudentByID)
}