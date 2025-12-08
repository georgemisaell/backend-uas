package services

import (
	"database/sql"
	"uas/app/models"
	"uas/app/repository"
	"uas/database"
	"uas/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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
            return c.Status(401).JSON(fiber.Map{"error": "Username salah"})
        }
        return c.Status(500).JSON(fiber.Map{"error": "Terjadi kesalahan pada server"})
    }

    if !utils.CheckPassword(req.Password, user.PasswordHash) {
        return c.Status(401).JSON(fiber.Map{"error": "password salah"})
    }

    if !user.IsActive {
        return c.Status(403).JSON(fiber.Map{"error": "Akun anda dinonaktifkan. Silahkan hubungi admin."})
    }

    accessToken, err := utils.GenerateToken(user)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal generate token"})
    }

    refreshToken, _ := utils.GenerateRefreshToken(user)

    userResponse := models.UserResponseDTO{
        ID:       user.ID,
        Username: user.Username,
        FullName: user.FullName,
        Role:     user.RoleName,
    }

    return c.Status(200).JSON(fiber.Map{
        "status": "success",
        "data": models.LoginResponse{
            Token:        accessToken,
            RefreshToken: refreshToken,
            User:         userResponse,
        },
    })
}

func Refresh(c *fiber.Ctx) error {
    var req struct {
        RefreshToken string `json:"refreshToken"`
    }

    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }

    // Parse token
    token, err := jwt.Parse(req.RefreshToken, func(t *jwt.Token) (interface{}, error) {
        return utils.JwtSecret, nil
    })

    if err != nil || !token.Valid {
        return c.Status(401).JSON(fiber.Map{"error": "Invalid refresh token"})
    }

    claims := token.Claims.(jwt.MapClaims)

    // Pastikan refresh token
    if claims["type"] != "refresh" {
        return c.Status(401).JSON(fiber.Map{"error": "Invalid token type"})
    }

    // Ambil userId dari JWT claim
    userIDStr := claims["userId"].(string)

    // Convert ke UUID
    userUUID, err := uuid.Parse(userIDStr)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
    }

    // Ambil DB instance (tanpa modif repository)
    db := database.ConnectDB()

    // Ambil user lengkap dari database
    user, err := repository.GetUserByID(db, userUUID)
    if err != nil {
        return c.Status(401).JSON(fiber.Map{"error": "User not found"})
    }

    // Generate access token baru
    newAccessToken, err := utils.GenerateToken(user)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to generate access token"})
    }

    return c.JSON(fiber.Map{
        "status": "success",
        "token":  newAccessToken,
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