package wrappers

import (
	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	helpers "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/helpers/pigeonhole-orchestrator"
	shorten_bulk "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/interfaces"
	repositories "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
)

type PostWrapper struct {
	hash string
}

func (p *PostWrapper) work(repository *shorten_bulk.ShortenBulkRepository) (
	*repositories.RepositoryDTO[entities.ShortenBulkEntity],
	error,
) {
	err := (*repository).Post(p.hash)
	return nil, err
}

func NewPostWrapper(hash string) *helpers.SingleOperation[shorten_bulk.ShortenBulkRepository, entities.ShortenBulkEntity] {
	return &PostWrapper{hash}
}
