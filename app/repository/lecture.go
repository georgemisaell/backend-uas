package repository

import (
	"context"
	"database/sql"
	"fmt"
	"uas/app/models"
)

type LecturerRepository interface {
	CreateLecture(ctx context.Context, tx *sql.Tx, lecture models.Lecture) error
	GetAllLecturersByRole(ctx context.Context, roleName string) ([]models.GetLecture, error)
	GetLecturerByID(ctx context.Context, id string) (models.GetLecture, error)
	GetAdviseesByLecturerID(ctx context.Context, lecturerID string) ([]models.GetStudent, error)
}

type lecturerRepository struct {
	db *sql.DB
}

func NewLecturerRepository(db *sql.DB) LecturerRepository {
	return &lecturerRepository{db: db}
}

func (r *lecturerRepository) CreateLecture(ctx context.Context, tx *sql.Tx, lecture models.Lecture) error {
	query := `
		INSERT INTO lecturers (
			id, user_id, lecturer_id, department, created_at
		) VALUES ($1, $2, $3, $4, $5)
	`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query,
			lecture.ID,
			lecture.UserID,
			lecture.LectureID,
			lecture.Department,
			lecture.CreatedAt,
		)
	} else {
		_, err = r.db.ExecContext(ctx, query,
			lecture.ID,
			lecture.UserID,
			lecture.LectureID,
			lecture.Department,
			lecture.CreatedAt,
		)
	}

	return err
}

func (r *lecturerRepository) GetAllLecturersByRole(ctx context.Context, roleName string) ([]models.GetLecture, error) {
	query := `
		SELECT 
			l.id, 
			l.user_id, 
			l.lecturer_id, 
			l.department, 
			u.full_name, 
			u.username, 
			u.email, 
			u.is_active,
			l.created_at,
			r.name as role_name
		FROM lecturers l
		JOIN users u ON l.user_id = u.id
		JOIN roles r ON u.role_id = r.id
		WHERE r.name = $1
		ORDER BY u.full_name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, roleName)
	if err != nil {
		return nil, fmt.Errorf("gagal query lecturers: %w", err)
	}
	defer rows.Close()

	var lecturers []models.GetLecture

	for rows.Next() {
		var l models.GetLecture
		err := rows.Scan(
			&l.ID,
			&l.UserID,
			&l.LecturerID,
			&l.Department,
			&l.FullName,
			&l.Username,
			&l.Email,
			&l.IsActive,
			&l.CreatedAt,
			&l.RoleName,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal scanning row dosen: %w", err)
		}
		lecturers = append(lecturers, l)
	}

	return lecturers, nil
}

func (r *lecturerRepository) GetLecturerByID(ctx context.Context, id string) (models.GetLecture, error) {
	query := `
		SELECT 
			l.id, l.user_id, l.lecturer_id, l.department, 
			u.full_name, u.username, u.email, u.is_active, l.created_at,
			r.name
		FROM lecturers l
		JOIN users u ON l.user_id = u.id
		JOIN roles r ON u.role_id = r.id
		WHERE l.id = $1
	`

	var l models.GetLecture
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&l.ID, &l.UserID, &l.LecturerID, &l.Department,
		&l.FullName, &l.Username, &l.Email, &l.IsActive, &l.CreatedAt,
		&l.RoleName,
	)

	if err != nil {
		return models.GetLecture{}, err
	}

	return l, nil
}

func (r *lecturerRepository) GetAdviseesByLecturerID(ctx context.Context, lecturerID string) ([]models.GetStudent, error) {
	query := `
		SELECT 
			s.id, 
			s.user_id, 
			s.student_id, 
			s.program_study, 
			s.academy_year,
			u.full_name, 
			u.username, 
			u.email,
			u.is_active,
			r.name as role_name
		FROM students s
		JOIN users u ON s.user_id = u.id
		JOIN roles r ON u.role_id = r.id
		WHERE s.advisor_id = $1
		ORDER BY u.full_name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, lecturerID)
	if err != nil {
		return nil, fmt.Errorf("gagal query advisees: %w", err)
	}
	defer rows.Close()

	var students []models.GetStudent

	for rows.Next() {
		var s models.GetStudent
		err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.NIM,
			&s.ProgramStudy,
			&s.AcademyYear, 
			&s.FullName,
			&s.Username,
			&s.Email,
			&s.IsActive,
			&s.RoleName,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal scanning row mahasiswa bimbingan: %w", err)
		}
		students = append(students, s)
	}

	return students, nil
}