package firestore_shorten_bulk

import (
	"os"
	"reflect"
	"testing"

	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	repository_helpers "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/helpers"
	shorten_bulk "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/interfaces"
	repositories "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
)

func getFirestore() shorten_bulk.ShortenBulkRepository {
	envName := os.Getenv("TEST_ENV")
	return NewShortenBulkFirestore(envName)
}

func TestGet(t *testing.T) {
	firestore := getFirestore()

	res, err := firestore.Get("test")
	if err != nil {
		t.Errorf("Test error at get element at hash 'test': %s", err.Error())
	}

	timestamp := "0000-01-01T00:00:00Z"
	timeObj, err := repository_helpers.NewTimeFrom10NanosecondsString(&timestamp)
	if err != nil {
		t.Errorf("Test error at creating a time: %s", err.Error())
	}
	locked := false
	entity := entities.NewShortenBulkEntity("firebase.google.com", 0)

	exp := &repositories.RepositoryDTO[entities.ShortenBulkEntity]{
		Entity:    entity,
		Locked:    locked,
		CreatedAt: *timeObj,
		UpdatedAt: *timeObj,
	}

	if !reflect.DeepEqual(res, exp) {
		t.Errorf("Test error at compare the get result")
	}
}
