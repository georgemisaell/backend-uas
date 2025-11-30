package test

import (
	"testing"
	"uas/config"
	"uas/database"
)

func TestDatabasePostgreSQL(t *testing.T) {
	// Menghubungkan ENV
	config.ConfigTest();

	postgreSQLDB := database.ConnectDB()

	// Cek apakah koneksi error
	if postgreSQLDB == nil {
		t.Fatal("Gagal menghubungkan ke database PostgreSQL")
	}

	err := postgreSQLDB.Ping()
	if err != nil {
		t.Fatalf("Failed to ping PostgreSQL: %v", err)
	}
}

func TestDatabaseMongoDB(t *testing.T) {
	mongoDB := database.ConnectMongoDB()

	// Cek apakah koneksi error
	if mongoDB == nil {
		t.Fatal("Gagal menghubungkan ke database MongoDB")
	}

	// err := mongoDB.Ping()
	// if err != nil {
	// 	t.Fatalf("Failed to ping MongoDB: %v", err)
	// }
}

