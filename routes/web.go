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
	protected.Put("/students/:id/advisor", middleware.RequirePermission("students:update"), studentService.UpdateStudentAdvisor)

	// Lectures (Admin)
	lecturerRepo := repository.NewLecturerRepository(postgreSQL)
  lecturerService := services.NewLecturerService(lecturerRepo)
	protected.Get("/lecturers", middleware.RequirePermission("lecturers:read"), lecturerService.GetLecturers)
	protected.Get("/lecturers/:id/advisees", middleware.RequirePermission("lecturers:read"), lecturerService.GetLecturerAdvisees)

	// Achievements (Mahasiswa)
	achRepo := repository.NewAchievementRepository(postgreSQL, mongoDB)
	achService := services.NewAchievementService(achRepo)
	protected.Post("/achievements", middleware.RequirePermission("achievements:create"), achService.CreateAchievement)
	protected.Put("/achievements/:id", middleware.RequirePermission("achievements:update"), achService.UpdateAchievement)
	protected.Delete("/achievements/:id", middleware.RequirePermission("achievements:delete"), achService.DeleteAchievement)
	protected.Post("/achievements/:id/submit", middleware.RequirePermission("achievements:update"), achService.SubmitAchievement)
	protected.Post("/achievements/:id/attachments", middleware.RequirePermission("achievements:update"), achService.UploadAttachment)

	// Achievements (Dosen Wali)
	protected.Post("/achievements/:id/verify", middleware.RequirePermission("achievements:verify"), achService.VerifyAchievement)
	protected.Post("/achievements/:id/reject", middleware.RequirePermission("achievements:reject"), achService.RejectAchievement)

	// Achievements (Admin)
	protected.Get("/achievements/:id", middleware.RequirePermission("achievements:read"), achService.GetAchievementDetail)
	protected.Get("/achievements/:id/history", middleware.RequirePermission("achievements:read"), achService.GetAchievementHistory)
	
	// Achievements (All Role)
	protected.Get("/achievements", middleware.RequirePermission("achievements:read"), achService.GetAllAchievements)

	// Reports & Analitycs 
	reportRepo := repository.NewReportRepository(postgreSQL)
	reportService := services.NewReportService(reportRepo, achRepo)
	reports := protected.Group("/reports")
	reports.Get("/statistics", middleware.RequirePermission("reports:read"), reportService.GetSystemStatistics)
	reports.Get("/student/:id", middleware.RequirePermission("reports:read"), reportService.GetStudentReport)
}