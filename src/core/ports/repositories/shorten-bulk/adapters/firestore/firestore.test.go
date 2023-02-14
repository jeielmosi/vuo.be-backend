package firestore_shorten_bulk

import (
	"errors"
	"math/rand"
	"os"
	"reflect"
	"testing"

	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	repository_helpers "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/helpers"
	shorten_bulk "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/interfaces"
	repositories "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
	randutil "go.step.sm/crypto/randutil"
)

// TODO: Move to interface, the test is generic to interface
const (
	getHash  = "get"
	postHash = "post"
)

func getFirestore() shorten_bulk.ShortenBulkRepository {
	envName := os.Getenv("TEST_ENV")
	return NewShortenBulkFirestore(envName)
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
	clicks := uint64(rand.Int63())
	url, err := randutil.ASCII(1009)
	if err != nil {
		return nil, err
	}

	entity := entities.NewShortenBulkEntity(url, clicks)
	exp := repositories.NewRepositoryDTO(entity, false)

	return exp, nil
}

func TestGet(t *testing.T) {
	firestore := getFirestore()

	res, err := firestore.Get(getHash)
	if err != nil {
		t.Errorf("Get element at hash 'test': %s", err.Error())
	}

	exp, err := getExpectedGetDTO()
	if err != nil {
		t.Errorf("Creating a time: %s", err.Error())
	}

	if !reflect.DeepEqual(res, exp) {
		t.Errorf("Result and expected are different")
	}
}

func TestGetOldest(t *testing.T) {
	firestore := getFirestore()

	res, err := firestore.GetOldest(1)
	if err != nil {
		t.Errorf("Test error at get oldest element at hash 'test': %s", err.Error())
	}

	dto, err := getExpectedGetDTO()
	if err != nil {
		t.Errorf("Test error at creating a time: %s", err.Error())
	}

	exp := map[string]*repositories.RepositoryDTO[entities.ShortenBulkEntity]{
		"get": dto,
	}

	if !reflect.DeepEqual(res, exp) {
		t.Errorf("Test error at compare the get result")
	}
}

// TODO: create a generic func for other methods
func createTest(
	executor func(shorten_bulk.ShortenBulkRepository) error,
	update func(*repositories.RepositoryDTO[entities.ShortenBulkEntity]) (
		*repositories.RepositoryDTO[entities.ShortenBulkEntity],
		error,
	),
) error {
	firestore := getFirestore()

	prev, err := firestore.Get(postHash)
	if err != nil {
		return err
	}

	err = executor(firestore)
	if err != nil {
		return err
	}

	exp, err := update(prev)
	if err != nil {
		return err
	}

	pos, err := firestore.Get(postHash)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(pos, exp) {
		return errors.New("Result and expected are different")
	}

	return nil
}

func TestIncrementClicks(t *testing.T) {
	executor := func(firestore shorten_bulk.ShortenBulkRepository) error {
		err := firestore.IncrementClicks(postHash)
		if err != nil {
			return err
		}

		return nil
	}

	update := func(prev *repositories.RepositoryDTO[entities.ShortenBulkEntity]) (
		*repositories.RepositoryDTO[entities.ShortenBulkEntity],
		error,
	) {
		if prev == nil {
			return nil, errors.New("No element available")
		}
		prev.Entity.Clicks += 1
		return prev, nil
	}

	err := createTest(executor, update)
	if err != nil {
		t.Errorf("Increment Clicks error: %s", err.Error())
	}
}

func TestPost(t *testing.T) {
	dto, err := getRandomDTO()
	if err != nil {
		t.Errorf("Error creating a dto: %s", err.Error())
	}

	executor := func(firestore shorten_bulk.ShortenBulkRepository) error {
		err := firestore.Post(postHash, *dto)
		if err != nil {
			return err
		}

		return nil
	}

	update := func(_ *repositories.RepositoryDTO[entities.ShortenBulkEntity]) (
		*repositories.RepositoryDTO[entities.ShortenBulkEntity],
		error,
	) {
		return dto, nil
	}

	err = createTest(executor, update)
	if err != nil {
		t.Errorf("Post error: %s", err.Error())
	}
}

func TestLock(t *testing.T) {
	executor := func(firestore shorten_bulk.ShortenBulkRepository) error {
		err := firestore.Lock(postHash)
		if err != nil {
			return err
		}

		return nil
	}

	update := func(prev *repositories.RepositoryDTO[entities.ShortenBulkEntity]) (
		*repositories.RepositoryDTO[entities.ShortenBulkEntity],
		error,
	) {
		if prev == nil {
			return nil, errors.New("Empty dto")
		}

		if prev.Locked {
			return nil, errors.New("Already locked")
		}

		prev.Locked = true
		return prev, nil
	}

	err := createTest(executor, update)

	firestore := getFirestore()
	defer firestore.Unlock(postHash)

	if err != nil {
		t.Errorf("Post error: %s", err.Error())
	}
}

func TestUnlock(t *testing.T) {
	executor := func(firestore shorten_bulk.ShortenBulkRepository) error {
		err := firestore.Lock(postHash)
		if err != nil {
			return err
		}

		err = firestore.Unlock(postHash)
		if err != nil {
			return err
		}

		return nil
	}

	update := func(prev *repositories.RepositoryDTO[entities.ShortenBulkEntity]) (
		*repositories.RepositoryDTO[entities.ShortenBulkEntity],
		error,
	) {
		if prev == nil {
			return nil, errors.New("Empty dto")
		}

		return prev, nil
	}

	err := createTest(executor, update)

	if err != nil {
		t.Errorf("Post error: %s", err.Error())
	}
}
