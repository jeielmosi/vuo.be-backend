package helpers

import (
	"errors"
	"sync"

	helpers "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/helpers"
	repositories "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
)

type PigeonholeOrchestrator[T any, K any] struct {
	worksSize    int
	repositories *[]*T
}

func (o *PigeonholeOrchestrator[T, K]) ExecuteSingleOperation(
	singleOperation func(*T) (*repositories.RepositoryDTO[K], error),
) (*repositories.RepositoryDTO[K], error) {
	if len(*o.repositories) < o.worksSize {
		return nil, errors.New("Internal error: Not enough repositories")
	}
	randomRepositories := helpers.NewRandomChannel(o.repositories)

	var wg sync.WaitGroup
	resultCh := make(chan *repositories.RepositoryDTO[K], o.worksSize)
	for w := 0; w < o.worksSize; w++ {
		wg.Add(1)
		go func() {
			for repository, ok := <-randomRepositories; ok; repository, ok = <-randomRepositories {
				//res, err := worker.work(repository)
				res, err := singleOperation(repository)
				if err == nil {
					resultCh <- res
					break
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if len(resultCh) != o.worksSize {
		return nil, errors.New("Internal error: Not enough successful workers")
	}

	ans := <-resultCh
	for result := range resultCh {
		if ans.IsOlderThan(result) {
			ans = result
		}
	}
	close(resultCh)

	return ans, nil
}

func (o *PigeonholeOrchestrator[T, K]) ExecuteMultipleOperation(
	multipleOperation func(*T) (map[string]*repositories.RepositoryDTO[K], error),
) (res map[string]*repositories.RepositoryDTO[K], err error) {
	if len(*o.repositories) < o.worksSize {
		return res, errors.New("Internal error: Not enough repositories")
	}
	randomRepositories := helpers.NewRandomChannel(o.repositories)

	var wg sync.WaitGroup
	resultCh := make(chan map[string]*repositories.RepositoryDTO[K], o.worksSize)
	for w := 0; w < o.worksSize; w++ {
		wg.Add(1)
		go func() {
			finished := false
			for repository, ok := <-randomRepositories; ok; repository, ok = <-randomRepositories {
				//res, err := worker.work(repository)
				res, err := multipleOperation(repository)
				if err == nil {
					resultCh <- res
					finished = true
					break
				}
			}
			if !finished {
				//use context
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if len(resultCh) != o.worksSize {
		return res, errors.New("Internal error: Not enough successful workers")
	}

	countMap := map[string]int{}
	valueMap := map[string]*repositories.RepositoryDTO[K]{}
	for resultMap := range resultCh {
		for key, aux := range resultMap {
			value, ok := valueMap[key]
			count := countMap[key]
			if !ok {
				value = aux
				count = 0
			} else if value.IsOlderThan(aux) {
				value = aux
			}
			count++
			valueMap[key] = value
			countMap[key] = count
		}
	}

	for key, count := range countMap {
		if count != o.worksSize {
			delete(valueMap, key)
		}
	}

	return valueMap, nil
}

func NewPigeonholeOrchestrator[T any, K any](
	repositories *[]*T,
) (*PigeonholeOrchestrator[T, K], error) {
	if repositories == nil {
		return nil, errors.New("Internal error: Repositories is a nil pointer")
	}

	size := len(*repositories)
	if size == 0 {
		return nil, errors.New("Internal error: Repositories array is empty")
	}

	return &PigeonholeOrchestrator[T, K]{
		worksSize:    size/2 + 1,
		repositories: repositories,
	}, nil
}
