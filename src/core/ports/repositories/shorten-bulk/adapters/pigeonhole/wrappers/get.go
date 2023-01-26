package wrappers

import (
	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	helpers "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/helpers/pigeonhole-orchestrator"
	shorten_bulk "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/interfaces"
	repositories "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
)

type GetWrapper struct {
	hash string
}

func (g *GetWrapper) work(repository *shorten_bulk.ShortenBulkRepository) (
	*repositories.RepositoryDTO[entities.ShortenBulkEntity],
	error,
) {

	return (*repository).Get(g.hash)
}

func NewGetWrapper(hash string) *helpers.SingleOperation[shorten_bulk.ShortenBulkRepository, entities.ShortenBulkEntity] {
	return &GetWrapper{hash}
}
