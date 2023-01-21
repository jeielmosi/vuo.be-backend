package ports

import (
	"errors"
	"sync"
)

type PigeonholeOrchestrator[T any, K any] struct {
	worksSize    int
	repositories *[]T
}

func (o *PigeonholeOrchestrator[T, K]) SingleOperation(
	worker SingleOperationFunction[T, *K],
) (res *K, err error) {
	if len(*o.repositories) < o.worksSize {
		return res, errors.New("Internal error: Not enough repositories")
	}
	randomRepositories := NewRandomChannel(o.repositories)

	var wg sync.WaitGroup
	resultCh := make(chan PigeonholeDTO[*K], o.worksSize)
	for w := 0; w < o.worksSize; w++ {
		wg.Add(1)
		go func() {
			for repository, ok := <-randomRepositories; ok; repository, ok = <-randomRepositories {
				res, err := worker.work(repository)
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
		return res, errors.New("Internal error: Not enough successful workers")
	}

	ans := <-resultCh
	for result := range resultCh {
		if ans < result {
			ans = result
		}
	}
	close(resultCh)

	return ans.Entity, nil
}

type valueCount[T any] struct {
	Value T
	Count uint
}

func (o *PigeonholeOrchestrator[T, K]) MultipleOperation(
	worker MultipleOperationFunction[T, K],
) (res map[string]K, err error) {
	if len(*o.repositories) < o.worksSize {
		return res, errors.New("Internal error: Not enough repositories")
	}
	randomRepositories := NewRandomChannel(o.repositories)

	var wg sync.WaitGroup
	resultCh := make(chan map[string]PigeonholeDTO[K], o.worksSize)
	for w := 0; w < o.worksSize; w++ {
		wg.Add(1)
		go func() {
			for repository, ok := <-randomRepositories; ok; repository, ok = <-randomRepositories {
				res, err := worker.work(repository)
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
		return res, errors.New("Internal error: Not enough successful workers")
	}

	valueCountMap := map[string]valueCount[PigeonholeDTO[K]]{}
	for resultMap := range resultCh {
		for key, newValue := range resultMap {
			if _, ok := valueCountMap[key]; !ok {
				valueCountMap[key] = valueCount[PigeonholeDTO[K]]{
					Value: newValue,
					Count: 0,
				}
			}
			if valueCountMap[key].Value < newValue {
				valueCountMap[key].Value = newValue
			}
			valueCountMap[key].Count++
		}
	}

	ans := map[string]K{}
	for key, valueCount := range valueCountMap {
		if valueCount.Count == o.worksSize {
			ans[key] = valueCount.Value.Entity
		}
	}

	return ans, nil
}

func NewPigeonholeOrchestrator[T any, K any](
	repositories *[]T,
) (*PigeonholeOrchestrator[T, K], error) {
	size := len(*repositories)
	if size == 0 {
		return nil, errors.New("Internal error: Not enough repositories")
	}

	return &PigeonholeOrchestrator[T, K]{
		worksSize:    size/2 + 1,
		repositories: repositories,
	}, nil
}
