package settings

import (
	"log"
	"os"
	"path/filepath"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

const projectDireName string = "jiva-guildes"

type Settings struct {
	DATABASE_URI        string `env:"DATABASE_URI"`
	CSV_FILE_FROM_SCRAP string `env:"CSV_FILE_FROM_SCRAP"`
	DATABASE_SCHEMA     string `env:"DATABASE_SCHEMA"`
}

var AppSettings *Settings

func init() {
	rootDirectory, err := os.Getwd()
	homeDirectory, errHome := os.UserHomeDir()

	if err != nil || errHome != nil {
		log.Fatalf("Error getting current or home directory: %v", err)
	}
	//  Determines the root directory of the project by traversing up the directory tree until it finds the project directory (where the .env is located).
	for rootDirectory != homeDirectory && filepath.Base(rootDirectory) != projectDireName {
		rootDirectory = filepath.Dir(rootDirectory)
		err = os.Chdir(rootDirectory)
	}

	envPath := filepath.Join(rootDirectory, ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Println("No .env file found")
	}

	AppSettings = &Settings{}
	if err := env.Parse(AppSettings); err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}
}
