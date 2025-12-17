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

type UserService interface {
	GetAllUsers(c *fiber.Ctx) error
	GetUserByID(c *fiber.Ctx) error
	CreateUser(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
	UpdateUserRole(c *fiber.Ctx) error
}

type userService struct {
	db           *sql.DB                      
	userRepo     repository.UserRepository    
	studentRepo  repository.StudentRepository 
	lecturerRepo repository.LecturerRepository
}


func NewUserService(
	db *sql.DB,
	userRepo repository.UserRepository,
	studentRepo repository.StudentRepository,
	lecturerRepo repository.LecturerRepository,
) UserService {
	return &userService{
		db:           db,
		userRepo:     userRepo,
		studentRepo:  studentRepo,
		lecturerRepo: lecturerRepo,
	}
}

// GetAllUsers godoc
// @Summary      Ambil Semua User
// @Description  Mengambil daftar semua user yang terdaftar (Admin Only)
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200  {array}   models.User
// @Failure      500  {object}  map[string]string
// @Router       /users [get]
func (s *userService) GetAllUsers(c *fiber.Ctx) error {
	users, err := s.userRepo.GetAllUsers(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi kesalahan pada server",
			"success": false,
		})
	}

	if len(users) == 0 {
		return c.JSON(fiber.Map{
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

// GetUserByID godoc
// @Summary      Ambil User berdasarkan ID
// @Description  Mengambil detail data satu user berdasarkan UUID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id   path      string  true  "User ID (UUID)"
// @Success      200  {object}  models.User
// @Failure      400  {object}  map[string]string "Format ID salah"
// @Failure      404  {object}  map[string]string "User tidak ditemukan"
// @Failure      500  {object}  map[string]string
// @Router       /users/{id} [get]
func (s *userService) GetUserByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Format ID tidak valid",
			"success": false,
		})
	}

	user, err := s.userRepo.GetUserByID(c.Context(), userID)
	if err == sql.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{
			"message": "User tidak ditemukan",
			"success": false,
		})
	} else if err != nil {
		return c.Status(500).JSON(fiber.Map{
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

// CreateUser godoc
// @Summary      Tambah User Baru (Admin)
// @Description  Membuat user baru beserta data Mahasiswa atau Dosen (Transactional)
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request body models.CreateUserRequest true "Data User Lengkap"
// @Success      201  {object} models.User
// @Failure      400  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /users [post]
func (s *userService) CreateUser(c *fiber.Ctx) error {
	var req models.CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Format data tidak valid",
			"success": false,
		})
	}

	// Hash password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal mengenkripsi password",
			"success": false,
		})
	}

	// 1. MULAI TRANSAKSI
	tx, err := s.db.Begin()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal memulai transaksi database",
			"success": false,
		})
	}
	defer tx.Rollback()

	userID := uuid.New()
	roleID, err := uuid.Parse(req.RoleID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
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
		RoleID:       roleID,
		RoleName:     req.RoleName,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// 2. INSERT USER
	err = s.userRepo.CreateUser(c.Context(), tx, newUser)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal menyimpan data user",
			"success": false,
			"error":   err.Error(),
		})
	}

	// 3. INSERT STUDENT
	if req.RoleName == "Mahasiswa" && req.Student != nil {
		newStudent := models.Student{
			ID:           uuid.New(),
			UserID:       userID,
			StudentID:    req.Student.StudentID,
			ProgramStudy: req.Student.ProgramStudy,
			AcademicYear: req.Student.AcademicYear,
			AdvisorID:    req.Student.AdvisorID,
			CreatedAt:    time.Now(),
		}

		if err := s.studentRepo.CreateStudent(c.Context(), tx, newStudent); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "Gagal menyimpan data mahasiswa",
				"success": false,
				"error":   err.Error(),
			})
		}
	}

	if req.RoleName == "Dosen Wali" && req.Lecture != nil {
		newLecture := models.Lecture{
			ID:         uuid.New(),
			UserID:     userID,
			LectureID:  req.Lecture.LectureID,
			Department: req.Lecture.Department,
			CreatedAt:  time.Now(),
		}
				
		if err := s.lecturerRepo.CreateLecture(c.Context(), tx, newLecture); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "Gagal menyimpan data dosen",
				"success": false,
				"error":   err.Error(),
			})
		}
	}

	// 5. COMMIT TRANSAKSI
	if err := tx.Commit(); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal commit transaksi",
			"success": false,
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "User berhasil dibuat",
		"success": true,
		"data":    newUser,
	})
}

func (s *userService) UpdateUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Format ID tidak valid",
			"success": false,
		})
	}

	var user models.UpdateUser
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Format data JSON tidak valid",
			"success": false,
		})
	}

	err = s.userRepo.UpdateUser(c.Context(), userID, user)
	if err == sql.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{
			"message": "User tidak ditemukan",
			"success": false,
		})
	} else if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal mengupdate data user",
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "User berhasil diupdate",
		"success": true,
		"data":    user,
	})
}

// DeleteUser godoc
// @Summary      Hapus User
// @Description  Menghapus data user secara permanen (atau Soft Delete tergantung implementasi repo)
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id   path      string  true  "User ID (UUID)"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users/{id} [delete]
func (s *userService) DeleteUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Format ID tidak valid",
			"success": false,
		})
	}

	err = s.userRepo.DeleteUser(c.Context(), userID)
	if err == sql.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{
			"message": "User tidak ditemukan",
			"success": false,
		})
	} else if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal menghapus data user",
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "User berhasil dihapus",
		"success": true,
	})
}

// UpdateUser godoc
// @Summary      Update Data User
// @Description  Memperbarui data profile user (Username, Email, Fullname, dll)
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id       path      string             true  "User ID (UUID)"
// @Param        request  body      models.UpdateUser  true  "Data yang diupdate"
// @Success      200      {object}  models.User
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /users/{id} [put]
func (s *userService) UpdateUserRole(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Format User ID tidak valid",
			"success": false,
		})
	}

	var req models.UpdateRole
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Format JSON tidak valid",
			"success": false,
		})
	}

	roleID, err := uuid.Parse(req.RoleID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Format Role ID tidak valid",
			"success": false,
		})
	}

	err = s.userRepo.UpdateUserRole(c.Context(), userID, roleID)
	if err == sql.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{
			"message": "User tidak ditemukan",
			"success": false,
		})
	} else if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal update role user",
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Role user berhasil diperbarui",
		"success": true,
	})
}