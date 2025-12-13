package helpers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetUserIDFromContext mengambil user_id dari c.Locals
func GetUserIDFromContext(c *fiber.Ctx) (string, error) {
	userIDLocal := c.Locals("user_id")
	
	if userIDLocal == nil {
		return "", fmt.Errorf("unauthorized: User ID tidak ditemukan di context")
	}

	switch v := userIDLocal.(type) {
	case string:
		return v, nil
	case uuid.UUID:
		return v.String(), nil
	default:
		return "", fmt.Errorf("internal error: Tipe data User ID tidak dikenali")
	}
}