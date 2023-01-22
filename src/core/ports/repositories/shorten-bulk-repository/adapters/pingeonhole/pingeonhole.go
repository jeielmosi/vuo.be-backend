package adapters

type PigeonholeShortenBulkRepository struct {
	orchestrator PigeonholeOrchestrator[ShortenBulkRepository, ShortenBulkEntity]
}

func (this *PigeonholeShortenBulkRepository) GetByHash(hash string) (*ShortenBulkEntity, error) {
	worker := NewGetHashWrapper(hash)
	res, err := this.orchestrator.SingleOperation(worker)
	if err != nil {
		return nil, err
	}

	return res.Entity, nil
}

func (this *PigeonholeShortenBulkRepository) GetOldest(size uint) (map[string]*ShortenBulkEntity, error) {
	worker := NewGetOldestsWrapper(size)
	res, err := this.orchestrator.MultipleOperationFunction(worker)
	if err != nil {
		return nil, err
	}

	ans := map[string]*ShortenBulkEntity{}
	for key, value := range res {
		ans[key] = value.Entity
	}

	return ans
}

func (this *PigeonholeShortenBulkRepository) PostAtHash(hash string) error {
	worker := NewPostAtHashWrapperWrapper(hash)
	_, err := this.orchestrator.SingleOperation(worker)
	return err
}

func (this *PigeonholeShortenBulkRepository) IncrementClicks(hash string) error {
	worker := NewIncrementClicksWrapper(hash)
	_, err := this.orchestrator.SingleOperation(worker)
	return err
}

func NewPigeonholeShortenBulkRepository(repositories []*ShortenBulkRepository) *PigeonholeShortenBulkRepository {
	return &PigeonholeShortenBulkRepository{
		orchestrator: NewPigeonholeOrchestrator[ShortenBulkRepository, ShortenBulkEntity](repositories),
	}
}
