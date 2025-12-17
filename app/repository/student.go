package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"uas/app/models"
)

type StudentRepository interface {
	CreateStudent(ctx context.Context, tx *sql.Tx, student models.Student) error
	GetAllStudentsByRole(ctx context.Context, roleName string) ([]models.GetStudent, error)
	GetStudentByID(ctx context.Context, id string) (models.GetStudent, error)
	UpdateStudentAdvisor(ctx context.Context, studentID string, advisorID string) error
}

type studentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) StudentRepository {
	return &studentRepository{db: db}
}

func (r *studentRepository) CreateStudent(ctx context.Context, tx *sql.Tx, student models.Student) error {
	query := `
		INSERT INTO students (
			id, user_id, student_id, program_study, academy_year, advisor_id, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query,
			student.ID,
			student.UserID,
			student.StudentID,
			student.ProgramStudy,
			student.AcademicYear,
			student.AdvisorID,
			student.CreatedAt,
		)
	} else {
		_, err = r.db.ExecContext(ctx, query,
			student.ID,
			student.UserID,
			student.StudentID,
			student.ProgramStudy,
			student.AcademicYear,
			student.AdvisorID,
			student.CreatedAt,
		)
	}

	return err
}

func (r *studentRepository) GetAllStudentsByRole(ctx context.Context, roleName string) ([]models.GetStudent, error) {
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
		WHERE r.name = $1
		ORDER BY u.full_name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, roleName)
	if err != nil {
		return nil, fmt.Errorf("gagal query students join roles: %w", err)
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
			return nil, fmt.Errorf("gagal scanning row: %w", err)
		}
		students = append(students, s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterasi rows: %w", err)
	}

	return students, nil
}

func (r *studentRepository) GetStudentByID(ctx context.Context, id string) (models.GetStudent, error) {
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
		WHERE s.id = $1
	`

	var s models.GetStudent

	err := r.db.QueryRowContext(ctx, query, id).Scan(
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
		return models.GetStudent{}, err
	}

	return s, nil
}

func (r *studentRepository) UpdateStudentAdvisor(ctx context.Context, studentID string, advisorID string) error {
	query := `
		UPDATE students 
		SET advisor_id = $1, updated_at = NOW()
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, advisorID, studentID)
	if err != nil {
		return fmt.Errorf("gagal update dosen wali: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("student_not_found")
	}

	return nil
}