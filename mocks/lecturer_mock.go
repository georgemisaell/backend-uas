package mocks

import (
	"context"
	"database/sql"
	"uas/app/models"

	"github.com/stretchr/testify/mock"
)

type MockLecturerRepo struct {
	mock.Mock
}

func (m *MockLecturerRepo) CreateLecture(ctx context.Context, tx *sql.Tx, lecture models.Lecture) error {
	args := m.Called(ctx, tx, lecture)
	return args.Error(0)
}

// Method lain (Dummy)
func (m *MockLecturerRepo) GetAllLecturersByRole(ctx context.Context, roleName string) ([]models.GetLecture, error) { return nil, nil }
func (m *MockLecturerRepo) GetLecturerByID(ctx context.Context, id string) (models.GetLecture, error) { return models.GetLecture{}, nil }
func (m *MockLecturerRepo) GetAdviseesByLecturerID(ctx context.Context, lecturerID string) ([]models.GetStudent, error) { return nil, nil }