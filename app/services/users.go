package services

import (
	"database/sql"
	"time"
	"uas/app/models"
	"uas/app/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GetAllUsers(c *fiber.Ctx, db *sql.DB) error {
	users, err := repository.GetAllUsers(db)
	
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Terjadi kesalahan pada server",
			"success": false,
		})
	}

	if len(users) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Data User tidak ditemukan",
			"success": true,
			"data":    []string{},
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data berhasil diambil",
		"success": true,
		"data":    users,
	})
}

func GetUserByID(c *fiber.Ctx, db *sql.DB) error {
    idParam := c.Params("id")

    userID, err := uuid.Parse(idParam)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Format ID tidak valid",
            "success": false,
        })
    }

    user, err := repository.GetUserByID(db, userID)
    
    if err == sql.ErrNoRows {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "message": "User tidak ditemukan",
            "success": false,
        })
    } else if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Terjadi kesalahan server",
            "success": false,
        })
    }

    return c.JSON(fiber.Map{
        "message": "Data user ditemukan",
        "success": true,
        "data":    user,
    })
}

func CreateUser(c *fiber.Ctx, db *sql.DB) error {
	var req models.CreateUser

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Format data tidak valid",
			"success": false,
		})
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengenkripsi password",
			"success": false,
		})
	}

	newUser := models.User{
		ID:           uuid.New(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPwd),
		FullName:     req.FullName,
		RoleID:       req.RoleID,
		IsActive:     req.IsActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = repository.CreateUser(db, newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menyimpan data user",
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User berhasil dibuat",
		"success": true,
		"data":    newUser,
	})
}

func UpdateUser(c *fiber.Ctx, db *sql.DB) error {
	idParam := c.Params("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Format ID tidak valid",
			"success": false,
		})
	}

	var user models.UpdateUser
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Format data JSON tidak valid",
			"success": false,
		})
	}

	err = repository.UpdateUser(db, userID, user)
	
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User tidak ditemukan, gagal update",
			"success": false,
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengupdate data user",
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "User berhasil diupdate",
		"success": true,
		"data": user,
	})
}

func DeleteUser(c *fiber.Ctx, db *sql.DB) error {
	idParam := c.Params("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Format ID tidak valid",
			"success": false,
		})
	}

	err = repository.DeleteUser(db, userID)

	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User tidak ditemukan, gagal menghapus",
			"success": false,
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Terjadi kesalahan server saat menghapus data",
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "User berhasil dihapus",
		"success": true,
	})
}