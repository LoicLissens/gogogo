package repositories

import (
	"errors"
	"fmt"
	"jiva-guildes/adapters/db"
	"jiva-guildes/adapters/db/tables"
	"jiva-guildes/adapters/db/test_utils"
	customerrors "jiva-guildes/domain/custom_errors"
	"jiva-guildes/domain/models"
	"jiva-guildes/domain/models/guilde"
	corerepo "jiva-guildes/domain/ports/repositories"
	"jiva-guildes/settings"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func setupTest(tb testing.TB) (corerepo.GuildeRepository, func(tb testing.TB)) {
	tables.InitAllTables()
	pool := db.MountDB(settings.AppSettings.DATABASE_URI)
	repo := NewGuildeRepository(pool)
	return &repo, func(tb testing.TB) {
		tables.DropAllTables()
		db.Teardown(pool)
	}
}
func TestGetSaveRepository(t *testing.T) {
	repo, teardownTest := setupTest(t)
	defer teardownTest(t)
	opts := guilde.GuildeOptions{Name: "Test", Page_url: "url", Exists: true, Validated: true, Active: nil, Creation_date: nil}
	entity, err := guilde.New(opts)
	if err != nil {
		t.Error(fmt.Errorf("Error creating entity: %w", err))
	}
	savedEntity, err := repo.Save(*entity)
	if err != nil {
		t.Fatal(fmt.Errorf("Error saving entity: %w", err))
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
	if savedEntity.Exists != entity.Exists {
		t.Fatalf("Expected %t, got %t", entity.Exists, savedEntity.Exists)

	}
	if savedEntity.Validated != entity.Validated {
		t.Fatalf("Expected %t, got %t", entity.Validated, savedEntity.Validated)
	}

	if savedEntity.Active != entity.Active {
		t.Fatalf("Expected %t, got %t", *entity.Active, *savedEntity.Active)
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
	repo, teardownTest := setupTest(t)
	defer teardownTest(t)
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
	creationDate := time.Now().UTC()
	entity, err := guilde.New(guilde.GuildeOptions{Name: "Test", Img_url: "img", Page_url: "page", Exists: true, Validated: true, Active: nil, Creation_date: &creationDate})
	if err != nil {
		t.Fatal(err)
	}
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

	creationDate := time.Now().UTC()
	entity, err := guilde.New(guilde.GuildeOptions{Name: "Test", Img_url: "img", Page_url: "page", Exists: true, Validated: true, Active: nil, Creation_date: &creationDate})
	if err != nil {
		t.Fatal(err)
	}
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

func TestUpdate(t *testing.T) {
	repo, teardownTest := setupTest(t)
	defer teardownTest(t)
	creationDate := time.Now().UTC()
	entity, err := guilde.New(guilde.GuildeOptions{Name: "Test", Img_url: "img", Page_url: "page", Exists: true, Validated: false, Active: nil, Creation_date: &creationDate})
	if err != nil {
		t.Fatal(err)
	}
	savedEntity, err := repo.Save(*entity)
	if err != nil {
		t.Fatal(err)
	}
	active := true
	savedEntity.Validated = true
	savedEntity.Active = &active
	savedEntity.Creation_date = nil
	firstUpdatedDate := savedEntity.Updated_at
	updatedEntity, err := repo.Update(savedEntity)
	if err != nil {
		t.Fatal(err)
	}
	if !updatedEntity.Validated {
		t.Fatal("Expected validated to be true, got false")
	}
	if *updatedEntity.Active != true {
		t.Fatal("Expected active to be true, got ", *updatedEntity.Active)
	}
	if updatedEntity.Creation_date != nil {
		t.Fatal("Expected creation date to be nil, got ", updatedEntity.Creation_date)
	}
	if updatedEntity.Updated_at == firstUpdatedDate {
		t.Fatal("Expected updated date to be different, got same")
	}
}

func TestUpdateNotFound(t *testing.T) {
	var expectedError customerrors.ErrorNotFound
	repo, teardownTest := setupTest(t)
	defer teardownTest(t)
	creationDate := time.Now().UTC()
	entity, err := guilde.New(guilde.GuildeOptions{Name: "Test", Img_url: "img", Page_url: "page", Exists: true, Validated: false, Active: nil, Creation_date: &creationDate})
	if err != nil {
		t.Fatal(err)
	}
	_, err = repo.Update(*entity)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if err != nil && !errors.As(err, &expectedError) {
		t.Fatal("Expected ErrorNotFound, got", reflect.TypeOf(err))
	}
}
func TestGetAll(t *testing.T) {
	repo, teardownTest := setupTest(t)
	defer teardownTest(t)
	samples := test_utils.SaveBasicSamples(repo)
	guildes, err := repo.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(guildes) != len(samples) {
		t.Fatalf("Expected %d, got %d", len(samples), len(guildes))
	}
}
func TestGetAllEmptyDB(t *testing.T) {
	repo, teardownTest := setupTest(t)
	defer teardownTest(t)
	guildes, err := repo.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(guildes) != 0 {
		t.Fatalf("Expected %d, got %d", 0, len(guildes))
	}
}
