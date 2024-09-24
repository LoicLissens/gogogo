package scripts

import (
	"fmt"
	"jiva-guildes/adapters/db"
	"jiva-guildes/adapters/db/repositories"
	"jiva-guildes/domain/models/guilde"
	"jiva-guildes/settings"

	"encoding/csv"
	"os"
	"path/filepath"
)

func PopulateDBFromCSV() {
	pool := db.MountDB(settings.AppSettings.DATABASE_URI)
	defer db.Teardown(pool)
	repo := repositories.NewGuildeRepository(pool)

	filePath, _ := filepath.Abs(settings.AppSettings.CSV_FILE_FROM_SCRAP)
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	i := 0
	fmt.Println("Inserting records...")
	for _, record := range records {
		g := guilde.New(record[0], record[1], record[2])
		repo.Save(*g) // may need bulk insert later
		i++
	}
	fmt.Printf("Inserted %d records\n", i)
}
