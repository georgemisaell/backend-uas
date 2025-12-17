package mocks

import (
	"context"
	"uas/app/models"

	"github.com/stretchr/testify/mock"
)

type MockReportRepo struct {
	mock.Mock
}

func (m *MockReportRepo) GetStatistics(ctx context.Context) (models.DashboardStatistics, error) {
	args := m.Called(ctx)

	return args.Get(0).(models.DashboardStatistics), args.Error(1)
}

func (m *MockReportRepo) GetStudentProfile(ctx context.Context, studentID string) (models.StudentReportProfile, error) {
	args := m.Called(ctx, studentID)
	return args.Get(0).(models.StudentReportProfile), args.Error(1)
}

func (m *MockReportRepo) GetVerifiedAchievementsByStudentID(ctx context.Context, studentID string) ([]models.AchievementReference, error) {
	args := m.Called(ctx, studentID)
	return args.Get(0).([]models.AchievementReference), args.Error(1)
}