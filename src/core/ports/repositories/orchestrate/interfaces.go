package adapters

type SingleOperationFunction[T any, K any] interface {
	work(T) (PigeonholeDTO[K], error)
}

type MultipleOperationFunction[T any, K any] interface {
	work(T) (map[string]PigeonholeDTO[K], error)
}
