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
	*PigeonholeDTO[K],
	error,
) {
	err :=  repository.IncrementClicks(this.hash)
	return nil, err
}
