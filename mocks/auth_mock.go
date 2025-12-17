package mocks

import (
	"context"
	"database/sql"
	"uas/app/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) GetByUsernameOrEmail(ctx context.Context, loginInput string) (models.User, error) {
	args := m.Called(ctx, loginInput)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserRepo) CreateUser(ctx context.Context, tx *sql.Tx, user models.User) error {
	args := m.Called(ctx, tx, user)
	return args.Error(0)
}

func (m *MockUserRepo) GetAllUsers(ctx context.Context) ([]models.User, error) { return nil, nil }
func (m *MockUserRepo) GetUserByID(ctx context.Context, id uuid.UUID) (models.User, error) { return models.User{}, nil }
func (m *MockUserRepo) UpdateUser(ctx context.Context, id uuid.UUID, user models.UpdateUser) error { return nil }
func (m *MockUserRepo) DeleteUser(ctx context.Context, id uuid.UUID) error { return nil }
func (m *MockUserRepo) UpdateUserRole(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error { return nil }