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

type valueCount[T any] struct {
	Value T
	Count int
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

	valueCountMap := map[string]valueCount[*repositories.RepositoryDTO[K]]{}
	for resultMap := range resultCh {
		for key, newValue := range resultMap {
			curr, ok := valueCountMap[key]
			if !ok {
				curr = valueCount[*repositories.RepositoryDTO[K]]{
					Value: newValue,
					Count: 0,
				}
			} else if curr.Value.IsOlderThan(newValue) {
				curr.Value = newValue
			}
			curr.Count++
			valueCountMap[key] = curr
		}
	}

	ans := map[string]*repositories.RepositoryDTO[K]{}
	for key, valueCount := range valueCountMap {
		if valueCount.Count == o.worksSize {
			ans[key] = valueCount.Value
		}
	}

	return ans, nil
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
