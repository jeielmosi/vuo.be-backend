package usecases

import (
	"errors"
	"math/rand"
	"time"

	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	helpers "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk/helpers"
	shorten_bulk "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/interfaces"
	repositories "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
)

const RECOVERY_SIZE int = 5
const HASH_SIZE uint = 7

func postAtNewHash(
	repository *shorten_bulk.ShortenBulkRepository,
	dto *repositories.RepositoryDTO[entities.ShortenBulkEntity],
) (string, error) {
	const TRY_SIZE int = 11

	for i := 0; i < TRY_SIZE; i++ {
		hash := helpers.NewRandomHash(HASH_SIZE)

		backup, err := (*repository).Get(hash)
		if err != nil {
			continue
		}

		err = (*repository).Post(hash, *dto)
		for r := 0; r < RECOVERY_SIZE; r++ {
			if err == nil {
				break
			}
			err = (*repository).Post(hash, *backup)
		}
		if err == nil {
			return hash, nil
		}
	}

	return "", errors.New("Internal error: Not found empty hash")
}

type KeyValue struct {
	Key   string
	Value *repositories.RepositoryDTO[entities.ShortenBulkEntity]
}

func postAtOldHash(
	repository *shorten_bulk.ShortenBulkRepository,
	dto *repositories.RepositoryDTO[entities.ShortenBulkEntity],
) (string, error) {
	const OLDESTS_SIZE uint = 101

	oldests, err := (*repository).GetOldests(OLDESTS_SIZE)
	if err != nil {
		return "", err
	}

	oldestsArray := make([]KeyValue, 0)
	for key, val := range oldests {
		oldestsArray = append(oldestsArray, KeyValue{
			Key:   key,
			Value: val,
		})
	}

	for size := len(oldestsArray); size > 0; size-- {
		idx := rand.Intn(size)
		last := size - 1
		if idx != last {
			oldestsArray[idx], oldestsArray[last] = oldestsArray[last], oldestsArray[idx]
		}

		hash := oldestsArray[last].Key
		dto := oldestsArray[last].Value

		err := (*repository).Post(hash, *dto)
		if err == nil {
			return hash, nil
		}
	}

	return "", errors.New("Internal error: Not found empty hash")
}

func PostShortenBulkEntity(
	repository *shorten_bulk.ShortenBulkRepository,
	entity entities.ShortenBulkEntity,
) (string, error) {
	now := time.Now()
	dto := repositories.NewRepositoryDTO(&entity, now, nil)

	res, err := postAtNewHash(repository, dto)
	if err == nil {
		return res, err
	}

	return postAtOldHash(repository, dto)
}
