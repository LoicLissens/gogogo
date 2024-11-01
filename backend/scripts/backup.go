package scripts

import (
	"context"
	"fmt"
	"io/fs"
	"jiva-guildes/adapters"
	"jiva-guildes/settings"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var bucketManager = adapters.BucketManager{}

func saveImgOnBucket() {
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
func saveDBDump() {
	folderPath, err := filepath.Abs(settings.AppSettings.DB_DUMP_FOLDER)
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
		dumpFile := splitedPath[len(splitedPath)-1]
		if strings.Contains(dumpFile, ".sql") {
			fmt.Println("\033[33m" + "Attempting to save" + dumpFile + "\033[0m")
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			override := false
			uploaderParam := uploader.UploadParams{
				PublicID:     dumpFile,
				AssetFolder:  settings.AppSettings.DB_DUMP_FOLDER,
				DisplayName:  dumpFile,
				ResourceType: "raw",
				Overwrite:    &override,
			}
			err = bucketManager.UploadFile(context.Background(), file, uploaderParam)
			if err != nil {
				return err
			}
			fmt.Println("\033[32m" + dumpFile + "saved" + "\033[0m")
			return nil
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

}
func main() {
	bucketManager.Setup(settings.AppSettings.BUCKET_API_KEY)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		saveImgOnBucket()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		saveDBDump()
	}()
	wg.Wait()
}
