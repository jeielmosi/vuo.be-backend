package usecases

import (
	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	shorten_bulk "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/interfaces"
)

func GetShortenBulkEntity(repository *shorten_bulk.ShortenBulkRepository, hash string) (
	*entities.ShortenBulkEntity,
	error,
) {
	res, err := (*repository).Get(hash)
	if err != nil {
		return nil, err
	}

	err = (*repository).IncrementClicks(hash)
	return res.Entity, err
}
