package services

import (
	"errors"
	"jiva-guildes/domain/commands"
	customerrors "jiva-guildes/domain/custom_errors"
	"jiva-guildes/domain/models/guilde"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestDeleteGuildeHandler(t *testing.T) {
	TestServiceManager, teardownTest := SetupTest(t)
	defer teardownTest(t)
	uow := TestServiceManager.UnitOfWorkManager.Start()
	guilde := guilde.New("GUnit", "https://www.googleimage.com", "https://www.google.com")
	g, err := uow.GuildeRepository().Save(*guilde)
	if err != nil {
		t.Fatal(err)
	}
	cmd := commands.DeleteGuildeCommand{
		Uuid: g.Uuid,
	}
	err = TestServiceManager.DeleteGuildeHandler(cmd)
	if err != nil {
		t.Fatal(err)
	}
	_, err = uow.GuildeRepository().GetByUUID(g.Uuid)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
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
