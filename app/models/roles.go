package models

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID          uuid.UUID
	Name        string
	Description string
	CreatedAt   time.Time
}

const (
	RoleAdmin     = "Admin"
	RoleMahasiswa = "Mahasiswa"
	RoleDosen     = "Dosen Wali"
)