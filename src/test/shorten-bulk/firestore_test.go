package shorten_bulk_test

import (
	"math/rand"
	"os"
	"testing"
	"time"

	config "github.com/jei-el/vuo.be-backend/src/config"
	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	repository_helpers "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/helpers"
	firestore_shorten_bulk "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/adapters/firestore"
	shorten_bulk "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/interfaces"
	repositories "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
	shorten_bulk_helpers_test "github.com/jei-el/vuo.be-backend/src/test/shorten-bulk/helpers"
	randutil "go.step.sm/crypto/randutil"
)

// TODO: Move to firestore folder
const (
	getHash             = "get"
	postHash            = "post"
	lockHash            = "lock"
	incrementClicksHash = "increment_clicks"
)

func getFirestore() shorten_bulk.ShortenBulkRepository {
	config.Load()
	envName := os.Getenv("TEST_ENV")
	return firestore_shorten_bulk.NewShortenBulkFirestore(envName)
}

func getExpectedGetDTO() (
	*repositories.RepositoryDTO[entities.ShortenBulkEntity],
	error,
) {
	timestamp := "0000-01-01T00:00:00Z"
	timeObj, err := repository_helpers.NewTimeFrom10NanosecondsString(&timestamp)
	if err != nil {
		return nil, err
	}
	locked := false
	entity := entities.NewShortenBulkEntity("firebase.google.com", 0)

	exp := &repositories.RepositoryDTO[entities.ShortenBulkEntity]{
		Entity:    entity,
		Locked:    locked,
		CreatedAt: *timeObj,
		UpdatedAt: *timeObj,
	}

	return exp, nil
}

func getRandomDTO() (
	*repositories.RepositoryDTO[entities.ShortenBulkEntity],
	error,
) {
	rnd := rand.New(
		rand.NewSource(
			time.Now().UTC().UnixNano(),
		),
	)

	lock := (rnd.Intn(2) == 1)
	clicks := rnd.Int63()
	url, err := randutil.ASCII(1009)
	if err != nil {
		return nil, err
	}

	entity := entities.NewShortenBulkEntity(url, clicks)
	exp := repositories.NewRepositoryDTO(entity, lock)

	return exp, nil
}

func TestGet(t *testing.T) {
	firestore := getFirestore()

	exp, err := getExpectedGetDTO()
	if err != nil {
		t.Errorf("Creating a time: %s", err.Error())
	}

	shorten_bulk_helpers_test.TestGet(getHash, firestore, exp, t)
}

func TestGetOldest(t *testing.T) {
	firestore := getFirestore()

	dto, err := getExpectedGetDTO()
	if err != nil {
		t.Errorf("Test error at creating a time: %s", err.Error())
	}

	shorten_bulk_helpers_test.TestGetOldest(firestore, dto, t)
}

func TestIncrementClicks(t *testing.T) {
	firestore := getFirestore()
	shorten_bulk_helpers_test.TestIncrementClicks(incrementClicksHash, firestore, t)
}

func TestPost(t *testing.T) {
	dto, err := getRandomDTO()
	if err != nil {
		t.Errorf("Error creating a dto: %s", err.Error())
	}

	firestore := getFirestore()
	shorten_bulk_helpers_test.TestPost(postHash, firestore, dto, t)
}

func TestLockUnlock(t *testing.T) {
	firestore := getFirestore()
	shorten_bulk_helpers_test.TestLockUnlock(lockHash, firestore, t)
}
