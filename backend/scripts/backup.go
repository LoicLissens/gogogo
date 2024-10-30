package main

import (
	"context"
	"fmt"
	"io/fs"
	"jiva-guildes/adapters"
	"jiva-guildes/settings"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func main() {
	bucketManager := adapters.BucketManager{}
	bucketManager.Setup(settings.AppSettings.BUCKET_API_KEY)
	folderPath, err := filepath.Abs(settings.AppSettings.IMG_FOLDER)
	if err != nil {
		panic(err)
	}
	err = filepath.WalkDir(folderPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == folderPath {
			return nil
		}
		splitedPath := strings.Split(path, "/")
		imageFile := splitedPath[len(splitedPath)-1]
		// walk dir will fire the func with the folder, we want to skip this iteration and only get the files
		if strings.Contains(imageFile, "jpg") || strings.Contains(imageFile, "jpeg") || strings.Contains(imageFile, "png") {
			fmt.Println("\033[33m" + "Attempting to save" + imageFile + "\033[0m")
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			override := false
			uuid := splitedPath[len(splitedPath)-2]
			uploaderParam := uploader.UploadParams{
				PublicID:     uuid + "/" + imageFile,
				AssetFolder:  "guildes" + "/" + uuid,
				DisplayName:  imageFile,
				ResourceType: "image",
				Overwrite:    &override,
			}
			err = bucketManager.UploadFile(context.Background(), file, uploaderParam)
			if err != nil {
				return err
			}
			fmt.Println("\033[32m" + imageFile + "saved" + "\033[0m")
			return nil
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
}
