package adapters

import (
	domain "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	helpers "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/helpers/pigeonhole-orchestrator"
	wrappers "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/adapters/pigeonhole/wrappers"
	shorten_bulk "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/interfaces"
)

type Orchestrator = helpers.PigeonholeOrchestrator[shorten_bulk.ShortenBulkRepository, domain.ShortenBulkEntity]

func NewOrchestrator(repos *[]*shorten_bulk.ShortenBulkRepository) (
	*Orchestrator,
	error,
) {
	return helpers.
		NewPigeonholeOrchestrator[shorten_bulk.ShortenBulkRepository, domain.ShortenBulkEntity](repos)
}

type PigeonholeShortenBulkRepository struct {
	orchestrator *Orchestrator
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

func NewPigeonholeShortenBulkRepository(repos *[]*shorten_bulk.ShortenBulkRepository) (
	*PigeonholeShortenBulkRepository,
	error,
) {
	orchestrator, err := NewOrchestrator(repos)
	if err != nil {
		return nil, err
	}

	return &PigeonholeShortenBulkRepository{
		orchestrator,
	}, nil
}
