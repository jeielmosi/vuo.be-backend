package adapters

import (
	domain "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	helpers "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/helpers/pigeonhole-orchestrator"
	wrappers "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/adapters/pigeonhole/wrappers"
	shorten_bulk "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/interfaces"
)

type PigeonholeShortenBulkRepository struct {
	orchestrator *helpers.PigeonholeOrchestrator[shorten_bulk.ShortenBulkRepository, domain.ShortenBulkEntity]
}

func (p *PigeonholeShortenBulkRepository) Get(hash string) (*domain.ShortenBulkEntity, error) {
	worker := wrappers.NewGetWrapper(hash)
	res, err := p.orchestrator.ExecuteSingleOperation(worker)
	if err != nil {
		return nil, err
	}

	return res.Entity, nil
}

func (p *PigeonholeShortenBulkRepository) GetOldest(size uint) (map[string]*domain.ShortenBulkEntity, error) {
	worker := wrappers.NewGetOldestsWrapper(size)
	res, err := p.orchestrator.ExecuteMultipleOperation(worker)
	if err != nil {
		return nil, err
	}

	ans := map[string]*domain.ShortenBulkEntity{}
	for key, value := range res {
		ans[key] = value.Entity
	}

	return ans, nil
}

func (p *PigeonholeShortenBulkRepository) Post(hash string) error {
	worker := wrappers.NewPostWrapper(hash)
	_, err := p.orchestrator.ExecuteSingleOperation(worker)
	return err
}

func (p *PigeonholeShortenBulkRepository) IncrementClicks(hash string) error {
	worker := wrappers.NewIncrementClicksWrapper(hash)
	_, err := p.orchestrator.ExecuteSingleOperation(worker)
	return err
}

func NewPigeonholeShortenBulkRepository(repos *[]*shorten_bulk.ShortenBulkRepository) (
	*PigeonholeShortenBulkRepository,
	error,
) {
	orchestrator, err := helpers.NewPigeonholeOrchestrator[shorten_bulk.ShortenBulkRepository, entities.ShortenBulkEntity](repos)
	if err != nil {
		return nil, err
	}

	return &PigeonholeShortenBulkRepository{
		orchestrator,
	}, nil
}
