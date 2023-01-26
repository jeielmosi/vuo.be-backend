package helpers

import (
	repositories "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
)

type SingleOperation[T any, K any] interface {
	work(*T) (*repositories.RepositoryDTO[K], error)
}

type MultipleOperation[T any, K any] interface {
	work(*T) (map[string]*repositories.RepositoryDTO[K], error)
}
