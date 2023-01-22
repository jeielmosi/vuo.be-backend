package adapters

type IncrementClicksWrapper struct {
	hash string
}

func NewIncrementClicksWrapper(hash string) *IncrementClicksWrapper {
	return &IncrementClicksWrapper{
		hash,
	}
}

func (this *IncrementClicksWrapper) work (repository *ShortenBulkRepository) 
(
	*DatabaseDTO[K],
	error,
) {
	err :=  repository.IncrementClicks(this.hash)
	return nil, err
}
