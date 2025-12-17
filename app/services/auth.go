package services

import (
	"database/sql"
	"uas/app/models"
	"uas/app/repository"
	"uas/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService interface {
	Login(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error
	GetProfile(c *fiber.Ctx) error
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Login(c *fiber.Ctx) error {
	var req models.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Request body tidak valid"})
	}

	if req.Username == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Username dan password harus diisi"})
	}

	user, err := s.userRepo.GetByUsernameOrEmail(c.Context(), req.Username)
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

func (s *authService) Refresh(c *fiber.Ctx) error {
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

	if claims["type"] != "refresh" {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token type"})
	}

	userIDStr := claims["userId"].(string)

	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	user, err := s.userRepo.GetUserByID(c.Context(), userUUID)
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

func (s *authService) GetProfile(c *fiber.Ctx) error {
	userIDLocal := c.Locals("user_id")
	if userIDLocal == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}
	userID := userIDLocal.(uuid.UUID)

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