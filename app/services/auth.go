package services

import (
	"database/sql"
	"uas/app/models"
	"uas/app/repository"
	"uas/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Login(c *fiber.Ctx) error {
    var req models.LoginRequest

    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Request body tidak valid"})
    }

    if req.Username == "" || req.Password == "" {
        return c.Status(400).JSON(fiber.Map{"error": "Username dan password harus diisi"})
    }

    user, err := repository.Login(req.Username)
    if err != nil {
        if err == sql.ErrNoRows {
            return c.Status(401).JSON(fiber.Map{"error": "Username atau password salah"})
        }
        return c.Status(500).JSON(fiber.Map{"error": "Terjadi kesalahan pada server"})
    }

    if !utils.CheckPassword(req.Password, user.PasswordHash) {
        return c.Status(401).JSON(fiber.Map{"error": "Username atau password salah"})
    }

    if !user.IsActive {
        return c.Status(403).JSON(fiber.Map{"error": "Akun anda dinonaktifkan. Silahkan hubungi admin."})
    }

    token, err := utils.GenerateToken(user)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal generate token"})
    }

    return c.JSON(fiber.Map{
        "success": true,
        "message": "Login berhasil",
        "data": models.LoginResponse{
            User:  user,
            Token: token,
        },
    })
}

func GetProfile(c *fiber.Ctx) error { 
    userID := c.Locals("user_id").(uuid.UUID) 
    username := c.Locals("username").(string) 
    role := c.Locals("role_name").(string) 
 
    return c.JSON(fiber.Map{ 
        "success": true, 
        "message": "Profile berhasil diambil", 
        "data": fiber.Map{ 
            "user_id":  userID, 
            "username": username, 
            "role":     role, 
        }, 
    }) 
}

func RefreshToken(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user session"})
	}
	
	username, ok := c.Locals("username").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user session"})
	}

	roleName, ok := c.Locals("role_name").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user session"})
	}

	user := models.User{
		ID:       userID,
		Username: username,
		RoleName: roleName,
	}

	newToken, err := utils.GenerateToken(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal merefresh token"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Token berhasil diperbarui",
		"data": fiber.Map{
			"token": newToken,
		},
	})
}