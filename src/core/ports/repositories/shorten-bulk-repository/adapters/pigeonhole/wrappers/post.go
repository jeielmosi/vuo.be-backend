package wrappers

import (
	repositories "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk-repository"
	ports "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
)

type PostWrapper struct {
	hash string
}

func NewPostWrapper(hash string) *PostWrapper {
	return &PostWrapper{
		hash,
	}
}

func (this *PostWrapper) work(repository *repositories.ShortenBulkRepository) (
	*ports.RepositoryDTO[any],
	error,
) {
	err := (*repository).Post(this.hash)
	return nil, err
}
