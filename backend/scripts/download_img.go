package main

import (
	"fmt"
	"io"
	"jiva-guildes/adapters/db"
	"jiva-guildes/adapters/db/repositories"
	"jiva-guildes/settings"
	"net/http"
	"os"
	"path/filepath"
)

func findImgUrl(slice []string, pathString string) bool {
	for _, item := range slice {
		if item == pathString {
			return true
		}
	}
	return false
}

func main() {
	pool := db.MountDB(settings.AppSettings.DATABASE_URI)
	defer db.Teardown(pool)
	repo := repositories.NewGuildeRepository(pool)
	folderPath, err := filepath.Abs(settings.AppSettings.IMG_FOLDER)
	if err != nil {
		panic(err)
	}
	guildes, err := repo.GetAll()
	if err != nil {
		panic(err)
	}
	for _, guilde := range guildes {
		url := guilde.Img_url
		if url == "" {
			continue
		}
		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			fmt.Println("\033[31m" + fmt.Sprintf("Unable to retrive image for %s", guilde.Name) + "\033[0m")
			continue
		}
		fmt.Println("\033[33m" + "Download img for " + guilde.Name + "\033[0m")
		dirPath := filepath.Join(folderPath, guilde.Uuid.String())
		err = os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			continue
		}

		filesInDir, err := filepath.Glob(filepath.Join(dirPath, "*"))
		if err != nil {
			panic(err)
		}
		path := filepath.Join(dirPath, fmt.Sprintf("%s.jpg", guilde.Name))
		if findImgUrl(filesInDir, path) {
			fmt.Println("\033[31m" + fmt.Sprintf("Image already exists for %s", guilde.Name) + "\033[0m")
			continue
		}
		file, err := os.Create(path)
		if err != nil {
			fmt.Println("Error creating file:", err)
			continue
		}
		defer file.Close()

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			fmt.Println("Error saving image:", err)
			continue
		}

		fmt.Println("\033[32m" + "Image saved for " + guilde.Name + "\033[0m")
	}
}
