package pigeonhole_shorten_bulk

import (
	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	helpers "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/helpers/pigeonhole-orchestrator"
	funcs "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/helpers/pigeonhole/funcs"
	shorten_bulk "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/interfaces"
	repositories "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
)

type PigeonholeShortenBulkRepository struct {
	orchestrator *helpers.PigeonholeOrchestrator[
		shorten_bulk.ShortenBulkRepository,
		entities.ShortenBulkEntity,
	]
}

func (p *PigeonholeShortenBulkRepository) Get(hash string) (
	*repositories.RepositoryDTO[entities.ShortenBulkEntity],
	error,
) {
	fn := funcs.NewGetFunc(hash)
	dto, err, _ := p.orchestrator.ExecuteSingleFn(fn)
	return dto, err
}

func (p *PigeonholeShortenBulkRepository) GetOldests(size uint) (
	map[string]*repositories.RepositoryDTO[entities.ShortenBulkEntity],
	error,
) {
	fn := funcs.NewGetOldestsFunc(size)
	res, err, _ := p.orchestrator.ExecuteMultipleFn(fn)
	if err != nil {
		return nil, err
	}

	ans := map[string]*repositories.RepositoryDTO[entities.ShortenBulkEntity]{}
	for key, value := range res {
		ans[key] = value
	}

	return ans, nil
}

func (p *PigeonholeShortenBulkRepository) Post(
	hash string,
	dto repositories.RepositoryDTO[entities.ShortenBulkEntity],
) (
	error,
	<-chan int,
) {
	fn := funcs.NewPostFunc(hash, dto)
	_, err, resCh := p.orchestrator.ExecuteSingleFn(fn)

	return err, resCh
}

func (p *PigeonholeShortenBulkRepository) UndoPost(
	hash string,
	dto repositories.RepositoryDTO[entities.ShortenBulkEntity],
	idxsCh <-chan int,
) (
	error,
	<-chan int,
) {
	fn := funcs.NewPostFunc(hash, dto)
	idxCh1 := make(chan int, len(idxsCh))
	addMp := map[int]bool{}

	for idx := range idxsCh {
		addMp[idx] = true
	}
	_, err, resCh := p.orchestrator.ExecuteSingleFnWithCh(fn, idxCh1)

	size := len(addMp) - len(resCh)
	for idx := range resCh {
		addMp[idx] = false
	}

	errCh := make(chan int, size)
	for idx, add := range addMp {
		if add {
			errCh <- idx
		}
	}
	close(errCh)

	return err, errCh
}

func (p *PigeonholeShortenBulkRepository) IncrementClicks(hash string) (
	error,
	<-chan int,
) {
	fn := funcs.NewIncrementClicksFunc(hash)
	_, err, resCh := p.orchestrator.ExecuteSingleFn(fn)
	return err, resCh
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
