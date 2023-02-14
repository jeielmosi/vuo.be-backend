package shorten_bulk_helpers_test

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"testing"

	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	shorten_bulk "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/interfaces"
	types "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/types"
	repositories "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
)

func TestGet(
	hash string,
	repo shorten_bulk.ShortenBulkRepository,
	exp *repositories.RepositoryDTO[entities.ShortenBulkEntity],
	t *testing.T,
) {

	res, err := repo.Get(hash)
	if err != nil {
		t.Errorf("Get element at hash 'test': %s", err.Error())
	}

	if !reflect.DeepEqual(res, exp) {
		t.Errorf("Result and expected are different")
	}
}

func TestGetOldest(
	repo shorten_bulk.ShortenBulkRepository,
	dto *repositories.RepositoryDTO[entities.ShortenBulkEntity],
	t *testing.T,
) {
	res, err := repo.GetOldest(1)
	if err != nil {
		t.Errorf("Test error at get oldest element at hash 'test': %s", err.Error())
	}

	exp := map[string]*repositories.RepositoryDTO[entities.ShortenBulkEntity]{
		"get": dto,
	}

	if !reflect.DeepEqual(res, exp) {
		fmt.Println(res)
		fmt.Println(exp)
		t.Errorf("Test error at compare the get result")
	}
}

func createTest(
	hash string,
	repo shorten_bulk.ShortenBulkRepository,
	executor func(shorten_bulk.ShortenBulkRepository) error,
	update func(*repositories.RepositoryDTO[entities.ShortenBulkEntity]) (
		*repositories.RepositoryDTO[entities.ShortenBulkEntity],
		error,
	),
) error {
	prev, err := repo.Get(hash)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = executor(repo)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	exp, err := update(prev)
	if err != nil {
		return err
	}

	res, err := repo.Get(hash)
	if err != nil {
		return err
	}

	if res == nil && exp == nil {
		return nil
	}

	if res == nil || exp == nil {
		if res == nil {
			return errors.New("Result and expected are different: result is nil, expected is not")
		} else {
			return errors.New("Result and expected are different: expected is nil, result is nil")
		}
	}

	resFlatten := types.NewShortenBulkFlattenDTO(*res)
	expFlatten := types.NewShortenBulkFlattenDTO(*exp)

	if !reflect.DeepEqual(resFlatten, expFlatten) {
		log.Println("res:", resFlatten)
		log.Println("exp:", expFlatten)

		return errors.New("Result and expected are different")
	}

	return nil
}

func TestIncrementClicks(
	hash string,
	repo shorten_bulk.ShortenBulkRepository,
	t *testing.T,
) {
	executor := func(repo shorten_bulk.ShortenBulkRepository) error {
		err := repo.IncrementClicks(hash)
		if err != nil {
			log.Println(err.Error())
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

	err := createTest(hash, repo, executor, update)
	if err != nil {
		t.Errorf("Increment Clicks error: %s", err.Error())
	}
}

func TestPost(
	hash string,
	repo shorten_bulk.ShortenBulkRepository,
	dto *repositories.RepositoryDTO[entities.ShortenBulkEntity],
	t *testing.T,
) {
	executor := func(firestore shorten_bulk.ShortenBulkRepository) error {
		err := firestore.Post(hash, *dto)
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

	err := createTest(hash, repo, executor, update)
	if err != nil {
		t.Errorf("Post error: %s", err.Error())
	}
}

func TestLockUnlock(
	hash string,
	repo shorten_bulk.ShortenBulkRepository,
	t *testing.T,
) {
	executorLock := func(firestore shorten_bulk.ShortenBulkRepository) error {
		return firestore.Lock(hash)
	}

	updateLock := func(prev *repositories.RepositoryDTO[entities.ShortenBulkEntity]) (
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

	err := createTest(hash, repo, executorLock, updateLock)
	if err != nil {
		t.Errorf("Lock error: %s", err.Error())
		return
	}

	executorUnlock := func(firestore shorten_bulk.ShortenBulkRepository) error {
		return firestore.Unlock(hash)
	}

	updateUnlock := func(prev *repositories.RepositoryDTO[entities.ShortenBulkEntity]) (
		*repositories.RepositoryDTO[entities.ShortenBulkEntity],
		error,
	) {
		if prev == nil {
			return nil, errors.New("Empty dto")
		}

		if !prev.Locked {
			return nil, errors.New("Already unlocked")
		}

		prev.Locked = false
		return prev, nil
	}

	err = createTest(hash, repo, executorUnlock, updateUnlock)
	if err != nil {
		t.Errorf("Unlock error: %s", err.Error())
	}
}
