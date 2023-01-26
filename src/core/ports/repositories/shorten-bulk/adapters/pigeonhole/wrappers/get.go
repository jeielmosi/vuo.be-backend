package wrappers

import (
	shorten_bulk "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/interfaces"
	types "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/types"
)

type GetWrapper struct {
	hash string
}

func NewGetWrapper(hash string) *GetWrapper {
	return &GetWrapper{
		hash,
	}
}

func (this *GetWrapper) work(repository *shorten_bulk.ShortenBulkRepository) (
	*types.ShortenBulkDTO,
	error,
) {

	return (*repository).Get(this.hash)
}
