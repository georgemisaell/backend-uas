package services

import (
	"database/sql"
	"uas/app/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LecturerService interface {
	GetLecturers(c *fiber.Ctx) error
	GetLecturerByID(c *fiber.Ctx) error
	GetLecturerAdvisees(c *fiber.Ctx) error
}

type lecturerService struct {
	repo repository.LecturerRepository
}

func NewLecturerService(repo repository.LecturerRepository) LecturerService {
	return &lecturerService{repo: repo}
}

func (s *lecturerService) GetLecturers(c *fiber.Ctx) error {
	const targetRole = "Dosen Wali"

	lecturers, err := s.repo.GetAllLecturersByRole(c.Context(), targetRole)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi kesalahan server",
			"success": false,
			"error":   err.Error(),
		})
	}

	if len(lecturers) == 0 {
		return c.JSON(fiber.Map{
			"message": "Data Dosen Wali tidak ditemukan",
			"success": true,
			"data":    []string{},
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data Dosen Wali berhasil diambil",
		"success": true,
		"data":    lecturers,
	})
}

func (s *lecturerService) GetLecturerByID(c *fiber.Ctx) error {
	idParam := c.Params("id")

	if _, err := uuid.Parse(idParam); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Format ID tidak valid",
			"success": false,
		})
	}

	lecturer, err := s.repo.GetLecturerByID(c.Context(), idParam)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{
				"message": "Dosen tidak ditemukan",
				"success": false,
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi kesalahan server",
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data dosen ditemukan",
		"success": true,
		"data":    lecturer,
	})
}

func (s *lecturerService) GetLecturerAdvisees(c *fiber.Ctx) error {
	// 1. Ambil ID Dosen dari parameter URL
	lecturerID := c.Params("id")

	// 2. Validasi UUID
	if _, err := uuid.Parse(lecturerID); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Format ID Dosen tidak valid",
			"success": false,
		})
	}

	// 3. Panggil Repository (Opsional: Cek dulu apakah Dosennya ada, tapi langsung query juga oke)
	students, err := s.repo.GetAdviseesByLecturerID(c.Context(), lecturerID)
	
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi kesalahan server saat mengambil data bimbingan",
			"success": false,
			"error":   err.Error(),
		})
	}

	// 4. Handle jika data kosong (belum punya anak bimbingan)
	if len(students) == 0 {
		return c.JSON(fiber.Map{
			"message": "Dosen ini belum memiliki mahasiswa bimbingan",
			"success": true,
			"data":    []string{},
		})
	}

	// 5. Return Data
	return c.JSON(fiber.Map{
		"message": "Data mahasiswa bimbingan berhasil diambil",
		"success": true,
		"data":    students,
	})
}