package wrappers

import (
	shorten_bulk "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/interfaces"
	types "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/types"
)

type GetOldestsWrapper struct {
	size uint
}

func (this *GetOldestsWrapper) work(repository *shorten_bulk.ShortenBulkRepository) (
	map[string]*types.ShortenBulkDTO,
	error,
) {
	return (*repository).GetOldests(this.size)
}

func NewGetOldestsWrapper(size uint) *GetOldestsWrapper {
	return &GetOldestsWrapper{
		size,
	}
}
