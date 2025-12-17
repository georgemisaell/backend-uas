package services_test

import (
	"errors"
	"net/http/httptest"
	"testing"
	"uas/app/models"
	"uas/app/services"
	"uas/mocks"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetSystemStatistics_Success(t *testing.T) {
	// 1. SETUP STUNTMAN (Mock)
	mockRepo := new(mocks.MockReportRepo)
	
	// Kita gak butuh AchievementRepo buat fungsi ini, jadi bisa nil
	reportService := services.NewReportService(mockRepo, nil) 

	// 2. SETUP DATA DUMMY (Skenario: DB mengembalikan data ini)
	dummyStats := models.DashboardStatistics{
		TotalPrestasi: 10,
		ByStatus: map[string]int64{
			"verified": 10,
		},
	}

	// 3. LATIH STUNTMAN
	// "Eh Repo, kalau kamu dipanggil fungsi GetStatistics, 
	// tolong balikin data dummyStats dan errornya nil ya!"
	mockRepo.On("GetStatistics", mock.Anything).Return(dummyStats, nil)

	// 4. SETUP FIBER (Pura-pura jadi Server)
	app := fiber.New()
	
	// Kita pasang middleware buat nipu c.Locals ("role_name")
	// Seolah-olah yang login adalah Admin
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("role_name", "Admin")
		return c.Next()
	})

	// Pasang route yang mau dites
	app.Get("/stats", reportService.GetSystemStatistics)

	// 5. EKSEKUSI (Action!)
	// Bikin request pura-pura (HTTP GET ke /stats)
	req := httptest.NewRequest("GET", "/stats", nil)
	resp, _ := app.Test(req)

	// 6. ASSERTION (Cek Hasil)
	// Harapannya: Status Code 200 OK
	assert.Equal(t, 200, resp.StatusCode)

	// Cek apakah Mock Repo beneran dipanggil?
	mockRepo.AssertExpectations(t)
}

func TestGetSystemStatistics_Forbidden(t *testing.T) {
	// Skenario: Yang akses Mahasiswa (Bukan Admin/Doswal)
	
	mockRepo := new(mocks.MockReportRepo)
	reportService := services.NewReportService(mockRepo, nil)

	app := fiber.New()
	// Nipu Locals jadi "Mahasiswa"
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("role_name", "Mahasiswa") 
		return c.Next()
	})
	app.Get("/stats", reportService.GetSystemStatistics)

	req := httptest.NewRequest("GET", "/stats", nil)
	resp, _ := app.Test(req)

	// Harapannya: Ditolak (403 Forbidden)
	assert.Equal(t, 403, resp.StatusCode)
	
	// Pastikan Repo TIDAK dipanggil (karena ditolak di gerbang service)
	mockRepo.AssertNotCalled(t, "GetStatistics")
}

func TestGetSystemStatistics_DBError(t *testing.T) {
	// Skenario: Admin akses, tapi Database Error
	
	mockRepo := new(mocks.MockReportRepo)
	reportService := services.NewReportService(mockRepo, nil)

	// Latih Stuntman buat balikin Error
	mockRepo.On("GetStatistics", mock.Anything).Return(models.DashboardStatistics{}, errors.New("database mati"))

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("role_name", "Admin")
		return c.Next()
	})
	app.Get("/stats", reportService.GetSystemStatistics)

	req := httptest.NewRequest("GET", "/stats", nil)
	resp, _ := app.Test(req)

	// Harapannya: Error Server (500)
	assert.Equal(t, 500, resp.StatusCode)
}