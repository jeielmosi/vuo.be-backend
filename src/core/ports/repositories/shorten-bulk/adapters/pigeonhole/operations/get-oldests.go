package operations

import (
	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	shorten_bulk "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/interfaces"
	repositories "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
)

func NewGetOldestsOperation(size uint) func(*shorten_bulk.ShortenBulkRepository) (map[string]*repositories.RepositoryDTO[entities.ShortenBulkEntity], error) {
	return func(repository *shorten_bulk.ShortenBulkRepository) (
		map[string]*repositories.RepositoryDTO[entities.ShortenBulkEntity],
		error,
	) {
		return (*repository).GetOldests(size)
	}
}
