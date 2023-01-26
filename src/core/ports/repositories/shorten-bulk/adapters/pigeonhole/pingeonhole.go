package adapters

import (
	domain "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	helpers "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/helpers/pigeonhole-orchestrator"
	operations "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/adapters/pigeonhole/operations"
	shorten_bulk "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/interfaces"
)

type PigeonholeShortenBulkRepository struct {
	orchestrator *helpers.PigeonholeOrchestrator[shorten_bulk.ShortenBulkRepository, domain.ShortenBulkEntity]
}

func (p *PigeonholeShortenBulkRepository) Get(hash string) (*domain.ShortenBulkEntity, error) {
	operation := operations.NewGetOperation(hash)
	res, err := p.orchestrator.ExecuteSingleOperation(operation)
	if err != nil {
		return nil, err
	}

	return res.Entity, nil
}

func (p *PigeonholeShortenBulkRepository) GetOldest(size uint) (map[string]*domain.ShortenBulkEntity, error) {
	operation := operations.NewGetOldestsOperation(size)
	res, err := p.orchestrator.ExecuteMultipleOperation(operation)
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
	operation := operations.NewPostOperation(hash)
	_, err := p.orchestrator.ExecuteSingleOperation(operation)
	return err
}

func (p *PigeonholeShortenBulkRepository) IncrementClicks(hash string) error {
	operation := operations.NewIncrementClicksOperation(hash)
	_, err := p.orchestrator.ExecuteSingleOperation(operation)
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
