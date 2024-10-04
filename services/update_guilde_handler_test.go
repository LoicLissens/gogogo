package services

import (
	"errors"
	"jiva-guildes/domain/commands"
	customerrors "jiva-guildes/domain/custom_errors"
	"jiva-guildes/domain/models/guilde"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestUpdateGuildeHandler(t *testing.T) {
	TestServiceManager, teardownTest := SetupTest(t)
	defer teardownTest(t)
	uow, close := TestServiceManager.UnitOfWorkManager.Start()
	guilde, err := guilde.New(guilde.GuildeOptions{Name: "GUnit",
		Img_url:       "https://www.googleimage.com",
		Page_url:      "https://www.google.com",
		Exists:        true,
		Active:        &[]bool{true}[0],
		Creation_date: &[]time.Time{time.Now()}[0],
		Validated:     true,
	})
	if err != nil {
		t.Fatal(err)
	}
	g, err := uow.GuildeRepository().Save(*guilde)
	if err != nil {
		t.Fatal(err)
	}
	close()
	img := "img"
	page := "page"
	exist := false
	active := false
	creation := time.Now().UTC()
	validated := false
	name := "name"
	cmd := commands.UpdateGuildeCommand{
		Uuid:         g.Uuid,
		Name:         name,
		Img_url:      img,
		Page_url:     page,
		Exists:       &exist,
		Active:       &active,
		CreationDate: creation,
		Validated:    &validated,
	}
	updatedGuide, err := TestServiceManager.UpdateGuildeHandler(cmd)
	if err != nil {
		t.Fatal(err)
	}
	if updatedGuide.Name != name {
		t.Fatalf("Expected %s, got %s", name, updatedGuide.Name)
	}
	if updatedGuide.Img_url != img {
		t.Fatalf("Expected %s, got %s", img, updatedGuide.Img_url)
	}
	if updatedGuide.Page_url != page {
		t.Fatalf("Expected %s, got %s", page, updatedGuide.Page_url)
	}
	if updatedGuide.Exists != exist {
		t.Fatalf("Expected %t, got %t", exist, updatedGuide.Exists)
	}
	if *updatedGuide.Active != active {
		t.Fatalf("Expected %t, got %t", active, *updatedGuide.Active)
	}
	if *updatedGuide.Creation_date != creation {
		t.Fatalf("Expected %s, got %s", creation, updatedGuide.Creation_date)
	}
	if updatedGuide.Validated != validated {
		t.Fatalf("Expected %t, got %t", validated, updatedGuide.Validated)
	}
}
func TestUpdateGuildeHandlerNotFound(t *testing.T) {
	var expectedError customerrors.ErrorNotFound

	TestServiceManager, teardownTest := SetupTest(t)
	defer teardownTest(t)
	cmd := commands.UpdateGuildeCommand{Uuid: uuid.New()}
	_, err := TestServiceManager.UpdateGuildeHandler(cmd)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
	if err != nil && !errors.As(err, &expectedError) {
		t.Fatalf("Expected ErrorNotFound, got %s", reflect.TypeOf(err))
	}
}
