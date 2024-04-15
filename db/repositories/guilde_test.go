package repositories

import (
	"jiva-guildes/db"
	"jiva-guildes/domain/models/guilde"
	"testing"
)

// TODO Try to close the connection at the end of all tests
func TestGetSaveRepository(t *testing.T) {
	pool := db.MountDB()
	defer pool.Close()
	var repo GuildeRepository = GuildeRepository{conn: pool}

	entity := guilde.New("Test", "test", "test")
	savedEntity, err := repo.Save(*entity)
	if err != nil {
		t.Fatal(err)
	}
	if savedEntity.Name != entity.Name {
		t.Fatalf("Expected %s, got %s", entity.Name, savedEntity.Name)
	}
	if savedEntity.Img_url != entity.Img_url {
		t.Fatalf("Expected %s, got %s", entity.Img_url, savedEntity.Img_url)
	}
	if savedEntity.Page_url != entity.Page_url {
		t.Fatalf("Expected %s, got %s", entity.Page_url, savedEntity.Page_url)
	}
	if savedEntity.Uuid != entity.Uuid {
		t.Fatalf("Expected %s, got %s", entity.Uuid, savedEntity.Uuid)
	}
	if savedEntity.Created_at != entity.Created_at {
		t.Fatalf("Expected %s, got %s", entity.Created_at, savedEntity.Created_at)
	}
	if savedEntity.Updated_at != entity.Updated_at {
		t.Fatalf("Expected %s, got %s", entity.Updated_at, savedEntity.Updated_at)
	}

	fetchedEntity, error := repo.GetByUUID(entity.Uuid)

	if error != nil {
		t.Fatal(error)
	}
	if fetchedEntity.Name != entity.Name {
		t.Fatalf("Expected %s, got %s", entity.Name, savedEntity.Name)
	}
	if fetchedEntity.Img_url != entity.Img_url {
		t.Fatalf("Expected %s, got %s", entity.Img_url, savedEntity.Img_url)
	}
	if fetchedEntity.Page_url != entity.Page_url {
		t.Fatalf("Expected %s, got %s", entity.Page_url, savedEntity.Page_url)
	}
	if fetchedEntity.Uuid != entity.Uuid {
		t.Fatalf("Expected %s, got %s", entity.Uuid, savedEntity.Uuid)
	}
	if fetchedEntity.Created_at != entity.Created_at {
		t.Fatalf("Expected %s, got %s", entity.Created_at, savedEntity.Created_at)
	}
	if fetchedEntity.Updated_at != entity.Updated_at {
		t.Fatalf("Expected %s, got %s", entity.Updated_at, savedEntity.Updated_at)
	}
}
