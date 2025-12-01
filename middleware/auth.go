package middleware

import (
	"strings"
	"uas/app/repository"
	"uas/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil token dari header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Token akses diperlukan",
			})
		}

		// Extract token dari "Bearer TOKEN"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Format token tidak valid",
			})
		}

		// Validasi token
		claims, err := utils.ValidateToken(tokenParts[1])
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "Token tidak valid atau expired",
			})
		}

		// Simpan informasi user di context
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("role_name", claims.RoleName)

		return c.Next()
	}
}

// Menerima parameter string 'perm' (misal: "achievement:create")
func RequirePermission(perm string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // 1. Ambil User ID dari Locals (yang diset oleh AuthRequired)
        userIDLocal := c.Locals("user_id")
        
        if userIDLocal == nil {
            return c.Status(401).JSON(fiber.Map{"error": "Unauthorized: User ID not found"})
        }

        // Pastikan tipe datanya UUID
        userID, ok := userIDLocal.(uuid.UUID)
        if !ok {
             return c.Status(500).JSON(fiber.Map{"error": "Internal Server Error: Invalid User ID format"})
        }

        // 2. Cek ke Database via Repository
        allowed, err := repository.CheckPermission(userID, perm)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Gagal memverifikasi izin"})
        }

        // 3. Logika Allow/Deny
        if !allowed {
            return c.Status(403).JSON(fiber.Map{
                "error": "Forbidden: Anda tidak memiliki izin '" + perm + "'",
            })
        }

        // 4. Lanjut ke Controller
        return c.Next()
    }
}