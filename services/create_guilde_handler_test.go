package services

import (
	"errors"
	"jiva-guildes/domain/commands"
	customerrors "jiva-guildes/domain/custom_errors"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

var (
	frozenTime = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	NowFunc    = time.Now
)

func TestCreateGuildeHandler(t *testing.T) {
	NowFunc = func() time.Time {
		return frozenTime
	}
	defer func() { NowFunc = time.Now }()

	TestServiceManager, teardownTest := SetupTest(t)
	defer teardownTest(t)

	cmd := commands.CreateGuildeCommand{
		Name:          "GUnit",
		Img_url:       "https://www.googleimage.com",
		Page_url:      "https://www.google.com",
		Exists:        true,
		Active:        &[]bool{true}[0],
		Creation_date: &[]time.Time{time.Now()}[0],
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
	if g.Created_at.IsZero() {
		t.Fatalf("Expected created_at, got zero value")
	} //TODO: fix this
	// if g.Created_at != frozenTime {
	// 	t.Fatalf("Expected created_at %s, got %s", frozenTime, g.Created_at)
	// }
	if g.Updated_at.IsZero() {
		t.Fatalf("Expected created_at, got zero value")
	}
	if *g.Active != *cmd.Active {
		t.Fatalf("Expected %t, got %t", *cmd.Active, *g.Active)
	}
	// if g.Creation_date != nil && cmd.Creation_date != nil && g.Creation_date.UTC() != cmd.Creation_date.UTC() {
	// 	t.Fatalf("Expected %s, got %s", cmd.Creation_date.UTC(), g.Creation_date.UTC())
	// }
	if g.Validated != true {
		t.Fatalf("Expected false, got %t", g.Validated)
	}
	if g.Exists != cmd.Exists {
		t.Fatalf("Expected %t, got %t", cmd.Exists, g.Exists)
	}
}
func TestCreateGuildeNotExistActive(t *testing.T) {
	var expectedError customerrors.ValueError

	TestServiceManager, teardownTest := SetupTest(t)
	defer teardownTest(t)

	cmd := commands.CreateGuildeCommand{
		Name:          "GUnit",
		Img_url:       "https://www.googleimage.com",
		Page_url:      "https://www.google.com",
		Exists:        false,
		Active:        &[]bool{true}[0],
		Creation_date: nil,
	}
	_, err := TestServiceManager.CreateGuildeHandler(cmd)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if err != nil && !errors.As(err, &expectedError) {
		t.Fatal("Expected ValueError, got", reflect.TypeOf(err))
	}
}
