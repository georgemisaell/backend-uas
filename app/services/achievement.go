package services

import (
	"fmt"
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
	GetAllAchievements(c *fiber.Ctx) error
	GetAchievementDetail(c *fiber.Ctx) error
    GetAchievementHistory(c *fiber.Ctx) error
    UploadAttachment(c *fiber.Ctx) error
}

type achievementService struct {
	repo repository.AchievementRepository
}

func NewAchievementService(repo repository.AchievementRepository) AchievementService {
	return &achievementService{repo: repo}
}

// CreateAchievement godoc
// @Summary      Buat Prestasi Baru (Draft)
// @Description  Mahasiswa membuat data prestasi baru. Status awal otomatis 'draft'.
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request body models.CreateAchievementRequest true "Data Prestasi"
// @Success      201  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /achievements [post]
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

// UpdateAchievement godoc
// @Summary      Edit Data Prestasi
// @Description  Mengubah data prestasi. Hanya bisa dilakukan jika status masih 'draft'.
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id      path string true "Achievement ID (UUID)"
// @Param        request body models.CreateAchievementRequest true "Data Update"
// @Success      200  {object} map[string]string
// @Failure      400  {object} map[string]string
// @Failure      403  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /achievements/{id} [put]
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

// DeleteAchievement godoc
// @Summary      Hapus Prestasi
// @Description  Menghapus data prestasi (Soft Delete). Hanya bisa dilakukan jika status masih 'draft'.
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id   path string true "Achievement ID (UUID)"
// @Success      200  {object} map[string]string
// @Failure      400  {object} map[string]string
// @Failure      403  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /achievements/{id} [delete]
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

// SubmitAchievement godoc
// @Summary      Ajukan Prestasi (Submit)
// @Description  Mengubah status prestasi dari 'draft' menjadi 'submitted'. Menunggu verifikasi dosen.
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id   path string true "Achievement ID (UUID)"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      403  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /achievements/{id}/submit [post]
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

// VerifyAchievement godoc
// @Summary      Verifikasi Prestasi (Dosen Wali)
// @Description  Menyetujui prestasi mahasiswa bimbingan. Status berubah menjadi 'verified'.
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id   path string true "Achievement ID (UUID)"
// @Success      200  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      403  {object} map[string]string "Bukan mahasiswa bimbingan anda"
// @Failure      500  {object} map[string]string
// @Router       /achievements/{id}/verify [post]
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

// RejectAchievement godoc
// @Summary      Tolak Prestasi (Dosen Wali)
// @Description  Menolak prestasi mahasiswa bimbingan dengan catatan. Status berubah menjadi 'rejected'.
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id      path string true "Achievement ID (UUID)"
// @Param        request body models.RejectAchievementRequest true "Alasan Penolakan"
// @Success      200  {object} map[string]string
// @Failure      400  {object} map[string]string "Catatan wajib diisi"
// @Failure      403  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /achievements/{id}/reject [post]
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

// GetAllAchievements godoc
// @Summary      List Semua Prestasi
// @Description  Mengambil daftar prestasi. Mahasiswa melihat miliknya sendiri, Dosen melihat anak walinya, Admin lihat semua.
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200  {object} map[string][]models.AchievementResponse
// @Failure      500  {object} map[string]string
// @Router       /achievements [get]
func (s *achievementService) GetAllAchievements(c *fiber.Ctx) error {
  
    // 1. Ambil Info User
    roleName := c.Locals("role_name").(string)
    userIDLocal := c.Locals("user_id") // ID User Login

    // Variable untuk filter query
    var filterUserID string = ""

    if roleName == "Mahasiswa" {
        switch v := userIDLocal.(type) {
        case string:
            filterUserID = v
        case uuid.UUID:
            filterUserID = v.String()
        }
    }

    pgRefs, names, nims, err := s.repo.GetAllReferences(c.Context(), filterUserID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"message": "Gagal mengambil data referensi"})
    }

    if len(pgRefs) == 0 {
        return c.JSON(fiber.Map{"success": true, "data": []string{}})
    }

    var mongoIDs []string
    for _, ref := range pgRefs {
        mongoIDs = append(mongoIDs, ref.MongoAchievementID)
    }

    mongoDocs, err := s.repo.GetMongoDetailsByIDs(c.Context(), mongoIDs)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"message": "Gagal mengambil detail prestasi"})
    }

    var responses []models.AchievementResponse
    for _, ref := range pgRefs {

        detail, ok := mongoDocs[ref.MongoAchievementID]
        
        res := models.AchievementResponse{
            ID:          ref.ID,
            MongoID:     ref.MongoAchievementID,
            StudentID:   ref.StudentID,
            StudentName: names[ref.ID],
            StudentNIM:  nims[ref.ID],
            Status:      ref.Status,
            CreatedAt:   ref.CreatedAt,
        }

        if ok {
            res.Title = detail.Title
            res.AchievementType = detail.AchievementType
        } else {
            res.Title = "[Detail Tidak Ditemukan]"
        }

        responses = append(responses, res)
    }

    return c.JSON(fiber.Map{
        "success": true,
        "data":    responses,
    })
}

// GetAchievementDetail godoc
// @Summary      Detail Prestasi Lengkap
// @Description  Mengambil detail lengkap prestasi (Gabungan data Postgres & MongoDB).
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id   path string true "Achievement ID (UUID)"
// @Success      200  {object} map[string]models.AchievementResponse
// @Failure      403  {object} map[string]string "Akses Ditolak"
// @Failure      404  {object} map[string]string "Tidak Ditemukan"
// @Failure      500  {object} map[string]string
// @Router       /achievements/{id} [get]
func (s *achievementService) GetAchievementDetail(c *fiber.Ctx) error {
    id := c.Params("id")

    refData, err := s.repo.GetAchievementReferenceWithDetail(c.Context(), id)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"message": "Data prestasi tidak ditemukan"})
    }

    userIDLocal := c.Locals("user_id")
    roleName := c.Locals("role_name").(string)

    if roleName == "Mahasiswa" && userIDLocal != nil {
        currentUserID := ""
        switch v := userIDLocal.(type) {
        case string: currentUserID = v
        case uuid.UUID: currentUserID = v.String()
        }

        myStudentID, _ := s.repo.GetStudentIDByUserID(c.Context(), currentUserID)
        if myStudentID != refData.StudentID {
            return c.Status(403).JSON(fiber.Map{"message": "Anda tidak memiliki akses ke detail prestasi ini"})
        }
    }

    mongoData, err := s.repo.GetMongoDetailByID(c.Context(), refData.MongoID)
    if err != nil {
        refData.Title = "[Detail Hilang]"
    } else {
        // Merge Data
        refData.AchievementType = mongoData.AchievementType
        refData.Title = mongoData.Title
        refData.Description = mongoData.Description
        refData.Details = mongoData.Details
        refData.Tags = mongoData.Tags
        refData.Points = mongoData.Points
    }

    return c.JSON(fiber.Map{
        "success": true,
        "data":    refData,
    })
}

// GetAchievementHistory godoc
// @Summary      Riwayat Prestasi
// @Description  Melihat log perubahan status prestasi (Draft -> Submitted -> Verified/Rejected).
// @Tags         Achievements
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id   path string true "Achievement ID (UUID)"
// @Success      200  {object} map[string][]models.HistoryItem
// @Failure      403  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Router       /achievements/{id}/history [get]
func (s *achievementService) GetAchievementHistory(c *fiber.Ctx) error {
    id := c.Params("id")

    data, err := s.repo.GetAchievementReferenceWithDetail(c.Context(), id)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"message": "Prestasi tidak ditemukan"})
    }

    // VALIDASI AKSES
    userIDLocal := c.Locals("user_id")
    roleName := c.Locals("role_name").(string)
    
    if roleName == "Mahasiswa" && userIDLocal != nil {
        currentUserID := ""
        switch v := userIDLocal.(type) {
        case string: currentUserID = v
        case uuid.UUID: currentUserID = v.String()
        }

        myStudentID, _ := s.repo.GetStudentIDByUserID(c.Context(), currentUserID)
        if myStudentID != data.StudentID {
            return c.Status(403).JSON(fiber.Map{"message": "Anda tidak memiliki akses ke detail prestasi ini"})
        }
    }

    var histories []models.HistoryItem

    // DRAFT (Pasti ada karena ada created_at)
    histories = append(histories, models.HistoryItem{
        Action:    "Draft Dibuat",
        Timestamp: data.CreatedAt,
        Actor:     data.StudentName,
        Note:      "Prestasi disimpan sebagai draft",
    })

    // SUBMITTED (Cek apakah submitted_at ada isinya?)
    if data.SubmittedAt != nil {
        histories = append(histories, models.HistoryItem{
            Action:    "Disubmit",
            Timestamp: *data.SubmittedAt, // Dereference pointer
            Actor:     data.StudentName,
            Note:      "Menunggu verifikasi Dosen Wali",
        })
    }

    // VERIFIED / REJECTED (Cek status akhir & verified_at)
    if data.VerifiedAt != nil {
        action := ""
        actor := "Dosen Wali"
        note := ""

      switch data.Status {
        case "verified":
            action = "Diverifikasi"
            note = "Prestasi valid dan disetujui"
        case "rejected":
            action = "Ditolak"
            note = data.RejectionNote
        }

        if action != "" {
            histories = append(histories, models.HistoryItem{
                Action:    action,
                Timestamp: *data.VerifiedAt,
                Actor:     actor, 
                Note:      note,
            })
        }
    }

    return c.JSON(fiber.Map{
        "success": true,
        "data":    histories,
    })
}

// UploadAttachment godoc
// @Summary      Upload Bukti Prestasi (File)
// @Description  Mengunggah file bukti (Sertifikat, Foto, dll) ke prestasi. Maksimal status 'draft'.
// @Tags         Achievements
// @Accept       multipart/form-data
// @Produce      json
// @Security     Bearer
// @Param        id   path string true "Achievement ID (UUID)"
// @Param        file formData file true "File Dokumen"
// @Success      200  {object} map[string]models.Attachment
// @Failure      400  {object} map[string]string "File missing / Status salah"
// @Failure      403  {object} map[string]string "Forbidden"
// @Failure      500  {object} map[string]string
// @Router       /achievements/{id}/attachments [post]
func (s *achievementService) UploadAttachment(c *fiber.Ctx) error {
    id := c.Params("id")

    userID, err := helpers.GetUserIDFromContext(c)
    if err != nil {
        return c.Status(401).JSON(fiber.Map{"message": err.Error()})
    }
    
    data, err := s.repo.GetAchievementByID(c.Context(), id)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"message": "Prestasi tidak ditemukan"})
    }

    studentID, _ := s.repo.GetStudentIDByUserID(c.Context(), userID)
    if data.StudentID != studentID {
        return c.Status(403).JSON(fiber.Map{"message": "Anda tidak berhak upload file ke prestasi ini"})
    }

    if data.Status != "draft" {
        return c.Status(400).JSON(fiber.Map{"message": "Upload gagal: Prestasi sudah disubmit/diverifikasi"})
    }

    file, err := c.FormFile("file")
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"message": "File tidak ditemukan. Gunakan key form-data 'file'"})
    }

    uniqueFileName := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
    filePath := fmt.Sprintf("./uploads/%s", uniqueFileName)

    if err := c.SaveFile(file, filePath); err != nil {
        return c.Status(500).JSON(fiber.Map{"message": "Gagal menyimpan file ke server"})
    }

    attachment := models.Attachment{
        FileName:   file.Filename,
        FileURL:    "/uploads/" + uniqueFileName,
        FileType:   file.Header.Get("Content-Type"),
        UploadedAt: time.Now(),
    }

    err = s.repo.AddAttachmentToMongo(c.Context(), data.MongoAchievementID, attachment)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"message": "Gagal mencatat file ke database"})
    }

    return c.JSON(fiber.Map{
        "success": true,
        "message": "File berhasil diupload",
        "data":    attachment,
    })
}