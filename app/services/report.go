package services

import (
	"uas/app/models"
	"uas/app/repository"
	"uas/helpers"

	"github.com/gofiber/fiber/v2"
)

type ReportService interface {
	GetSystemStatistics(c *fiber.Ctx) error
	GetStudentReport(c *fiber.Ctx) error
}

type reportService struct {
	reportRepo      repository.ReportRepository
	achievementRepo repository.AchievementRepository
}

func NewReportService(reportRepo repository.ReportRepository, achievementRepo repository.AchievementRepository) ReportService {
	return &reportService{
		reportRepo:      reportRepo,
		achievementRepo: achievementRepo,
	}
}

// GetSystemStatistics godoc
// @Summary      Dashboard Statistik
// @Description  Menampilkan ringkasan jumlah prestasi berdasarkan status (Draft, Verified, Rejected). Admin & Dosen Wali Only.
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200  {object}  models.DashboardStatistics
// @Failure      403  {object}  map[string]string "Akses Ditolak"
// @Failure      500  {object}  map[string]string
// @Router       /reports/statistics [get]
func (s *reportService) GetSystemStatistics(c *fiber.Ctx) error {
	roleName := c.Locals("role_name").(string)

	if roleName != "Admin" && roleName != "Dosen Wali" {
		return c.Status(403).JSON(fiber.Map{"message": "Akses ditolak"})
	}

	stats, err := s.reportRepo.GetStatistics(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal mengambil statistik"})
	}

	return c.JSON(fiber.Map{"success": true, "data": stats})
}

// GetStudentReport godoc
// @Summary      Rapor Prestasi Mahasiswa (Transkrip)
// @Description  Menampilkan profil, total poin (SKP), dan daftar prestasi verified mahasiswa. Mahasiswa hanya bisa lihat punya sendiri. Dosen Wali hanya anak bimbingan.
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id   path      string  true  "Student ID (UUID)"
// @Success      200  {object}  models.StudentReportResponse
// @Failure      403  {object}  map[string]string "Bukan hak akses anda"
// @Failure      404  {object}  map[string]string "Mahasiswa tidak ditemukan"
// @Failure      500  {object}  map[string]string
// @Router       /reports/student/{id} [get]
func (s *reportService) GetStudentReport(c *fiber.Ctx) error {
	targetStudentID := c.Params("id")
	requesterID, _ := helpers.GetUserIDFromContext(c)
	requesterRole := c.Locals("role_name").(string)

	if requesterRole == "Mahasiswa" {
		myStudentID, _ := s.achievementRepo.GetStudentIDByUserID(c.Context(), requesterID)
		if myStudentID != targetStudentID {
			return c.Status(403).JSON(fiber.Map{"message": "Akses ditolak"})
		}
	} else if requesterRole == "Dosen Wali" {
		lecturerID, _ := s.achievementRepo.GetLecturerIDByUserID(c.Context(), requesterID)
		isAdvisor, _ := s.achievementRepo.CheckStudentAdvisorRelationship(c.Context(), lecturerID, targetStudentID)
		if !isAdvisor {
			return c.Status(403).JSON(fiber.Map{"message": "Bukan mahasiswa bimbingan Anda"})
		}
	}

	profile, err := s.reportRepo.GetStudentProfile(c.Context(), targetStudentID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Mahasiswa tidak ditemukan"})
    }

	refs, err := s.reportRepo.GetVerifiedAchievementsByStudentID(c.Context(), targetStudentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal mengambil data prestasi"})
	}

	var mongoIDs []string
	for _, ref := range refs {
		mongoIDs = append(mongoIDs, ref.MongoAchievementID)
	}
	mongoDocs, _ := s.achievementRepo.GetMongoDetailsByIDs(c.Context(), mongoIDs)

	var totalPoints int
	var achievementList []models.AchievementResponse

	for _, ref := range refs {
		detail, ok := mongoDocs[ref.MongoAchievementID]
		item := models.AchievementResponse{
			ID:        ref.ID,
			Status:    ref.Status,
			CreatedAt: ref.CreatedAt,
            Title:     "[Data Corrupt]",
		}
		if ok {
			item.Title = detail.Title
			item.AchievementType = detail.AchievementType
			item.Points = detail.Points
			totalPoints += detail.Points
		}
		achievementList = append(achievementList, item)
	}

	profile.TotalPoints = totalPoints
	profile.TotalItems = len(achievementList)

	return c.JSON(fiber.Map{
		"success": true,
		"data": models.StudentReportResponse{
			Profile:      profile,
			Achievements: achievementList,
		},
	})
}