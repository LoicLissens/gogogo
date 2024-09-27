package views

import (
	"errors"
	"fmt"
	"jiva-guildes/adapters"
	"jiva-guildes/adapters/db"
	"jiva-guildes/adapters/db/tables"
	customerrors "jiva-guildes/domain/custom_errors"
	"jiva-guildes/domain/models/guilde"
	"jiva-guildes/domain/ports"
	"jiva-guildes/domain/ports/views"
	"jiva-guildes/domain/ports/views/dtos"
	"jiva-guildes/settings"

	"reflect"
	"testing"

	"github.com/google/uuid"
)

// TODO Try to use test suite (see:https://medium.com/nerd-for-tech/setup-and-teardown-unit-test-in-go-bd6fa1b785cd)

func setupTest(tb testing.TB) (ports.UnitOfWork, views.ViewsManager, func(tb testing.TB)) {
	tables.InitAllTables()

	var connectionPool = db.MountDB(settings.AppSettings.DATABASE_URI)
	var UnitOfWorkManager = adapters.NewUnitOfWorkManager(connectionPool)
	uow := UnitOfWorkManager.Start()
	var ViewsManager = NewViewsManager(connectionPool)

	return uow, ViewsManager, func(tb testing.TB) {
		tables.DropAllTables()
		db.Teardown(connectionPool)
	}
}

func createGuilde(uow ports.UnitOfWork, name string, img_url string, page_url string) guilde.Guilde {
	g := guilde.New(name, img_url, page_url)
	savedGuilde, err := uow.GuildeRepository().Save(*g)
	if err != nil {
		panic(err)
	}
	return savedGuilde
}
func TestFetchGuilde(t *testing.T) {
	uow, ViewsManager, teardownTest := setupTest(t)
	defer teardownTest(t)

	g := createGuilde(uow, "GUnit", "https://www.googleimage.com", "https://www.google.com")
	fetchedGuilde, err := ViewsManager.Guilde().Fetch(g.Uuid)
	if err != nil {
		t.Fatal(err)
	}
	if reflect.TypeOf(fetchedGuilde) != reflect.TypeOf(dtos.GuildeViewDTO{}) {
		t.Fatalf("Expected %s, got %s", reflect.TypeOf(dtos.GuildeViewDTO{}), reflect.TypeOf(fetchedGuilde))
	}
	if fetchedGuilde.Name != g.Name {
		t.Fatalf("Expected %s, got %s", g.Name, fetchedGuilde.Name)
	}
	if fetchedGuilde.Img_url != g.Img_url {
		t.Fatalf("Expected %s, got %s", g.Img_url, fetchedGuilde.Img_url)
	}
	if fetchedGuilde.Page_url != g.Page_url {
		t.Fatalf("Expected %s, got %s", g.Page_url, fetchedGuilde.Page_url)
	}
	if fetchedGuilde.Uuid != g.Uuid {
		t.Fatalf("Expected %s, got %s", g.Uuid, fetchedGuilde.Uuid)
	}
	if fetchedGuilde.Created_at != g.Created_at {
		t.Fatalf("Expected %s, got %s", g.Created_at, fetchedGuilde.Created_at)
	}
	if fetchedGuilde.Updated_at != g.Updated_at {
		t.Fatalf("Expected %s, got %s", g.Updated_at, fetchedGuilde.Updated_at)
	}
}

func TestFetchGuildeNotFound(t *testing.T) {
	var expectedError customerrors.ErrorNotFound
	_, ViewsManager, teardownTest := setupTest(t)
	defer teardownTest(t)

	_, err := ViewsManager.Guilde().Fetch(uuid.New())
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if err != nil && !errors.As(err, &expectedError) {
		t.Fatal("Expected ErrorNotFound, got", reflect.TypeOf(err))
	}
}
func TestListGuildes(t *testing.T) {
	uow, ViewsManager, teardownTest := setupTest(t)
	defer teardownTest(t)
	samples := []map[string]string{
		{"name": "Gunit", "url": "url", "img": "img"},
		{"name": "D12", "url": "url2", "img": "img2"},
		{"name": "Neptuns", "url": "url3", "img": "img3"},
	}
	for _, sample := range samples {
		createGuilde(uow, sample["name"], sample[""], "https://www.google.com")
	}
	fetchedGuildes, err := ViewsManager.Guilde().List(1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if reflect.TypeOf(fetchedGuildes) != reflect.TypeOf(dtos.GuildeListViewDTO{}) {
		t.Fatalf("Expected %s, got %s", reflect.TypeOf(dtos.GuildeListViewDTO{}), reflect.TypeOf(fetchedGuildes))
	}
	guildes := fetchedGuildes.Items
	fmt.Printf("%v\n", guildes)
	if len(guildes) != len(samples) {
		t.Fatalf("Expected %d, got %d", 4, len(guildes))
	}
	for i, guilde := range guildes {
		if guilde.Name != samples[i]["name"] {
			t.Fatalf("Expected %s, got %s", samples[i]["name"], guilde.Name)
		}
		if guilde.Img_url != samples[i]["img"] {
			t.Fatalf("Expected %s, got %s", samples[i]["img"], guilde.Img_url)
		}
		if guilde.Page_url != samples[i]["url"] {
			t.Fatalf("Expected %s, got %s", samples[i]["url"], guilde.Page_url)
		}
	}
}
