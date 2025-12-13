package services

import (
	"time"
	"uas/app/models"
	"uas/app/repository"
	"uas/helpers"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementService interface {
	CreateAchievement(c *fiber.Ctx) error
	UpdateAchievement(c *fiber.Ctx) error
  DeleteAchievement(c *fiber.Ctx) error
	SubmitAchievement(c *fiber.Ctx) error
	VerifyAchievement(c *fiber.Ctx) error
	RejectAchievement(c *fiber.Ctx) error
}

type achievementService struct {
	repo repository.AchievementRepository
}

func NewAchievementService(repo repository.AchievementRepository) AchievementService {
	return &achievementService{repo: repo}
}

func (s *achievementService) CreateAchievement(c *fiber.Ctx) error {

	var req models.CreateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Format data tidak valid",
			"success": false,
		})
	}

	userIDLocal := c.Locals("user_id")
	if userIDLocal == nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized: User ID tidak ditemukan",
			"success": false,
		})
	}

	var userID string
	switch v := userIDLocal.(type) {
	case string:
		userID = v
	case uuid.UUID:
		userID = v.String()
	default:
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error: Tipe data User ID tidak dikenali",
			"success": false,
		})
	}

	studentID, err := s.repo.GetStudentIDByUserID(c.Context(), userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Data mahasiswa tidak ditemukan untuk user ini",
			"success": false,
		})
	}

	mongoData := models.AchievementMongo{
		ID:              primitive.NewObjectID(),
		StudentID:       studentID,
		AchievementType: req.AchievementType,
		Title:           req.Title,
		Description:     req.Description,
		Details:         req.Details,
		Tags:            req.Tags,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	mongoID, err := s.repo.CreateAchievementMongo(c.Context(), mongoData)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal menyimpan data detail prestasi",
			"success": false,
			"error":   err.Error(),
		})
	}

	pgRef := models.AchievementReference{
		ID:                 uuid.New().String(),
		StudentID:          studentID,
		MongoAchievementID: mongoID,
		Status:             "draft",
	}

	err = s.repo.CreateAchievementReference(c.Context(), pgRef)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal menyimpan referensi prestasi",
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Prestasi berhasil dibuat (Draft)",
		"success": true,
		"data": fiber.Map{
			"id":                   pgRef.ID,
			"mongo_achievement_id": mongoID,
			"status":               "draft",
			"created_at":           time.Now(),
		},
	})
}

func (s *achievementService) UpdateAchievement(c *fiber.Ctx) error {
    id := c.Params("id")

    var req models.CreateAchievementRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"message": "Format data salah", "success": false})
    }

    userIDLocal := c.Locals("user_id")
    if userIDLocal == nil {
        return c.Status(401).JSON(fiber.Map{"message": "Unauthorized"})
    }

    var userID string
    switch v := userIDLocal.(type) {
    case string: userID = v
    case uuid.UUID: userID = v.String()
    }

    studentID, err := s.repo.GetStudentIDByUserID(c.Context(), userID)
    if err != nil {
        return c.Status(403).JSON(fiber.Map{"message": "User bukan mahasiswa"})
    }

    existingData, err := s.repo.GetAchievementByID(c.Context(), id)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"message": "Prestasi tidak ditemukan"})
    }

    if existingData.StudentID != studentID {
        return c.Status(403).JSON(fiber.Map{"message": "Anda tidak berhak mengedit data ini"})
    }

    if existingData.Status != "draft" {
        return c.Status(400).JSON(fiber.Map{
            "message": "Gagal update: Hanya status 'draft' yang boleh diedit",
            "current_status": existingData.Status,
        })
    }

    mongoData := models.AchievementMongo{
        AchievementType: req.AchievementType,
        Title:           req.Title,
        Description:     req.Description,
        Details:         req.Details,
        Tags:            req.Tags,
    }

    err = s.repo.UpdateAchievement(c.Context(), existingData.ID, existingData.MongoAchievementID, mongoData)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"message": "Gagal mengupdate data"})
    }

    return c.JSON(fiber.Map{"message": "Prestasi berhasil diupdate", "success": true})
}

func (s *achievementService) DeleteAchievement(c *fiber.Ctx) error {
    id := c.Params("id")

    userIDLocal := c.Locals("user_id")
    var userID string
    switch v := userIDLocal.(type) {
    case string: userID = v
    case uuid.UUID: userID = v.String()
    }

    studentID, err := s.repo.GetStudentIDByUserID(c.Context(), userID)
    if err != nil {
        return c.Status(403).JSON(fiber.Map{"message": "User bukan mahasiswa"})
    }

    existingData, err := s.repo.GetAchievementByID(c.Context(), id)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"message": "Prestasi tidak ditemukan"})
    }

    if existingData.StudentID != studentID {
        return c.Status(403).JSON(fiber.Map{"message": "Anda tidak berhak menghapus data ini"})
    }

    if existingData.Status != "draft" {
        return c.Status(400).JSON(fiber.Map{
            "message": "Gagal hapus: Hanya status 'draft' yang boleh dihapus",
        })
    }

    err = s.repo.SoftDeleteAchievement(c.Context(), existingData.ID, existingData.MongoAchievementID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"message": "Gagal menghapus data"})
    }

    return c.JSON(fiber.Map{"message": "Prestasi berhasil dihapus", "success": true})
}

func (s *achievementService) SubmitAchievement(c *fiber.Ctx) error {
    id := c.Params("id")

    // 1. Ambil User ID & Student ID (Standard Auth Check)
    userIDLocal := c.Locals("user_id")
    var userID string
    switch v := userIDLocal.(type) {
    case string: userID = v
    case uuid.UUID: userID = v.String()
    }

    studentID, err := s.repo.GetStudentIDByUserID(c.Context(), userID)
    if err != nil {
        return c.Status(403).JSON(fiber.Map{"message": "User bukan mahasiswa"})
    }

    // 2. Cek Data Existing
    achievement, err := s.repo.GetAchievementByID(c.Context(), id)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"message": "Prestasi tidak ditemukan"})
    }

    // 3. Validasi Kepemilikan
    if achievement.StudentID != studentID {
        return c.Status(403).JSON(fiber.Map{"message": "Anda tidak berhak mensubmit data ini"})
    }

    if achievement.Status != "draft" {
        return c.Status(400).JSON(fiber.Map{
            "message": "Gagal submit: Hanya prestasi berstatus 'draft' yang bisa disubmit",
            "current_status": achievement.Status,
        })
    }

    // 5. Lakukan Submit
    err = s.repo.SubmitAchievement(c.Context(), id)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"message": "Gagal melakukan submit prestasi"})
    }

    return c.JSON(fiber.Map{
        "success": true,
        "message": "Prestasi berhasil disubmit dan menunggu verifikasi",
        "data": fiber.Map{
            "id": id,
            "status": "submitted",
            "submitted_at": time.Now(),
        },
    })
}

func (s *achievementService) VerifyAchievement(c *fiber.Ctx) error {
	achievementID := c.Params("id")

	verifierUserID, err := helpers.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"message": err.Error()})
	}

	if err := helpers.ValidateAdvisorAccess(c.Context(), s.repo, achievementID, verifierUserID); err != nil {
		return c.Status(403).JSON(fiber.Map{"message": err.Error()})
	}

	err = s.repo.VerifyAchievement(c.Context(), achievementID, verifierUserID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal memverifikasi prestasi"})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Prestasi berhasil diverifikasi"})
}

func (s *achievementService) RejectAchievement(c *fiber.Ctx) error {
	achievementID := c.Params("id")

	var req models.RejectAchievementRequest
	if err := c.BodyParser(&req); err != nil || req.RejectionNote == "" {
		return c.Status(400).JSON(fiber.Map{"message": "Alasan penolakan (rejection_note) wajib diisi"})
	}

	verifierUserID, err := helpers.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"message": err.Error()})
	}

	if err := helpers.ValidateAdvisorAccess(c.Context(), s.repo, achievementID, verifierUserID); err != nil {
		return c.Status(403).JSON(fiber.Map{"message": err.Error()})
	}

	err = s.repo.RejectAchievement(c.Context(), achievementID, verifierUserID, req.RejectionNote)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal menolak prestasi"})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Prestasi berhasil ditolak"})
}