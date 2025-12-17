package mocks

import (
	"context"
	"uas/app/models"

	"github.com/stretchr/testify/mock"
)

type MockAchievementRepo struct {
	mock.Mock
}

func (m *MockAchievementRepo) GetStudentIDByUserID(ctx context.Context, userID string) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *MockAchievementRepo) GetLecturerIDByUserID(ctx context.Context, userID string) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *MockAchievementRepo) CheckStudentAdvisorRelationship(ctx context.Context, lecturerID string, studentID string) (bool, error) {
	args := m.Called(ctx, lecturerID, studentID)
	return args.Bool(0), args.Error(1)
}

func (m *MockAchievementRepo) GetAchievementByID(ctx context.Context, id string) (models.AchievementReference, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.AchievementReference), args.Error(1)
}

func (m *MockAchievementRepo) SubmitAchievement(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAchievementRepo) VerifyAchievement(ctx context.Context, id string, verifierUserID string) error {
	args := m.Called(ctx, id, verifierUserID)
	return args.Error(0)
}

func (m *MockAchievementRepo) RejectAchievement(ctx context.Context, id string, verifierUserID string, note string) error {
	args := m.Called(ctx, id, verifierUserID, note)
	return args.Error(0)
}

func (m *MockAchievementRepo) GetMongoDetailsByIDs(ctx context.Context, mongoIDs []string) (map[string]models.AchievementMongo, error) {
	args := m.Called(ctx, mongoIDs)
	return args.Get(0).(map[string]models.AchievementMongo), args.Error(1)
}

func (m *MockAchievementRepo) CreateAchievementReference(ctx context.Context, ref models.AchievementReference) error { return nil }
func (m *MockAchievementRepo) CreateAchievementMongo(ctx context.Context, data models.AchievementMongo) (string, error) { return "mongo-id", nil }
func (m *MockAchievementRepo) UpdateAchievement(ctx context.Context, pgID string, mongoID string, data models.AchievementMongo) error { return nil }
func (m *MockAchievementRepo) SoftDeleteAchievement(ctx context.Context, pgID string, mongoID string) error { return nil }
func (m *MockAchievementRepo) AddAttachmentToMongo(ctx context.Context, mongoID string, attachment models.Attachment) error { return nil }
func (m *MockAchievementRepo) GetAllReferences(ctx context.Context, filterUserID string) ([]models.AchievementReference, map[string]string, map[string]string, error) { return nil, nil, nil, nil }
func (m *MockAchievementRepo) GetAchievementReferenceWithDetail(ctx context.Context, id string) (models.AchievementResponse, error) { return models.AchievementResponse{}, nil }
func (m *MockAchievementRepo) GetMongoDetailByID(ctx context.Context, mongoID string) (models.AchievementMongo, error) { return models.AchievementMongo{}, nil }