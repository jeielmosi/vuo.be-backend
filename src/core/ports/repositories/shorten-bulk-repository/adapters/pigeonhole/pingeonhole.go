package adapters

import (
	domain "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	pigeonhole_orchestrator "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/pigeonhole-orchestrator"
	repositories "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk-repository"
	wrappers "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk-repository/adapters/pigeonhole/wrappers"
)

type PigeonholeShortenBulkRepository struct {
	orchestrator pigeonhole_orchestrator.PigeonholeOrchestrator[repositories.ShortenBulkRepository, domain.ShortenBulkEntity]
}

func (this *PigeonholeShortenBulkRepository) Get(hash string) (*domain.ShortenBulkEntity, error) {
	worker := wrappers.NewGetWrapper(hash)
	res, err := this.orchestrator.SingleOperation(worker)
	if err != nil {
		return nil, err
	}

	return res.Entity, nil
}

func (this *PigeonholeShortenBulkRepository) GetOldest(size uint) (map[string]*domain.ShortenBulkEntity, error) {
	worker := wrappers.NewGetOldestsWrapper(size)
	res, err := this.orchestrator.MultipleOperation(worker)
	if err != nil {
		return nil, err
	}

	ans := map[string]*domain.ShortenBulkEntity{}
	for key, value := range res {
		ans[key] = value.Entity
	}

	return ans, nil
}

func (this *PigeonholeShortenBulkRepository) Post(hash string) error {
	worker := wrappers.NewPostWrapper(hash)
	_, err := this.orchestrator.SingleOperation(worker)
	return err
}

func (this *PigeonholeShortenBulkRepository) IncrementClicks(hash string) error {
	worker := wrappers.NewIncrementClicksWrapper(hash)
	_, err := this.orchestrator.SingleOperation(worker)
	return err
}

func NewPigeonholeShortenBulkRepository(repositories []*repositories.ShortenBulkRepository) *PigeonholeShortenBulkRepository {
	return &PigeonholeShortenBulkRepository{
		orchestrator: pigeonhole_orchestrator.
			NewPigeonholeOrchestrator[repositories.ShortenBulkRepository, domain.ShortenBulkEntity](repositories),
	}
}
