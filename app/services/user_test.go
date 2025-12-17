package services_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"uas/app/models"
	"uas/app/services"
	"uas/mocks"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser_Student_Success(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockUserRepo := new(mocks.MockUserRepo)
	mockStudentRepo := new(mocks.MockStudentRepo)
	mockLecturerRepo := new(mocks.MockLecturerRepo)

	userService := services.NewUserService(db, mockUserRepo, mockStudentRepo, mockLecturerRepo)

	app := fiber.New()
	app.Post("/users", userService.CreateUser)

	mockDB.ExpectBegin()

	mockUserRepo.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockStudentRepo.On("CreateStudent", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockDB.ExpectCommit()

	// UUID Dummy
  dummyAdvisorID := uuid.MustParse("00000000-0000-0000-0000-000000000002")

	input := models.CreateUserRequest{
		Username: "maba_2025",
		Email:    "maba@kampus.ac.id",
		Password: "password123",
		FullName: "Maba Ganteng",
		RoleName: "Mahasiswa",
		RoleID:   "00000000-0000-0000-0000-000000000001",
		Student: &models.Student{
			StudentID:    "A11.2025.00001",
			ProgramStudy: "Teknik Informatika",
			AcademicYear: "2025",
			AdvisorID: dummyAdvisorID,
		},
	}
	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	// 6. ASSERT
	assert.Equal(t, 201, resp.StatusCode)
	
	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	mockUserRepo.AssertExpectations(t)
	mockStudentRepo.AssertExpectations(t)
}

func TestCreateUser_Rollback_OnError(t *testing.T) {
	db, mockDB, _ := sqlmock.New()
	defer db.Close()

	mockUserRepo := new(mocks.MockUserRepo)
	userService := services.NewUserService(db, mockUserRepo, nil, nil)

	app := fiber.New()
	app.Post("/users", userService.CreateUser)

	mockDB.ExpectBegin()

	mockUserRepo.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError)

	mockDB.ExpectRollback()

	input := models.CreateUserRequest{
		Username: "error_user",
		Password: "123",
		RoleID:   "00000000-0000-0000-0000-000000000001",
	}
	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 500, resp.StatusCode)
	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}