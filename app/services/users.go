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
    var req models.CreateUserRequest

    // 1. Parsing Request
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Format data tidak valid",
            "success": false,
        })
    }

    // 2. Hash Password
    hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Gagal mengenkripsi password",
            "success": false,
        })
    }

    // 3. Mulai Transaksi
    tx, err := db.Begin()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Gagal memulai transaksi database",
            "success": false,
        })
    }
    
    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()

    userID := uuid.New() 

    parsedRoleID, err := uuid.Parse(req.RoleID)
    if err != nil {
         return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Format Role ID salah",
            "success": false,
        })
    }

    newUser := models.User{
        ID:           userID,
        Username:     req.Username,
        Email:        req.Email,
        PasswordHash: string(hashedPwd),
        FullName:     req.FullName,
        RoleID:       parsedRoleID,
        RoleName:     req.RoleName,
        IsActive:     true,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }

    err = repository.CreateUser(tx, newUser)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Gagal menyimpan data user",
            "success": false,
            "error":   err.Error(),
        })
    }

    if req.RoleName == "Mahasiswa" && req.Student != nil {
        newStudent := models.Student{
            ID:           uuid.New(),
            UserID:       userID,
            StudentID:    req.Student.StudentID,
            ProgramStudy: req.Student.ProgramStudy,
            AcademicYear: req.Student.AcademicYear,
            AdvisorID: req.Student.AdvisorID,
            CreatedAt:    time.Now(),
        }

        if err = repository.CreateStudent(tx, newStudent); err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "message": "Gagal menyimpan data mahasiswa",
                "success": false,
                "error":   err.Error(),
            })
        }

    } else if req.RoleName == "Dosen Wali" && req.Lecture != nil {
        newLecture := models.Lecture{
            ID:         uuid.New(),
            UserID:     userID,
            LectureID:  req.Lecture.LectureID,
            Department: req.Lecture.Department,
            CreatedAt:  time.Now(),
        }

        if err = repository.CreateLecture(tx, newLecture); err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "message": "Gagal menyimpan data dosen",
                "success": false,
                "error":   err.Error(),
            })
        }
    }

    if err = tx.Commit(); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Gagal commit transaksi",
            "success": false,
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