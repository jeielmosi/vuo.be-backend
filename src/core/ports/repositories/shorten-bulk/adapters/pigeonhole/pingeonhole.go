package adapters

import (
	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	helpers "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/helpers/pigeonhole-orchestrator"
	funcs "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/adapters/pigeonhole/funcs"
	shorten_bulk "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/interfaces"
	repositories "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
)

type PigeonholeShortenBulkRepository struct {
	orchestrator *helpers.PigeonholeOrchestrator[shorten_bulk.ShortenBulkRepository, entities.ShortenBulkEntity]
}

func (p *PigeonholeShortenBulkRepository) Get(hash string) (*repositories.RepositoryDTO[entities.ShortenBulkEntity], error) {
	fn := funcs.NewGetFunc(hash)
	res, err := p.orchestrator.ExecuteSingleFunc(fn)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *PigeonholeShortenBulkRepository) GetOldests(size uint) (map[string]*repositories.RepositoryDTO[entities.ShortenBulkEntity], error) {
	fn := funcs.NewGetOldestsFunc(size)
	res, err := p.orchestrator.ExecuteMultipleFunc(fn)
	if err != nil {
		return nil, err
	}

	ans := map[string]*repositories.RepositoryDTO[entities.ShortenBulkEntity]{}
	for key, value := range res {
		ans[key] = value
	}

	return ans, nil
}

func (p *PigeonholeShortenBulkRepository) Post(hash string, dto repositories.RepositoryDTO[entities.ShortenBulkEntity]) error {
	fn := funcs.NewPostFunc(hash, dto)
	_, err := p.orchestrator.ExecuteSingleFunc(fn)

	return err
}

func (p *PigeonholeShortenBulkRepository) IncrementClicks(hash string) error {
	fn := funcs.NewIncrementClicksFunc(hash)
	_, err := p.orchestrator.ExecuteSingleFunc(fn)
	return err
}

func NewPigeonholeShortenBulkRepository(repos *[]*shorten_bulk.ShortenBulkRepository) (
	shorten_bulk.ShortenBulkRepository,
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
