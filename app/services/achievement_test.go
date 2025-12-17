package services_test

import (
	"net/http/httptest"
	"testing"
	"uas/app/models"
	"uas/app/services"
	"uas/mocks"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- TEST SUBMIT (Mahasiswa) ---
func TestSubmitAchievement_Success(t *testing.T) {
	mockRepo := new(mocks.MockAchievementRepo)
	service := services.NewAchievementService(mockRepo)
	
	mockRepo.On("GetStudentIDByUserID", mock.Anything, "user-mhs").Return("std-1", nil)
	mockRepo.On("GetAchievementByID", mock.Anything, "ach-1").Return(models.AchievementReference{
		ID: "ach-1", StudentID: "std-1", Status: "draft",
	}, nil)
	mockRepo.On("SubmitAchievement", mock.Anything, "ach-1").Return(nil)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-mhs")
		return c.Next()
	})
	app.Post("/submit/:id", service.SubmitAchievement)

	req := httptest.NewRequest("POST", "/submit/ach-1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestSubmitAchievement_Fail_NotOwner(t *testing.T) {
	mockRepo := new(mocks.MockAchievementRepo)
	service := services.NewAchievementService(mockRepo)

	mockRepo.On("GetStudentIDByUserID", mock.Anything, "user-maling").Return("std-2", nil)
	mockRepo.On("GetAchievementByID", mock.Anything, "ach-1").Return(models.AchievementReference{
		ID: "ach-1", StudentID: "std-1", Status: "draft",
	}, nil)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-maling")
		return c.Next()
	})
	app.Post("/submit/:id", service.SubmitAchievement)

	req := httptest.NewRequest("POST", "/submit/ach-1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 403, resp.StatusCode)
}

func TestVerifyAchievement_Success(t *testing.T) {
	mockRepo := new(mocks.MockAchievementRepo)
	service := services.NewAchievementService(mockRepo)

	mockRepo.On("GetLecturerIDByUserID", mock.Anything, "user-dosen").Return("lec-1", nil)
	mockRepo.On("GetAchievementByID", mock.Anything, "ach-1").Return(models.AchievementReference{
		ID: "ach-1", StudentID: "std-1", Status: "submitted",
	}, nil)
	mockRepo.On("CheckStudentAdvisorRelationship", mock.Anything, "lec-1", "std-1").Return(true, nil)
	mockRepo.On("VerifyAchievement", mock.Anything, "ach-1", "user-dosen").Return(nil)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-dosen")
		return c.Next()
	})
	app.Post("/verify/:id", service.VerifyAchievement)

	req := httptest.NewRequest("POST", "/verify/ach-1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

func TestVerifyAchievement_Fail_NotAdvisor(t *testing.T) {
	mockRepo := new(mocks.MockAchievementRepo)
	service := services.NewAchievementService(mockRepo)

	mockRepo.On("GetLecturerIDByUserID", mock.Anything, "user-dosen-asing").Return("lec-99", nil)
	mockRepo.On("GetAchievementByID", mock.Anything, "ach-1").Return(models.AchievementReference{
		ID: "ach-1", StudentID: "std-1", Status: "submitted",
	}, nil)
	
	mockRepo.On("CheckStudentAdvisorRelationship", mock.Anything, "lec-99", "std-1").Return(false, nil)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-dosen-asing")
		return c.Next()
	})
	app.Post("/verify/:id", service.VerifyAchievement)

	req := httptest.NewRequest("POST", "/verify/ach-1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 403, resp.StatusCode)
}