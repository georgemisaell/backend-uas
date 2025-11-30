package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	PasswordHash string `json:"-"`
	FullName string `json:"full_name"`
	RoleID uuid.UUID `json:"role_id"`
	RoleName string `json:"role_name"`
	Student *Student `json:"student"` 
  Lecture *Lecture `json:"lecture"`
	IsActive bool `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Ini struct untuk menangkap Request JSON
type CreateUserRequest struct {
    // Data User
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
    FullName string `json:"full_name"`
    RoleID   string `json:"role_id"`
    RoleName string `json:"role_name"`

    // Data Terhubung (Nested Objects)
    // Gunakan Pointer (*)
    Student *Student `json:"student"` 
    Lecture *Lecture `json:"lecture"`
}

type CreateUser struct {
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	RoleID uuid.UUID `json:"role_id"`
	RoleName string `json:"role_name"`
	IsActive bool `json:"is_active"`
}

type UpdateUser struct {
	Username string `json:"username"`
	Email string `json:"email"`
	FullName string `json:"full_name"`
	RoleID uuid.UUID `json:"role_id"`
	IsActive bool `json:"is_active"`
}

type LoginRequest struct { 
	Username string `json:"username"` 
	Password string `json:"password"` 
}

type LoginResponse struct { 
	User  User   `json:"user"` 
	Token string `json:"token"`
}

type JWTClaims struct { 
	UserID   uuid.UUID  `json:"user_id"` 
	Username string `json:"username"` 
	RoleName string `json:"role_name"` 
	jwt.RegisteredClaims
}