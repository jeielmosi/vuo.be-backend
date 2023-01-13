package adapters

type MultipleOperationFunction[T any] interface {
	work(*ShortenBulkRepository, <-chan *ShortenBulkRepository, chan<- T)
}
