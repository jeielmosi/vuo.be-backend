package funcs

import (
	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	shorten_bulk "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/interfaces"
	repositories "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
)

func NewUnlockFunc(hash string) func(*shorten_bulk.ShortenBulkRepository) (*repositories.RepositoryDTO[entities.ShortenBulkEntity], error) {
	return func(repository *shorten_bulk.ShortenBulkRepository) (
		*repositories.RepositoryDTO[entities.ShortenBulkEntity],
		error,
	) {
		err := (*repository).Unlock(hash)
		return nil, err
	}
}
