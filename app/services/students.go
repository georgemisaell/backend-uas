package services

import (
	"database/sql"
	"strings"
	"uas/app/models"
	"uas/app/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type StudentService interface {
	GetStudents(c *fiber.Ctx) error
	GetStudentByID(c *fiber.Ctx) error
	UpdateStudentAdvisor(c *fiber.Ctx) error
}

type studentService struct {
	repo repository.StudentRepository
}

func NewStudentService(repo repository.StudentRepository) StudentService {
	return &studentService{repo: repo}
}

// GetStudents godoc
// @Summary      Ambil Semua Data Mahasiswa
// @Description  Mengambil daftar lengkap mahasiswa (terfilter role 'Mahasiswa'). Admin Only.
// @Tags         Students
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200  {object}  map[string][]models.GetStudent
// @Failure      500  {object}  map[string]string
// @Router       /students [get]
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

// GetStudentByID godoc
// @Summary      Ambil Detail Mahasiswa
// @Description  Mendapatkan data detail satu mahasiswa berdasarkan ID (UUID).
// @Tags         Students
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id   path      string  true  "Student ID (UUID)"
// @Success      200  {object}  map[string]models.GetStudent
// @Failure      400  {object}  map[string]string "Format ID salah"
// @Failure      404  {object}  map[string]string "Tidak ditemukan"
// @Failure      500  {object}  map[string]string
// @Router       /students/{id} [get]
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

// UpdateStudentAdvisor godoc
// @Summary      Update Dosen Wali
// @Description  Mengganti dosen wali (Advisor) untuk mahasiswa tertentu.
// @Tags         Students
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id       path      string                       true  "Student ID (UUID)"
// @Param        request  body      models.UpdateAdvisorRequest  true  "Advisor ID Baru"
// @Success      200      {object}  map[string]string
// @Failure      400      {object}  map[string]string "ID salah / Dosen tidak ada"
// @Failure      404      {object}  map[string]string "Mahasiswa tidak ditemukan"
// @Failure      500      {object}  map[string]string
// @Router       /students/{id}/advisor [put]
func (s *studentService) UpdateStudentAdvisor(c *fiber.Ctx) error {
	
	studentID := c.Params("id")
	if _, err := uuid.Parse(studentID); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Format Student ID tidak valid",
			"success": false,
		})
	}

	var req models.UpdateAdvisorRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Format data JSON tidak valid",
			"success": false,
		})
	}

	if _, err := uuid.Parse(req.AdvisorID); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Format Advisor ID tidak valid",
			"success": false,
		})
	}

	err := s.repo.UpdateStudentAdvisor(c.Context(), studentID, req.AdvisorID)
	if err != nil {
		if err.Error() == "student_not_found" {
			return c.Status(404).JSON(fiber.Map{
				"message": "Data mahasiswa tidak ditemukan",
				"success": false,
			})
		}
		
		if strings.Contains(err.Error(), "foreign key") {
			return c.Status(400).JSON(fiber.Map{
				"message": "ID Dosen Wali tidak ditemukan di database",
				"success": false,
			})
		}

		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal mengganti dosen wali",
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Dosen Wali berhasil diperbarui",
		"success": true,
	})
}