package repository

import (
	"database/sql"
	"uas/app/models"
)

func CreateLecture(tx *sql.Tx, lecture models.Lecture) error {
	query := `
		INSERT INTO lecturers (
			id, user_id, lecturer_id, department, created_at
		) VALUES ($1, $2, $3, $4, $5)
	`
	_, err := tx.Exec(query,
		lecture.ID,
		lecture.UserID,
		lecture.LectureID,
		lecture.Department,
		lecture.CreatedAt,
	)

	return err
}