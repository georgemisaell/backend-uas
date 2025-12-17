package services_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"uas/app/models"
	"uas/app/services"
	"uas/mocks"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(pwd string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(bytes)
}

func TestLogin_Success(t *testing.T) {
	// 1. SETUP
	mockRepo := new(mocks.MockUserRepo)
	authService := services.NewAuthService(mockRepo)
	app := fiber.New()
	app.Post("/login", authService.Login)

	// Data Dummy (Password asli: "123456")
	dummyUser := models.User{
		ID:           uuid.New(),
		Username:     "george_ganteng",
		Email:        "george@gmail.com",
		PasswordHash: hashPassword("123456"),
		FullName:     "George Misael",
		RoleName:     "Mahasiswa",
		IsActive:     true,
	}

	mockRepo.On("GetByUsernameOrEmail", mock.Anything, "george_ganteng").Return(dummyUser, nil)

	input := map[string]string{
		"username": "george_ganteng",
		"password": "123456", // Password asli
	}
	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)

	var responseBody map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&responseBody)
	
	data := responseBody["data"].(map[string]interface{})
	assert.NotEmpty(t, data["token"])
	assert.Equal(t, "george_ganteng", data["user"].(map[string]interface{})["username"])
	
	mockRepo.AssertExpectations(t)
}

func TestLogin_WrongPassword(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	authService := services.NewAuthService(mockRepo)
	app := fiber.New()
	app.Post("/login", authService.Login)

	// User ada di DB
	dummyUser := models.User{
		Username:     "george_ganteng",
		PasswordHash: hashPassword("1234567"),
		IsActive:     true,
	}

	mockRepo.On("GetByUsernameOrEmail", mock.Anything, "george_ganteng").Return(dummyUser, nil)

	// Tapi input password SALAH
	input := map[string]string{
		"username": "george_ganteng",
		"password": "SALAH_BOS", 
	}
	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 401, resp.StatusCode)
}

func TestLogin_UserNotFound(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	authService := services.NewAuthService(mockRepo)
	app := fiber.New()
	app.Post("/login", authService.Login)

	mockRepo.On("GetByUsernameOrEmail", mock.Anything, "hantu_laut").Return(models.User{}, sql.ErrNoRows)

	input := map[string]string{
		"username": "hantu_laut",
		"password": "bebas",
	}
	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 401, resp.StatusCode)
}

func TestLogin_AccountInactive(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	authService := services.NewAuthService(mockRepo)
	app := fiber.New()
	app.Post("/login", authService.Login)

	dummyUser := models.User{
		Username:     "george_cuti",
		PasswordHash: hashPassword("123456"),
		IsActive:     false,
	}

	mockRepo.On("GetByUsernameOrEmail", mock.Anything, "george_cuti").Return(dummyUser, nil)

	input := map[string]string{
		"username": "george_cuti",
		"password": "123456",
	}
	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 403, resp.StatusCode)
}