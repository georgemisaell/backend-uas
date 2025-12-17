package main

import (
	"log"
	"uas/config"
	"uas/database"
	"uas/routes"

	_ "uas/docs"

	"github.com/gofiber/fiber/v2"
)

// @title           UAS API Documentation
// @version         1.0
// @description     Dokumentasi lengkap API untuk Sistem Prestasi Mahasiswa.
// @termsOfService  http://swagger.io/terms/

// @contact.name    George Misael
// @contact.email   georgemisaelgantume@gmail.com

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @host            localhost:3000
// @BasePath        /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Masukkan token dengan format "Bearer <token_jwt_disini>"

func main() {

	// Menghubungkan ENV
	config.Config();

	// Database postgre	SQL
	postgreSQL := database.ConnectDB()
	mongoDB := database.ConnectMongoDB()

	// Inisialisasi fiber
	app := fiber.New(fiber.Config{
		ErrorHandler: func (c *fiber.Ctx, err error) error {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// routes
	routes.SetupRoutes(app, postgreSQL, mongoDB)

	// Server
	log.Fatal(app.Listen(":3000"))
}