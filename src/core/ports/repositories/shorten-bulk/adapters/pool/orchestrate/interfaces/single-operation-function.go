package adapters

type SingleOperationFunction[T any] interface {
	work(*ShortenBulkRepository, <-chan *ShortenBulkRepository, chan<- T)
}
