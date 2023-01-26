package wrappers

import (
	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	helpers "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/helpers/pigeonhole-orchestrator"
	shorten_bulk "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/interfaces"
	repositories "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
)

type GetOldestsWrapper struct {
	size uint
}

func (g *GetOldestsWrapper) work(repository *shorten_bulk.ShortenBulkRepository) (
	map[string]*repositories.RepositoryDTO[entities.ShortenBulkEntity],
	error,
) {
	return (*repository).GetOldests(g.size)
}

func NewGetOldestsWrapper(size uint) *helpers.SingleOperation[shorten_bulk.ShortenBulkRepository, entities.ShortenBulkEntity] {
	return &GetOldestsWrapper{size}
}
