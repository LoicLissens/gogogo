package views

import (
	"errors"
	"jiva-guildes/adapters"
	"jiva-guildes/adapters/db"
	"jiva-guildes/adapters/db/tables"
	customerrors "jiva-guildes/domain/custom_errors"
	"jiva-guildes/domain/models/guilde"
	"jiva-guildes/domain/ports"
	"jiva-guildes/domain/ports/views"
	"jiva-guildes/domain/ports/views/dtos"
	"jiva-guildes/settings"
	"time"

	"reflect"
	"testing"

	"github.com/google/uuid"
)

// UTILS

func setupTest(tb testing.TB) (ports.UnitOfWorkManager, views.ViewsManager, func(tb testing.TB)) {
	tables.InitAllTables()

	var connectionPool = db.MountDB(settings.AppSettings.DATABASE_URI)
	var UnitOfWorkManager = adapters.NewUnitOfWorkManager(connectionPool)
	var ViewsManager = NewViewsManager(connectionPool)

	return UnitOfWorkManager, ViewsManager, func(tb testing.TB) {
		tables.DropAllTables()
		db.Teardown(connectionPool)
	}
}

func createGuilde(uowm ports.UnitOfWorkManager, opts guilde.GuildeOptions) guilde.Guilde {
	uow, close := uowm.Start()
	defer close()
	g, err := guilde.New(opts)
	if err != nil {
		panic(err)
	}
	savedGuilde, err := uow.GuildeRepository().Save(*g)
	if err != nil {
		panic(err)
	}
	return savedGuilde
}
func saveBatchSamples(uowm ports.UnitOfWorkManager, samples []guilde.GuildeOptions) []guilde.GuildeOptions {

	for _, sample := range samples {
		createGuilde(uowm, sample)
	}
	return samples
}
func saveBasicSamples(uowm ports.UnitOfWorkManager) []guilde.GuildeOptions {
	active := true
	notActive := false
	today := time.Now().UTC()
	date := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	samples := []guilde.GuildeOptions{
		{Name: "GUnit", Page_url: "https://www.googleimage.com", Exists: true, Validated: true, Active: nil, Creation_date: nil},
		{Name: "D12", Page_url: "img1", Exists: true, Validated: true, Active: nil, Creation_date: nil},
		{Name: "eminem", Exists: true, Validated: true, Active: &active, Creation_date: &today},
		{Name: "AS$AP", Page_url: "img3", Exists: false, Validated: false, Active: &notActive, Creation_date: &date},
	}
	return saveBatchSamples(uowm, samples)
}

// TESTS
func TestFetchGuilde(t *testing.T) {
	uowm, ViewsManager, teardownTest := setupTest(t)
	defer teardownTest(t)

	g := createGuilde(uowm, guilde.GuildeOptions{Name: "GUnit", Page_url: "https://www.googleimage.com", Exists: true, Validated: true, Active: nil, Creation_date: nil})
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
	uowm, ViewsManager, teardownTest := setupTest(t)
	defer teardownTest(t)

	samples := saveBasicSamples(uowm)
	fetchedGuildes, err := ViewsManager.Guilde().List(views.ListGuildesViewOpts{BaseListViewOpts: views.BaseListViewOpts{
		Page:  1,
		Limit: 10,
	}})
	if err != nil {
		t.Fatal(err)
	}
	if reflect.TypeOf(fetchedGuildes) != reflect.TypeOf(dtos.GuildeListViewDTO{}) {
		t.Fatalf("Expected %s, got %s", reflect.TypeOf(dtos.GuildeListViewDTO{}), reflect.TypeOf(fetchedGuildes))
	}
	guildes := fetchedGuildes.Items
	if len(guildes) != len(samples) {
		t.Fatalf("Expected %d, got %d", len(samples), len(guildes))
	}
	for i, guilde := range guildes {
		if guilde.Name != samples[i].Name {
			t.Fatalf("Expected %s, got %s", samples[i].Name, guilde.Name)
		}
		if guilde.Img_url != samples[i].Img_url {
			t.Fatalf("Expected %s, got %s", samples[i].Img_url, guilde.Img_url)
		}
		if guilde.Page_url != samples[i].Page_url {
			t.Fatalf("Expected %s, got %s", samples[i].Page_url, guilde.Page_url)
		}
	}
	// Retrieve nothing
	fetchedGuildes, err = ViewsManager.Guilde().List(views.ListGuildesViewOpts{BaseListViewOpts: views.BaseListViewOpts{
		Page:  2,
		Limit: 10,
	}})
	if err != nil {
		t.Fatal(err)
	}
	guildes = fetchedGuildes.Items
	if len(guildes) != 0 {
		t.Fatalf("Expected %d, got %d", 0, len(guildes))
	}
	// Retrieve 2
	limit := 2
	fetchedGuildes, err = ViewsManager.Guilde().List(views.ListGuildesViewOpts{BaseListViewOpts: views.BaseListViewOpts{
		Page:  1,
		Limit: 2,
	}})
	if err != nil {
		t.Fatal(err)
	}
	guildes = fetchedGuildes.Items
	if len(guildes) != limit {
		t.Fatalf("Expected %d, got %d", limit, len(guildes))
	}
	// Tetst absurd value for pagination
	fetchedGuildes, err = ViewsManager.Guilde().List(views.ListGuildesViewOpts{BaseListViewOpts: views.BaseListViewOpts{
		Page:  0,
		Limit: -2,
	}})
	if err != nil {
		t.Fatal(err)
	}
	guildes = fetchedGuildes.Items
	if len(guildes) != len(samples) {
		t.Fatalf("Expected %d, got %d", len(samples), len(guildes))
	}
}
func TestListGuildesOrdering(t *testing.T) {
	uowm, ViewsManager, teardownTest := setupTest(t)
	defer teardownTest(t)

	saveBasicSamples(uowm)
	fetchedGuildes, err := ViewsManager.Guilde().List(views.ListGuildesViewOpts{BaseListViewOpts: views.BaseListViewOpts{
		OrderingMethod: views.DESC,
	},
		OrderBy: views.OrderByName,
	})
	if err != nil {
		t.Fatal(err)
	}
	if reflect.TypeOf(fetchedGuildes) != reflect.TypeOf(dtos.GuildeListViewDTO{}) {
		t.Fatalf("Expected %s, got %s", reflect.TypeOf(dtos.GuildeListViewDTO{}), reflect.TypeOf(fetchedGuildes))
	}
}
func TestListGuildesWithFilters(t *testing.T) {
	uowm, ViewsManager, teardownTest := setupTest(t)
	defer teardownTest(t)

	saveBasicSamples(uowm)
	// Filter by name
	fetchedGuildes, err := ViewsManager.Guilde().List(views.ListGuildesViewOpts{Name: "min"})
	if err != nil {
		panic(err)
	}
	if reflect.TypeOf(fetchedGuildes) != reflect.TypeOf(dtos.GuildeListViewDTO{}) {
		t.Fatalf("Expected %s, got %s", reflect.TypeOf(dtos.GuildeListViewDTO{}), reflect.TypeOf(fetchedGuildes))
	}
	guildes := fetchedGuildes.Items
	if guildes[0].Name != "eminem" {
		t.Fatalf("Expected %s, got %s", "eminem", guildes[0].Name)
	}
	// Filter by exists
	fetchedGuildes, err = ViewsManager.Guilde().List(views.ListGuildesViewOpts{Exists: &[]bool{true}[0]})
	if err != nil {
		panic(err)
	}
	guildes = fetchedGuildes.Items
	if len(guildes) != 3 {
		t.Fatalf("Expected %d, got %d", 3, len(guildes))
	}
	// Filter by validated
	fetchedGuildes, err = ViewsManager.Guilde().List(views.ListGuildesViewOpts{Validated: &[]bool{false}[0]})
	if err != nil {
		panic(err)
	}
	guildes = fetchedGuildes.Items
	if len(guildes) != 1 {
		t.Fatalf("Expected %d, got %d", 1, len(guildes))
	}
	// Filter by active
	fetchedGuildes, err = ViewsManager.Guilde().List(views.ListGuildesViewOpts{Active: &[]bool{true}[0]})
	if err != nil {
		panic(err)
	}
	guildes = fetchedGuildes.Items
	if len(guildes) != 1 {
		t.Fatalf("Expected %d, got %d", 1, len(guildes))
	}
	// Filter by creation date
	after := time.Date(2020, 12, 1, 0, 0, 0, 0, time.UTC)
	before := time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)
	fetchedGuildes, err = ViewsManager.Guilde().List(views.ListGuildesViewOpts{CreationDateSince: after, CreationDateUntil: before})
	if err != nil {
		panic(err)
	}
	guildes = fetchedGuildes.Items
	if len(guildes) != 1 {
		t.Fatalf("Expected %d, got %d", 1, len(guildes))
	}
	// Test multiple filters
	fetchedGuildes, err = ViewsManager.Guilde().List(views.ListGuildesViewOpts{Name: "min", Exists: &[]bool{true}[0], Validated: &[]bool{true}[0], Active: &[]bool{true}[0], CreationDateSince: after})
	if err != nil {
		panic(err)
	}
	guildes = fetchedGuildes.Items
	if len(guildes) != 1 {
		t.Fatalf("Expected %d, got %d", 1, len(guildes))
	}
	// Test multiple filters 2
	fetchedGuildes, err = ViewsManager.Guilde().List(views.ListGuildesViewOpts{Exists: &[]bool{true}[0], Validated: &[]bool{true}[0], Active: &[]bool{true}[0], CreationDateSince: after, CreationDateUntil: before})
	if err != nil {
		panic(err)
	}
	guildes = fetchedGuildes.Items
	if len(guildes) != 0 {
		t.Fatalf("Expected %d, got %d", 0, len(guildes))
	}

}
