package services

import (
	"jiva-guildes/domain/commands"
	"testing"

	"github.com/google/uuid"
)

func TestCreateGuildeHandler(t *testing.T) {
	TestServiceManager, teardownTest := SetupTest(t)
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
