package wrappers

import (
	repositories "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk-repository"
	"github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk-repository/types"
)

type GetWrapper struct {
	hash string
}

func NewGetWrapper(hash string) *GetWrapper {
	return &GetWrapper{
		hash,
	}
}

func (this *GetWrapper) work(repository *repositories.ShortenBulkRepository) (
	*types.ShortenBulkDTO,
	error,
) {

	return (*repository).Get(this.hash)
}
