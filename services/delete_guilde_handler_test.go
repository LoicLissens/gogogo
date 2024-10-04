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

func TestDeleteGuildeHandler(t *testing.T) {
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
	cmd := commands.DeleteGuildeCommand{
		Uuid: g.Uuid,
	}
	err = TestServiceManager.DeleteGuildeHandler(cmd)
	if err != nil {
		t.Fatal(err)
	}
	uow, close = TestServiceManager.UnitOfWorkManager.Start()
	_, err = uow.GuildeRepository().GetByUUID(g.Uuid)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
	close()
}
func TestDeleteGuildeHandlerNotFound(t *testing.T) {
	var expectedError customerrors.ErrorNotFound

	TestServiceManager, teardownTest := SetupTest(t)
	defer teardownTest(t)
	cmd := commands.DeleteGuildeCommand{Uuid: uuid.New()}
	err := TestServiceManager.DeleteGuildeHandler(cmd)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
	if err != nil && !errors.As(err, &expectedError) {
		t.Fatalf("Expected ErrorNotFound, got %s", reflect.TypeOf(err))
	}
}
