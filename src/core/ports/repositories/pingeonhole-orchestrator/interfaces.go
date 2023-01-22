package ports

type SingleOperationFunction[T any, K any] interface {
	work(*T) (*DatabaseDTO[K], error)
}

type MultipleOperationFunction[T any, K any] interface {
	work(*T) (map[string]*DatabaseDTO[K], error)
}
