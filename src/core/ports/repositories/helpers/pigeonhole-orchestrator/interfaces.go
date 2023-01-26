package helpers

import (
	types "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
)

type SingleOperation[T any, K any] interface {
	work(*T) (*types.RepositoryDTO[K], error)
}

type MultipleOperation[T any, K any] interface {
	work(*T) (map[string]*types.RepositoryDTO[K], error)
}
