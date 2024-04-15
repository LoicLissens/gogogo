package repositories

import (
	"jiva-guildes/settings"
)

var schema string = settings.AppSettings.DATABASE_SCHEMA

// func TestSaveEntity(t *testing.T) {
// 	pool := db.MountDB()
// 	defer pool.Close()
// 	_, err := pool.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS  test_table (id INT PRIMARY KEY, name TEXT);")
// 	pool.Exec(context.Background(), "TRUNCATE TABLE test_table")

// 	if err != nil {
// 		log.Fatalf("Failed to execute query: %v", err)
// 	}

// 	type TestEntity struct {
// 		ID   int    `db:"id"`
// 		Name string `db:"name"`
// 	}

// 	entity := TestEntity{
// 		ID:   1,
// 		Name: "Test Entity",
// 	}
// 	// Call the SaveEntity function
// 	savedEntity, err := repositories.SaveEntity(pool, "test_table", schema, entity)
// 	if err != nil {
// 		t.Fatalf("SaveEntity returned an error: %v", err)
// 	}

// 	// Assert the saved entity
// 	expectedEntity := TestEntity{
// 		ID:   1,
// 		Name: "Test Entity",
// 	}
// 	if !reflect.DeepEqual(savedEntity, expectedEntity) {
// 		t.Errorf("Saved entity does not match the expected entity. Got: %+v, Expected: %+v", savedEntity, expectedEntity)
// 	}
// }
