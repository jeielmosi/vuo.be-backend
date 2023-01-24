package ports

type SingleOperationFunction[T any, K any] interface {
	work(*T) (*RepositoryDTO[K], error)
}

type MultipleOperationFunction[T any, K any] interface {
	work(*T) (map[string]*RepositoryDTO[K], error)
}
