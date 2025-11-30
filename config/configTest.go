package config

import (
	"log"

	"github.com/joho/godotenv"
)

func ConfigTest() {
	err := godotenv.Load("../.env")
	if err != nil{
		log.Fatal("Error loading .env ", err)
	}
}