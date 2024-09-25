package repositories

import (
	"errors"
	"jiva-guildes/adapters/db"
	customerrors "jiva-guildes/domain/custom_errors"
	"jiva-guildes/domain/models"
	"jiva-guildes/domain/models/guilde"
	"jiva-guildes/settings"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func setupTest(tb testing.TB) (GuildeRepository, func(tb testing.TB)) {
	pool := db.MountDB(settings.AppSettings.DATABASE_URI)
	var repo GuildeRepository = GuildeRepository{conn: pool}
	return repo, func(tb testing.TB) {
		db.Teardown(pool)
	}
}
func TestGetSaveRepository(t *testing.T) {
	repo, teardownTest := setupTest(t)
	defer teardownTest(t)
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

func TestGetNotfound(t *testing.T) {
	var expectedError customerrors.ErrorNotFound
	pool := db.MountDB(settings.AppSettings.DATABASE_URI)
	defer db.Teardown(pool)
	var repo GuildeRepository = GuildeRepository{conn: pool}
	_, err := repo.GetByUUID(uuid.New())
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if err != nil && !errors.As(err, &expectedError) {
		t.Fatal("Expected ErrorNotFound, got", reflect.TypeOf(err))
	}
}

func TestSaveDuplicated(t *testing.T) {
	var expectedError customerrors.ErrorAlreadyExists
	repo, teardownTest := setupTest(t)
	defer teardownTest(t)

	entity := guilde.New("Test", "test", "test")
	savedEntity, err := repo.Save(*entity)
	if err != nil {
		t.Fatal(err)
	}
	duplicatedUuid := savedEntity.Uuid
	_, duplictedErr := repo.Save(guilde.Guilde{
		Name:     "TestDup",
		Img_url:  "test",
		Page_url: "test",
		BaseModel: models.BaseModel{
			Uuid:       duplicatedUuid,
			Created_at: time.Now().UTC(),
			Updated_at: time.Now().UTC(),
		}})
	if duplictedErr == nil {
		t.Fatal("Expected error, got nil")
	}
	if err != nil && !errors.As(err, &expectedError) {
		t.Fatal("Expected ErrorNotFound, got", reflect.TypeOf(err))
	}
}

func TestDelete(t *testing.T) {
	repo, teardownTest := setupTest(t)
	defer teardownTest(t)

	entity := guilde.New("Test", "test", "test")
	savedEntity, err := repo.Save(*entity)
	if err != nil {
		t.Fatal(err)
	}
	delError := repo.Delete(savedEntity.Uuid)
	if delError != nil {
		t.Fatal(err)
	}
	_, err = repo.GetByUUID(savedEntity.Uuid)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	var expectedError customerrors.ErrorNotFound
	if err != nil && !errors.As(err, &expectedError) {
		t.Fatal("Expected ErrorNotFound, got", reflect.TypeOf(err))
	}
}
