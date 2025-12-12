package models

import (
	"time"

	"github.com/google/uuid"
)

type Lecture struct {
	ID uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	LectureID string `json:"lecturer_id"`
	Department string `json:"department"`
	CreatedAt time.Time `json:"created_at"`
}

type GetLecture struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	NIM string `json:"student_id"`
	ProgramStudy string `json:"program_study"`
	AcademyYear string `json:"academy_year"`
	LecturerID string    `json:"lecturer_id"`
	FullName   string    `json:"full_name"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	RoleName   string    `json:"role_name"`
	Department string    `json:"department"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
}