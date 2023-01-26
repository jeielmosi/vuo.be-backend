package shorten_bulk

import (
	domain "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	types "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
)

type ShortenBulkDTO = types.RepositoryDTO[domain.ShortenBulkEntity]

func NewShortenBulkDTO(
	entity *domain.ShortenBulkEntity,
	lockedAt *string,
	updatedAt string,
) (*ShortenBulkDTO, error) {
	return types.NewRepositoryDTO(entity, lockedAt, updatedAt)
}
