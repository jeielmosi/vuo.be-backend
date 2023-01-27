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

func postAtOldHash(
	repository *shorten_bulk.ShortenBulkRepository,
	dto *repositories.RepositoryDTO[entities.ShortenBulkEntity],
) (string, error) {
	const OLDESTS_SIZE uint = 101

	mp, err := (*repository).GetOldests(OLDESTS_SIZE)
	if err != nil {
		return "", err
	}

	arr := helpers.MapToSlice(mp)

	for size := len(arr); size > 0; size-- {
		idx := rand.Intn(size)
		last := size - 1
		if idx != last {
			arr[idx], arr[last] = arr[last], arr[idx]
		}

		hash := arr[last].Key
		dto := arr[last].Value

		err := (*repository).Post(hash, *dto)
		if err == nil {
			return hash, nil
		}
	}

	return "", errors.New("Internal error: Hash not found")
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
