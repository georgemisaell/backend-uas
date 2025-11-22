package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	PasswordHash string `json:"-"`
	FullName string `json:"full_name"`
	RoleID uuid.UUID `json:"role_id"`
	IsActive bool `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUser struct {
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	RoleID uuid.UUID `json:"role_id"`
	IsActive bool `json:"is_active"`
}

type UpdateUser struct {
	Username string `json:"username"`
	Email string `json:"email"`
	FullName string `json:"full_name"`
	RoleID uuid.UUID `json:"role_id"`
	IsActive bool `json:"is_active"`
}