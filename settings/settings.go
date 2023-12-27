package settings

import (
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Settings struct {
	DATABASE_URL        string `env:"DATABASE_URL"`
	CSV_FILE_FROM_SCRAP string `env:"CSV_FILE_FROM_SCRAP"`
}

var AppSettings *Settings

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	AppSettings = &Settings{}
	if err := env.Parse(AppSettings); err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}
}
