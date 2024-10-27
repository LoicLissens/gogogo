package main

import (
	"jiva-guildes/adapters/db"
	"jiva-guildes/adapters/db/repositories"
	"jiva-guildes/settings"
	"path/filepath"
)

func main() {
	pool := db.MountDB(settings.AppSettings.DATABASE_URI)
	defer db.Teardown(pool)
	repo := repositories.NewGuildeRepository(pool)
	folderPath, err := filepath.Abs(settings.AppSettings.IMG_FOLDER)
	if err != nil {
		panic(err)
	}
}
