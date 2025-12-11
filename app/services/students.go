package services

import (
	"database/sql"
	"uas/app/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type StudentService interface {
	GetStudents(c *fiber.Ctx) error
	GetStudentByID(c *fiber.Ctx) error
}

type studentService struct {
	repo repository.StudentRepository
}

func NewStudentService(repo repository.StudentRepository) StudentService {
	return &studentService{repo: repo}
}

func (s *studentService) GetStudents(c *fiber.Ctx) error {
	const targetRole = "Mahasiswa"
	
	students, err := s.repo.GetAllStudentsByRole(c.Context(), targetRole)
	
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi kesalahan pada server saat mengambil data mahasiswa",
			"success": false,
            "error":   err.Error(),
		})
	}

	if len(students) == 0 {
		return c.JSON(fiber.Map{
			"message": "Data Mahasiswa tidak ditemukan",
			"success": true,
			"data":    []string{},
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data Mahasiswa berhasil diambil",
		"success": true,
		"data":    students,
	})
}

func (s *studentService) GetStudentByID(c *fiber.Ctx) error {
	idParam := c.Params("id")

	if _, err := uuid.Parse(idParam); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Format ID tidak valid",
			"success": false,
		})
	}

	student, err := s.repo.GetStudentByID(c.Context(), idParam)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{
				"message": "Data mahasiswa tidak ditemukan",
				"success": false,
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi kesalahan pada server",
			"success": false,
            "error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data mahasiswa ditemukan",
		"success": true,
		"data":    student,
	})
}