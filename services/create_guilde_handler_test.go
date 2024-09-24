package services

import (
	"jiva-guildes/adapters"
	"jiva-guildes/adapters/db"
	"jiva-guildes/domain/commands"
	"jiva-guildes/settings"
	"testing"

	"github.com/google/uuid"
)

// TODO Try to use test suite (see:https://medium.com/nerd-for-tech/setup-and-teardown-unit-test-in-go-bd6fa1b785cd)

func setupTest(tb testing.TB) (ServiceManager, func(tb testing.TB)) {
	var connectionPool = db.MountDB(settings.AppSettings.DATABASE_URI)
	var UnitOfWorkManager = adapters.NewUnitOfWorkManager(connectionPool)
	var TestServiceManager = ServiceManager{UnitOfWorkManager: &UnitOfWorkManager}

	return TestServiceManager, func(tb testing.TB) {
		db.Teardown(connectionPool)
	}
}
func TestCreateGuildeHandler(t *testing.T) {
	TestServiceManager, teardownTest := setupTest(t)
	defer teardownTest(t)

	cmd := commands.CreateGuildeCommand{
		Name:     "GUnit",
		Img_url:  "https://www.googleimage.com",
		Page_url: "https://www.google.com",
	}
	g, err := TestServiceManager.CreateGuildeHandler(cmd)
	if err != nil {
		t.Fatal(err)
	}
	if g.Name != cmd.Name {
		t.Fatalf("Expected %s, got %s", cmd.Name, g.Name)
	}
	if g.Img_url != cmd.Img_url {
		t.Fatalf("Expected %s, got %s", cmd.Img_url, g.Img_url)
	}
	if g.Page_url != cmd.Page_url {
		t.Fatalf("Expected %s, got %s", cmd.Page_url, g.Page_url)
	}
	if g.Uuid == uuid.Nil {
		t.Fatalf("Expected uuid, got nil")
	}
	if g.Created_at.IsZero() { //TODO can go further and freeze
		t.Fatalf("Expected created_at, got zero value")
	}
	if g.Updated_at.IsZero() {
		t.Fatalf("Expected created_at, got zero value")
	}
}
