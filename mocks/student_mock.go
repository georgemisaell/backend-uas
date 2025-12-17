package mocks

import (
	"context"
	"database/sql"
	"uas/app/models"

	"github.com/stretchr/testify/mock"
)

type MockStudentRepo struct {
	mock.Mock
}

func (m *MockStudentRepo) CreateStudent(ctx context.Context, tx *sql.Tx, student models.Student) error {
	args := m.Called(ctx, tx, student)
	return args.Error(0)
}

func (m *MockStudentRepo) GetAllStudentsByRole(ctx context.Context, roleName string) ([]models.GetStudent, error) { return nil, nil }
func (m *MockStudentRepo) GetStudentByID(ctx context.Context, id string) (models.GetStudent, error) { return models.GetStudent{}, nil }
func (m *MockStudentRepo) UpdateStudentAdvisor(ctx context.Context, studentID string, advisorID string) error { return nil }