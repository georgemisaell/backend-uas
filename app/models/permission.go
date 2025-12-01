package models

import "github.com/google/uuid"

type Permission struct {
	ID          uuid.UUID
	Name        string
	Resource    string
	Action      string
	Description string
}