package repository

import (
	"context"
	"database/sql"
	"uas/app/models"
)

type ReportRepository interface {
	GetStatistics(ctx context.Context) (models.DashboardStatistics, error)
	GetStudentProfile(ctx context.Context, studentID string) (models.StudentReportProfile, error)
	GetVerifiedAchievementsByStudentID(ctx context.Context, studentID string) ([]models.AchievementReference, error)
}

type reportRepository struct {
	pg *sql.DB
}

func NewReportRepository(pg *sql.DB) ReportRepository {
	return &reportRepository{pg: pg}
}

func (r *reportRepository) GetStatistics(ctx context.Context) (models.DashboardStatistics, error) {
	query := `
		SELECT status, COUNT(*) 
		FROM achievement_references 
		WHERE deleted_at IS NULL 
		GROUP BY status
	`
	rows, err := r.pg.QueryContext(ctx, query)
	if err != nil {
		return models.DashboardStatistics{}, err
	}
	defer rows.Close()

	stats := models.DashboardStatistics{
		TotalPrestasi: 0,
		ByStatus:      map[string]int64{"draft": 0, "submitted": 0, "verified": 0, "rejected": 0},
	}

	for rows.Next() {
		var status string
		var count int64
		if err := rows.Scan(&status, &count); err == nil {
			stats.ByStatus[status] = count
			stats.TotalPrestasi += count
		}
	}
	return stats, nil
}

func (r *reportRepository) GetStudentProfile(ctx context.Context, studentID string) (models.StudentReportProfile, error) {
	query := `
		SELECT s.id, s.student_id as nim, u.full_name
		FROM students s
		JOIN users u ON s.user_id = u.id
		WHERE s.id = $1
	`
	var profile models.StudentReportProfile
	err := r.pg.QueryRowContext(ctx, query, studentID).Scan(&profile.StudentID, &profile.NIM, &profile.FullName)
	return profile, err
}

func (r *reportRepository) GetVerifiedAchievementsByStudentID(ctx context.Context, studentID string) ([]models.AchievementReference, error) {
	query := `
		SELECT id, student_id, mongo_achievement_id, status, created_at
		FROM achievement_references
		WHERE student_id = $1 AND deleted_at IS NULL AND status = 'verified'
		ORDER BY created_at DESC
	`
	rows, err := r.pg.QueryContext(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []models.AchievementReference
	for rows.Next() {
		var ref models.AchievementReference
		rows.Scan(&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status, &ref.CreatedAt)
		refs = append(refs, ref)
	}
	return refs, nil
}