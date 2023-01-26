package wrappers

import (
	shorten_bulk "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/interfaces"
	ports "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
)

type IncrementClicksWrapper struct {
	hash string
}

func NewIncrementClicksWrapper(hash string) *IncrementClicksWrapper {
	return &IncrementClicksWrapper{
		hash,
	}
}

func (this *IncrementClicksWrapper) work(repository *shorten_bulk.ShortenBulkRepository) (
	*ports.RepositoryDTO[any],
	error,
) {
	err := (*repository).IncrementClicks(this.hash)
	return nil, err
}
