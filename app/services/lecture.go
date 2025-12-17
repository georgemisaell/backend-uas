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

// GetLecturers godoc
// @Summary      Ambil Semua Dosen Wali
// @Description  Mengambil daftar lengkap dosen wali yang aktif.
// @Tags         Lecturers
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200  {object}  map[string][]models.GetLecture
// @Failure      500  {object}  map[string]string
// @Router       /lecturers [get]
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

// GetLecturerByID godoc
// @Summary      Ambil Detail Dosen
// @Description  Mendapatkan data detail satu dosen berdasarkan ID (UUID).
// @Tags         Lecturers
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id   path      string  true  "Lecturer ID (UUID)"
// @Success      200  {object}  map[string]models.GetLecture
// @Failure      400  {object}  map[string]string "Format ID salah"
// @Failure      404  {object}  map[string]string "Tidak ditemukan"
// @Failure      500  {object}  map[string]string
// @Router       /lecturers/{id} [get]
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

// GetLecturerAdvisees godoc
// @Summary      Ambil Mahasiswa Bimbingan
// @Description  Melihat daftar mahasiswa yang dibimbing oleh dosen tertentu.
// @Tags         Lecturers
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id   path      string  true  "Lecturer ID (UUID)"
// @Success      200  {object}  map[string][]models.GetStudent
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /lecturers/{id}/advisees [get]
func (s *lecturerService) GetLecturerAdvisees(c *fiber.Ctx) error {
	lecturerID := c.Params("id")

	if _, err := uuid.Parse(lecturerID); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Format ID Dosen tidak valid",
			"success": false,
		})
	}

	students, err := s.repo.GetAdviseesByLecturerID(c.Context(), lecturerID)
	
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi kesalahan server saat mengambil data bimbingan",
			"success": false,
			"error":   err.Error(),
		})
	}

	if len(students) == 0 {
		return c.JSON(fiber.Map{
			"message": "Dosen ini belum memiliki mahasiswa bimbingan",
			"success": true,
			"data":    []string{},
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data mahasiswa bimbingan berhasil diambil",
		"success": true,
		"data":    students,
	})
}