package helpers

import (
	"context"
	"fmt"
	"uas/app/repository"
)

// ValidateAdvisorAccess mengecek apakah user adalah Dosen Wali yang sah untuk prestasi tersebut
func ValidateAdvisorAccess(ctx context.Context, repo repository.AchievementRepository, achievementID string, userID string) error {
	
	lecturerID, err := repo.GetLecturerIDByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("akses ditolak: akun Anda tidak terdaftar sebagai Dosen Wali")
	}

	ach, err := repo.GetAchievementByID(ctx, achievementID)
	if err != nil {
		return fmt.Errorf("data prestasi tidak ditemukan")
	}

	// Cek Status (Harus Submitted)
	if ach.Status != "submitted" {
		return fmt.Errorf("gagal memproses: hanya prestasi berstatus 'submitted' yang bisa diverifikasi. Status saat ini: %s", ach.Status)
	}

	// Cek Hubungan Dosen Wali - Mahasiswa
	isAdvisor, err := repo.CheckStudentAdvisorRelationship(ctx, lecturerID, ach.StudentID)
	if err != nil {
		return fmt.Errorf("terjadi kesalahan saat memvalidasi data perwalian")
	}

	if !isAdvisor {
		return fmt.Errorf("akses ditolak: Anda bukan Dosen Wali dari mahasiswa yang mengajukan prestasi ini")
	}

	return nil
}