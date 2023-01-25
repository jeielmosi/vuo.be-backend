package wrappers

import (
	domain "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	repositories "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk-repository"
	"github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
)

type GetOldestsWrapper struct {
	size uint
}

func NewGetOldestsWrapper(size uint) *GetOldestsWrapper {
	return &GetOldestsWrapper{
		size,
	}
}

func (this *GetOldestsWrapper) work(repository *repositories.ShortenBulkRepository) (
	map[string]*types.RepositoryDTO[domain.ShortenBulkEntity],
	error,
) {
	return (*repository).GetOldests(this.size)
}
