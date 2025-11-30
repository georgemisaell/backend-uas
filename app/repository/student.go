package repository

import (
	"database/sql"
	"uas/app/models"
)

func CreateStudent (tx *sql.Tx, student models.Student) error {
	    query := `
        INSERT INTO students (
            id, user_id, student_id, program_study, academy_year, advisor_id, created_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7)`

		_, err := tx.Exec(query,
			student.ID,
			student.UserID,
			student.StudentID,
			student.ProgramStudy,
			student.AcademicYear,
			student.AdvisorID,
			student.CreatedAt,
    )

		return err
}