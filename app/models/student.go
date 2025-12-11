package models

import (
	"time"

	"github.com/google/uuid"
)

type Student struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	StudentID    string    `json:"student_id"`
	ProgramStudy string    `json:"program_study"`
	AcademicYear string    `json:"academy_year"`
	AdvisorID    uuid.UUID `json:"advisor_id"`	
	CreatedAt    time.Time `json:"created_at"`
}

type GetStudent struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	NIM          string    `json:"nim"`
	FullName     string    `json:"full_name"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	RoleName     string    `json:"role_name"`
	ProgramStudy string    `json:"program_study"`
	AcademyYear  string    `json:"academy_year"`
	IsActive     bool      `json:"is_active"`
}